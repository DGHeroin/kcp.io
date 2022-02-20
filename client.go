package kio

import (
    "net"
    "strconv"
)

type client struct {
    conn     net.Conn
    id       int64
    ctx      interface{}
    sendChan chan *iEventMsg
}

func (c *client) ID() string {
    return strconv.FormatInt(c.id, 36)
}

func (c *client) Close() error {
    return c.conn.Close()
}

func (c *client) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

func (c *client) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}

func (c *client) Context() interface{} {
    return c.ctx
}

func (c *client) SetContext(v interface{}) {
    c.ctx = v
}

func (c *client) Emit(msg string, payload []byte) {
    e := askEventMsg()
    e.conn = c
    e.payload = payload
    c.sendChan <- e
}

//func (c *client) Join(room string) {
//    panic("implement me")
//}
//
//func (c *client) Leave(room string) {
//    panic("implement me")
//}
//
//func (c *client) LeaveAll() {
//    panic("implement me")
//}
//
//func (c *client) Rooms() []string {
//    panic("implement me")
//}
