package interfaces

import (
	"context"
	"io"
)

type Request interface {
	// Input grabs input from the request body by name
	Input(name string) string

	// Input grabs input from the query string by name
	Query(name string) string

	// Cookie grabs input from cookies by name
	Cookie(name string) string

	// File grabs the file
	File(name string) io.Reader

	// Context returns the context.Context of the current request
	Context() context.Context
}
