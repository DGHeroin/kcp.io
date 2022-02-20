package main

import kio "github.com/DGHeroin/kcp.io"

func main() {
    s := kio.NewServer(nil)
    s.OnConnect(func(conn kio.Conn) error {
        return nil
    })
    s.OnDisconnect(func(conn kio.Conn) {

    })
    s.OnError(func(conn kio.Conn, err error) {

    })
    s.OnEvent(func(conn kio.Conn, bytes []byte) {

    })
}
