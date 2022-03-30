package kio

import (
    "net"
    "strconv"
    "sync/atomic"
)

type (
    remoteClient struct {
        conn     net.Conn
        id       int64
        ctx      interface{}
        sendChan chan *iEventMsg
        msgId    int32
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
    e.payload = &packet{
        Type:    pTypeRequest,
        Payload: payload,
    }
    c.sendChan <- e
}

func (c *remoteClient) NextId() int {
    return int(atomic.AddInt32(&c.msgId, 1))
}
