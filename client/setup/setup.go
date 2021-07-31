package setup

import (
	"os"
)

type S struct{}

func New() *S {
	return &S{}
}

func (s *S) Host() string {
	return readString(envHost, envHostDefault)
}

func readString(key, defaultValue string) (value string) {
	if value = os.Getenv(key); len(value) > 0 {
		return
	}
	return defaultValue
}

func (s *S) ServerHost() string {
	return readString(envServerHost, envServerHostDefault)
}
