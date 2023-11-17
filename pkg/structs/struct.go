package structs

import "time"

type Value struct {
	Val        any
	Expiration time.Time
}

func NewValue(value any, expiration time.Duration) Value {
	return Value{Val: value, Expiration: time.Now().Add(expiration)}
}
