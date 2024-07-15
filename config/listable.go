package config

import (
	"encoding/json"
	"errors"
)

type Listable[T any] []T

var (
	_ json.Marshaler   = (*Listable[struct{}])(nil)
	_ json.Unmarshaler = (*Listable[struct{}])(nil)
)

func (l *Listable[T]) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*[]T)(l))
	if err == nil {
		return nil
	}
	var single T
	err2 := json.Unmarshal(bytes, &single)
	if err2 != nil {
		return errors.Join(err, err2)
	}
	*l = []T{single}
	return nil
}

func (l Listable[T]) MarshalJSON() ([]byte, error) {
	if len(l) == 1 {
		return json.Marshal(l[0])
	}
	return json.Marshal([]T(l))
}
