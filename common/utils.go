package common

var CRLF = []byte("\r\n")

func Must[T any](s T, err error) T {
	if err != nil {
		panic(err)
	}
	return s
}

func MustOK[T any](s T, ok bool) T {
	if ok {
		return s
	}
	panic("assertion failed")
}

func Must0(err error) {
	if err != nil {
		panic(err)
	}
}

func GetSecond[T any](_ any, r T) T {
	return r
}
