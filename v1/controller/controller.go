package controller

import (
	"github.com/gorilla/mux"
)

type Controller interface {
	Hook(router *mux.Router)
}
