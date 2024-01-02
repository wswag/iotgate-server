package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"wswagner.visualstudio.com/iotgate-server/v1/services"
)

// DevicesController handles device management related api features
type PublicKeyController struct {
	Dds services.DeviceDataService
}

// Hook hooks onto a given subrouter and provides a device management api
func (dc PublicKeyController) Hook(mux *mux.Router) {
	mux.HandleFunc("/", dc.handleGetPublicKey).Methods("GET")
}

func (dc PublicKeyController) handleGetPublicKey(writer http.ResponseWriter, r *http.Request) {
	content, err := services.GetPublicKeyPEM()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(writer, string(content))
}
