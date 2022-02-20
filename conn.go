package kio

import (
    "net"
)

type Conn interface {
    ID() string
    Close() error

    LocalAddr() net.Addr
    RemoteAddr() net.Addr

    Context() interface{}
    SetContext(v interface{})

    Emit(msg string, payload []byte)
}
