package kio

import (
    "github.com/xtaci/kcp-go/v5"
    "net"
)

type (
    Client struct {
        onConnect    func()
        onDisconnect func()
        onError      func(err error)
        onEvent      func([]byte)

        conn     net.Conn
        id       int64
        ctx      interface{}
        sendChan chan *iEventMsg
    }
    ClientOption struct {
    }
)

func NewClient(opt *ClientOption) *Client {
    cli := &Client{}

    return cli
}

func (s *Client) OnConnect(cb func()) {
    s.onConnect = cb
}
func (s *Client) OnDisconnect(cb func()) {
    s.onDisconnect = cb
}
func (s *Client) OnError(cb func(error)) {
    s.onError = cb
}
func (s *Client) OnEvent(cb func([]byte)) {
    s.onEvent = cb
}

func (s *Client) Connect(raddr string) error {
    conn, err := kcp.Dial(raddr)
    if err != nil {
        return err
    }
    s.conn = conn
    writeHeadMessage(conn, []byte("hello"), 0)
    return nil
}
