package gapt

import (
	"crypto/sha512"
	"time"
)

func nowAsUnixMilliseconds() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}

func hash(data []byte) [32]byte {
	return sha512.Sum512_256(data)
}

// assert panics if the assertion applies.
func assertNil(any interface{}) {
	if err, ok := any.(error); ok {
		panic(err)
	}
}
