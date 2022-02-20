package kio

import "sync"

type (
    iEventMsg struct {
        conn    *client
        payload []byte
    }
    iErrMsg struct {
        conn Conn
        err  error
    }
)

var (
    iEventMsgPool = sync.Pool{
        New: func() interface{} {
            return iEventMsg{}
        },
    }
)

func askEventMsg() *iEventMsg {
    return iEventMsgPool.Get().(*iEventMsg)
}
func relEventMsg(p *iEventMsg) {
    iEventMsgPool.Put(p)
}
