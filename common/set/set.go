package set

import "encoding/json"

type StringSet map[string]struct{}

var (
	_ json.Marshaler   = (StringSet)(nil)
	_ json.Unmarshaler = (*StringSet)(nil)
)

func (s StringSet) Has(item string) (ok bool) {
	_, ok = s[item]
	return
}

func (s StringSet) Add(item string) {
	s[item] = struct{}{}
}

func (s StringSet) Delete(item string) {
	delete(s, item)
}

func (s StringSet) MarshalJSON() (data []byte, err error) {
	slice := make([]string, 0, len(s))
	for item := range s {
		slice = append(slice, item)
	}
	return json.Marshal(slice)
}

func (s *StringSet) UnmarshalJSON(data []byte) error {
	var slice []string
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	*s = NewStringSetFromSlice(slice)
	return nil
}

func NewStringSetFromSlice(slice []string) StringSet {
	s := make(StringSet, len(slice))
	for _, item := range slice {
		s.Add(item)
	}
	return s
}
