package jsonx

import (
	"encoding"
	"time"
)

type Duration time.Duration

var (
	_ encoding.TextMarshaler   = Duration(0)
	_ encoding.TextUnmarshaler = (*Duration)(nil)
)

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	*d = Duration(duration)
	return err
}
