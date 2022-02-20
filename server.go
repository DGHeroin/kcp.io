package kio

import (
    "encoding/binary"
    "encoding/json"
    "github.com/xtaci/kcp-go/v5"
    "io"
    "net"
    "time"
)

type (
    Server struct {
        ln net.Listener
    }
)

func NewServer() *Server {
    return &Server{}
}
func (s *Server) Serve(addrs ...string) {
    var addr = ":6677"
    if len(addrs) == 1 {
        addr = addrs[0]
    }
    ln, err := kcp.Listen(addr)
    if err != nil {
        return
    }
    s.ln = ln
    for {
        conn, err := ln.Accept()
        if err == nil {
            s.handleEcho(conn)
        }
    }
}
func (s *Server) handleEcho(conn net.Conn) {
    for {
        payload, err := readHeadMessage(conn, time.Second*10)
        if err != nil {
            break
        }
        var msgs []GameMessage
        json.Unmarshal(payload, msgs)

    }
}

// 146 166 228 189 160 229 165 189 100
type GameMessage struct {
    Type int
    Body []byte
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
    var header = make([]byte, 4)
    sz := uint32(len(payload))
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
