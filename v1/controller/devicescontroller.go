package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"wswagner.visualstudio.com/iotgate-server/v1/model"

	"github.com/gorilla/mux"
	"wswagner.visualstudio.com/iotgate-server/v1/services"
)

// DevicesController handles device management related api features
type DevicesController struct {
	Dds services.DeviceDataService
}

// Hook hooks onto a given subrouter and provides a device management api
func (dc DevicesController) Hook(mux *mux.Router) {
	mux.HandleFunc("/", dc.handleGetHome).Methods("GET")
	mux.HandleFunc("/{id}", dc.handleGetID).Methods("GET")
	mux.HandleFunc("/{id}/register", dc.handlePostIDRegister).Methods("POST")
}

func (dc DevicesController) handleGetHome(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(writer, "yet to come")
}

func (dc DevicesController) handleGetID(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	deviceID := v["id"]
	deviceData, err := dc.Dds.GetDeviceData(deviceID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	deviceData.PrivateKey = "secret" // obfuscate the private key
	data, _ := json.Marshal(deviceData)
	writer.Write(data)
}

func (dc DevicesController) handlePostIDRegister(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	//fmt.Fprintln(writer, "device %s registered", v["id"])
	//log.Println(fmt.Sprintf("device %s registered", v["id"]))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	deviceData := model.DeviceData{}
	json.Unmarshal(body, &deviceData)
	deviceData.DeviceID = v["id"] // override id
	dc.Dds.SetDeviceData(deviceData)
}
