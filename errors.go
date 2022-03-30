package kio

import "errors"

var (
    ErrorMessageTooLarge = errors.New("message payload too large")
    ErrorMessageCorrupt = errors.New("message format corrupt")
)
