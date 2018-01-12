package hemlock

import (
	"os"
)

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
