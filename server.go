package kio

import (
    "encoding/json"
    "github.com/xtaci/kcp-go/v5"
    "log"
    "net"
    "sync/atomic"
    "time"
)

type (
    Server struct {
        onConnect      func(conn Conn) error
        onDisconnect   func(conn Conn)
        onError        func(conn Conn, err error)
        onEvent        func(Conn, []byte)
        opt            *ServerOption
        ln             net.Listener
        rChan          chan *iEventMsg
        wChan          chan *iEventMsg
        writeTimeout   time.Duration
        clientId       int64
        remoteClients  map[int64]*remoteClient
        clientCount    int64
        sendBytesCount uint64
        recvBytesCount uint64
    }
    ServerOption struct {
        RecvBufferQueue uint
        SendBufferQueue uint
        ReadTimeout     time.Duration
        WriteTimeout    time.Duration
    }
)

func DefaultServerOption() *ServerOption {
    return &ServerOption{
        RecvBufferQueue: 0,
        SendBufferQueue: 0,
        ReadTimeout:     30 * time.Second,
        WriteTimeout:    30 * time.Second,
    }
}
func NewServer(opt *ServerOption) *Server {
    if opt == nil {
        opt = DefaultServerOption()
    }
    return &Server{
        rChan: make(chan *iEventMsg, opt.RecvBufferQueue),
        wChan: make(chan *iEventMsg, opt.SendBufferQueue),
    }
}
func (s *Server) Serve(addrs ...string) {
    var addr = ":1989"
    if len(addrs) == 1 {
        addr = addrs[0]
    }
    ln, err := kcp.Listen(addr)
    if err != nil {
        return
    }
    go s.handleRead()
    go s.handleWrite()
    s.handleAccept(ln)
}
func (s *Server) serveConn(conn net.Conn) {
    var cli = &remoteClient{
        id:       atomic.AddInt64(&s.clientId, 1),
        sendChan: s.wChan,
    }
    atomic.AddInt64(&s.clientCount, 1)
    if err := s.onConnect(cli); err != nil {
        return
    }

    for {
        payload, err := readHeadMessage(conn, time.Second*10, &s.recvBytesCount)
        if err != nil {
            s.onError(cli, err)
            break
        }
        var pkt packet
        if err := json.Unmarshal(payload, &pkt); err != nil {
            s.onError(cli, ErrorMessageCorrupt)
            cli.Close()
            break
        }
        evt := askEventMsg()
        evt.conn = cli
        evt.payload = &packet{
            Type:    0,
            Payload: nil,
        }
        s.rChan <- evt
    }
    s.onDisconnect(cli)
    atomic.AddInt64(&s.clientCount, -1)
}
func (s *Server) OnConnect(cb func(conn Conn) error) {
    s.onConnect = cb
}
func (s *Server) OnDisconnect(cb func(Conn)) {
    s.onDisconnect = cb
}
func (s *Server) OnError(cb func(Conn, error)) {
    s.onError = cb
}
func (s *Server) OnEvent(cb func(conn Conn, payload []byte)) {
    s.onEvent = cb
}
func (s *Server) Count() int64 {
    return atomic.LoadInt64(&s.clientCount)
}
func (s *Server) handleRead() {
    pOnEvent := func(cli *remoteClient, msg *iEventMsg) {
        defer func() {
            recover()
        }()
        p := msg.payload
        switch p.Type {
        case 0: // ping
            msg := askEventMsg()
            msg.payload = &packet{
                Id:      cli.NextId(),
                Type:    0,
                Payload: nil,
            }
            cli.sendChan <- msg
        case 1: // pong
            return
        case 2: // req
            s.onEvent(cli, p.Payload)
        case 3: // rsp
        }
        log.Println(p)

    }
    go func() {
        for {
            e := <-s.rChan
            if e == nil {
                break
            }
            atomic.AddUint64(&s.recvBytesCount, e.size)
            pOnEvent(e.conn, e)
            relEventMsg(e)
        }
    }()
}
func (s *Server) handleWrite() {
    pOnEvent := func(cli *remoteClient, msg *iEventMsg) {
        defer func() {
            recover()
        }()
        var bin []byte
        n, err := writeHeadMessage(cli.conn, bin, s.writeTimeout, &s.sendBytesCount)
        if err != nil {
            s.onError(cli, err)
            return
        }
        atomic.AddUint64(&s.sendBytesCount, uint64(n))
    }
    go func() {
        for {
            e := <-s.wChan
            if e == nil {
                break
            }
            pOnEvent(e.conn, e)
            relEventMsg(e)
        }
    }()
}
func (s *Server) handleAccept(ln net.Listener) {
    for {
        conn, err := ln.Accept()
        if err == nil {
            s.serveConn(conn)
        }
    }
}
