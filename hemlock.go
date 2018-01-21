package hemlock

import (
	"fmt"
	"os"
	"time"
)

var CacheBustKey string

func Env(name string) string {
	return os.Getenv(name)
}

func EnvOr(name, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return value
}

func init () {
	CacheBustKey = fmt.Sprintf("%d", time.Now().Unix())
}
