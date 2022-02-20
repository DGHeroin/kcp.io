package kio

import (
    "net"
    "strconv"
)

type (
    remoteClient struct {
        conn     net.Conn
        id       int64
        ctx      interface{}
        sendChan chan *iEventMsg
    }
)

func (c *remoteClient) ID() string {
    return strconv.FormatInt(c.id, 36)
}

func (c *remoteClient) Close() error {
    return c.conn.Close()
}

func (c *remoteClient) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

func (c *remoteClient) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}

func (c *remoteClient) Context() interface{} {
    return c.ctx
}

func (c *remoteClient) SetContext(v interface{}) {
    c.ctx = v
}

func (c *remoteClient) Emit(msg string, payload []byte) {
    e := askEventMsg()
    e.conn = c
    e.payload = payload
    c.sendChan <- e
}
