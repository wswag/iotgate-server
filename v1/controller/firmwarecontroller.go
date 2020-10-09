package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"wswagner.visualstudio.com/iotgate-server/v1/services"
)

// FirmwareController provides device firmware related apis
type FirmwareController struct {
	Dds services.DeviceDataService
	Fws services.FirmwareService
}

// Hook hooks onto a given subrouter and provides a device management api
func (fc FirmwareController) Hook(mux *mux.Router) {
	mux.HandleFunc("/", fc.handleGetHome).Methods("GET")
	mux.HandleFunc("/{id}", fc.handleGetID).Methods("GET")
	mux.HandleFunc("/{id}/image", fc.handleGetIDImage).Methods("GET")
	mux.HandleFunc("/{id}/image", fc.handlePostIDImage).Methods("POST")
	mux.HandleFunc("/{id}/image/chunk", fc.handleGetIDImageChunk).Queries("start", "{start:[0-9]+}", "len", "{len:[0-9]+}")
}

func (fc FirmwareController) handleGetHome(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(writer, "yet to come")
}

func (fc FirmwareController) handleGetID(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	deviceID := v["id"]
	meta, err := fc.Fws.GetFirmwareMetadata(deviceID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(meta)
	writer.Write(data)
}

func (fc FirmwareController) handleGetIDImage(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	deviceID := v["id"]
	err := fc.Fws.DownloadFirmware(deviceID, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func (fc FirmwareController) handlePostIDImage(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	deviceID := v["id"]

	deviceData, err := fc.Dds.GetDeviceData(deviceID)
	if err != nil {
		// most likely, device is not registered
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	//contentType := r.Header.Get("Content-Type")
	//if (contentType != "application/octet-stream")

	// TODO: where to put business logic?
	meta, _ := fc.Fws.GetFirmwareMetadata(deviceID)
	meta.DeviceID = deviceID
	meta.Iteration++
	meta.Size = uint32(r.ContentLength)
	meta.Timestamp = time.Now().Unix()

	// reset the hash
	meta.SHAHash = ""
	err = fc.Fws.UploadFirmware(deviceID, r.Body, &meta)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	h := hmac.New(sha256.New, []byte(deviceData.PrivateKey))
	h.Write([]byte(meta.SHAHash))

	meta.Signature = hex.EncodeToString(h.Sum(nil))

	fc.Fws.SetFirmwareMetadata(meta)
}

func (fc FirmwareController) handleGetIDImageChunk(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	deviceID := v["id"]
	start, _ := strconv.Atoi(r.FormValue("start"))
	len, _ := strconv.Atoi(r.FormValue("len"))
	data, err := fc.Fws.DownloadFirmwareChunk(deviceID, uint32(start), uint32(len))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	writer.Write(data)
}
