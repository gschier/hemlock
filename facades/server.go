package facades

import (
	"net/http"
)

type Server interface {
	http.Server
}