package kio

const (
    pTypePing = iota
    pTypePong
    pTypeRequest
    pTypeResponse
)

type packet struct {
    Id      int    `json:"id"`
    Type    int    `json:"type"`
    Payload []byte `json:"payload"`
}
