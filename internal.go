package kio

import "sync"

type (
    iEventMsg struct {
        conn    *remoteClient
        payload *packet
        size    uint64
    }
)

var (
    iEventMsgPool = sync.Pool{
        New: func() interface{} {
            return &iEventMsg{}
        },
    }
)

func askEventMsg() *iEventMsg {
    return iEventMsgPool.Get().(*iEventMsg)
}
func relEventMsg(p *iEventMsg) {
    iEventMsgPool.Put(p)
}
