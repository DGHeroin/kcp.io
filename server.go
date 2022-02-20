package kio

import (
    "encoding/binary"
    "github.com/xtaci/kcp-go/v5"
    "io"
    "net"
    "time"
)

type (
    Server struct {
        onConnect    func(conn Conn) error
        onDisconnect func(conn Conn)
        onError      func(conn Conn, err error)
        onEvent      func(Conn, []byte)
        opt          *ServerOption
        ln           net.Listener
        rChan        chan *iEventMsg
        wChan        chan *iEventMsg
        writeTimeout time.Duration
    }
    ServerOption struct {
        BufferQueue  uint
        ReadTimeout  time.Duration
        WriteTimeout time.Duration
    }
)

var (
    DefaultServerOption = &ServerOption{
        BufferQueue:  0,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
)

func NewServer(opt *ServerOption) *Server {
    if opt == nil {
        opt = DefaultServerOption
    }
    return &Server{
        rChan: make(chan *iEventMsg, opt.BufferQueue),
        wChan: make(chan *iEventMsg, opt.BufferQueue),
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
    var cli = &client{
        sendChan: s.wChan,
    }
    if err := s.onConnect(cli); err != nil {
        return
    }

    for {
        payload, err := readHeadMessage(conn, time.Second*10)
        if err != nil {
            break
        }
        evt := askEventMsg()
        evt.conn = cli
        evt.payload = payload
        s.rChan <- evt
    }
    s.onDisconnect(cli)
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
func (s *Server) OnEvent(cb func(Conn, []byte)) {
    s.onEvent = cb
}

func (s *Server) handleRead() {
    pOnEvent := func(cli *client, payload []byte) {
        defer func() {
            recover()
        }()
        s.onEvent(cli, payload)
    }
    go func() {
        for {
            e := <-s.rChan
            if e == nil {
                break
            }
            pOnEvent(e.conn, e.payload)
            relEventMsg(e)
        }
    }()
}

func (s *Server) handleWrite() {
    pOnEvent := func(cli *client, payload []byte) {
        defer func() {
            recover()
        }()
        _, err := writeHeadMessage(cli.conn, payload, s.writeTimeout)
        if err != nil {
            s.onError(cli, err)
            return
        }
    }
    go func() {
        for {
            e := <-s.wChan
            if e == nil {
                break
            }
            pOnEvent(e.conn, e.payload)
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

func readHeadMessage(conn net.Conn, timeout time.Duration) ([]byte, error) {
    var header = make([]byte, 4)
    if timeout > 0 {
        if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
            return nil, err
        }
    }
    if n, err := io.ReadFull(conn, header); n != 4 || err != nil {
        return nil, err
    }
    sz := binary.BigEndian.Uint32(header)
    if sz > cMessagePayloadSize {
        return nil, ErrorMessageTooLarge
    }
    payload := make([]byte, sz)
    if timeout > 0 {
        if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
            return nil, err
        }
    }
    if _, err := io.ReadFull(conn, payload); err != nil {
        return nil, err
    }
    return payload, nil
}
func writeHeadMessage(conn net.Conn, payload []byte, timeout time.Duration) (int, error) {
    sz := uint32(len(payload))

    if sz > cMessagePayloadSize {
        return 0, ErrorMessageTooLarge
    }

    var header = make([]byte, 4)
    binary.BigEndian.PutUint32(header, sz)

    if timeout > 0 {
        if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
            return -0, err
        }
    }
    if sz == 0 {
        return conn.Write(header)
    } else {
        return conn.Write(append(header, payload...))
    }

}
