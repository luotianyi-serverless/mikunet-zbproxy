package config

import (
	"encoding/json"
	"fmt"

	"github.com/phuslu/log"
)

type _Log struct {
	Level string
}

type Log struct {
	Level log.Level
}

var (
	_ json.Marshaler   = (*Log)(nil)
	_ json.Unmarshaler = (*Log)(nil)
)

func (l *Log) UnmarshalJSON(bytes []byte) error {
	var rawLog _Log
	err := json.Unmarshal(bytes, &rawLog)
	if err != nil {
		return err
	}
	if rawLog.Level == "" {
		l.Level = log.DebugLevel // default log level
	} else {
		l.Level = log.ParseLevel(rawLog.Level)
		if l.Level > log.PanicLevel { // the last enum value
			return fmt.Errorf("unknown log level: %s", rawLog.Level)
		}
	}
	return nil
}

func (l Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(_Log{
		Level: l.Level.String(),
	})
}
