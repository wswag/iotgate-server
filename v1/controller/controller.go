package controller

import (
	"github.com/gorilla/mux"
)

// Controller interface exposes a Hook function to bind to a given router
type Controller interface {
	Hook(router *mux.Router)
}
