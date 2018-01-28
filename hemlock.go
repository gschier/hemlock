package hemlock

import (
	"fmt"
	"os"
	"time"
)

// CacheBustKey is the cache busting key
var CacheBustKey string

// Env returns the value of the 'name'd environment variable or an empty string
func Env(name string) string {
	return os.Getenv(name)
}

// EnvOr will return the value for 'name' or the fallback if it doesn't exist
func EnvOr(name, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return value
}

// EnvOrPanic will return the value for 'name' or panic if not present
func EnvOrPanic(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic("Environment variable " + name + " must be set")
	}

	return value
}

func init() {
	CacheBustKey = fmt.Sprintf("%d", time.Now().Unix())
}
