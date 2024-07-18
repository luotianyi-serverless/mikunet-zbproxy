package jsonx

import "encoding/json"

type RawJSON []byte

var (
	_ json.Marshaler   = RawJSON{}
	_ json.Unmarshaler = (*RawJSON)(nil)
)

func (r RawJSON) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	return r, nil
}

func (r *RawJSON) UnmarshalJSON(bytes []byte) error {
	*r = bytes
	return nil
}
