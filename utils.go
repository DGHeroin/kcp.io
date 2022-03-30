package kio

import (
    "encoding/binary"
    "io"
    "net"
    "sync/atomic"
    "time"
)

func readHeadMessage(conn net.Conn, timeout time.Duration, count *uint64) ([]byte, error) {
    var header = make([]byte, 4)
    if timeout > 0 {
        if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
            return nil, err
        }
    }
    if n, err := io.ReadFull(conn, header); n != 4 || err != nil {
        return nil, err
    }
    if count != nil {
        atomic.AddUint64(count, 4)
    }
    sz := binary.BigEndian.Uint32(header)
    if sz > cMessagePayloadSize {
        return nil, ErrorMessageTooLarge
    }
    payload := make([]byte, sz)
    if timeout > 0 {
        if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
            return nil, err
        }
    }
    if n, err := io.ReadFull(conn, payload); err != nil {
        return nil, err
    } else {
        if count != nil {
            atomic.AddUint64(count, uint64(n))
        }
        return payload, nil
    }

}
func writeHeadMessage(conn net.Conn, payload []byte, timeout time.Duration, count *uint64) (int, error) {
    sz := uint32(len(payload))

    if sz > cMessagePayloadSize {
        return 0, ErrorMessageTooLarge
    }

    var header = make([]byte, 4)
    binary.BigEndian.PutUint32(header, sz)

    if timeout > 0 {
        if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
            return -0, err
        }
    }
    var (
        n   int
        err error
    )
    if sz == 0 {
        n, err = conn.Write(header)
    } else {
        n, err = conn.Write(append(header, payload...))
    }
    if err == nil && count != nil {
        atomic.AddUint64(count, uint64(n))
    }
    return n, err
}
