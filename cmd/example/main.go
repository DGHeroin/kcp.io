package main

import (
    kio "github.com/DGHeroin/kcp.io"
    "log"
    "time"
)

func main() {
    go runClient()
    runServer()
}

func runClient() {
    time.Sleep(time.Second)
    c := kio.NewClient(nil)
    c.OnConnect(func() {
        log.Println("连接成功")
    })
    c.OnDisconnect(func() {
        log.Println("链接断开")
    })
    c.OnError(func(err error) {
        log.Println("发送错误:", err)
    })
    c.OnEvent(func(payload []byte) {
        log.Println("收到数据:", payload)
    })
    err := c.Connect("127.0.0.1:1989")
    if err != nil {
        log.Println("连接错误")
        return
    }
}

func runServer() {
    s := kio.NewServer(nil)
    s.OnConnect(func(conn kio.Conn) error {
        log.Println("接受客户端:", conn)
        return nil
    })
    s.OnDisconnect(func(conn kio.Conn) {
        log.Println("客户端断开:", conn)
    })
    s.OnError(func(conn kio.Conn, err error) {
        log.Println("发生错误:", conn, err)
    })
    s.OnEvent(func(conn kio.Conn, payload []byte) {
        log.Println("收到:", conn, payload)
    })

    s.Serve()
}
