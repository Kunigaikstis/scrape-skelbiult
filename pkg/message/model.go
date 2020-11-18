package message

import "time"

type Chat struct {
	Id       int64
	Username string
	Created  time.Time
}
