package model

import "time"

type RawTime []byte

func (t *RawTime) Time() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(*t))
}
