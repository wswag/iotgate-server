package main

import (
	"log"
	"net/http"
	"os"

	"wswagner.visualstudio.com/iotgate-server/v1/controller"
	"wswagner.visualstudio.com/iotgate-server/v1/middleware"
	"wswagner.visualstudio.com/iotgate-server/v1/services"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err.Error())
	}

	mdfs := services.ManagedDocumentFileService{Basepath: homedir + "/.iotgate-server", StoragePattern: "%s.%s"}
	dds := services.MFDeviceDataService{FileService: mdfs}
	fws := services.MFFirmwareService{FileService: mdfs}

	dc := controller.DevicesController{Dds: dds}
	fc := controller.FirmwareController{Dds: dds, Fws: fws}

	devices := r.PathPrefix("/devices/").Subrouter()
	firmware := r.PathPrefix("/firmware/").Subrouter()

	dc.Hook(devices)
	fc.Hook(firmware)

	r.Use(middleware.LoggingMiddleware)

	http.ListenAndServe(":8080", r)
	log.Println("server stopped.")
}
