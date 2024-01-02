package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"wswagner.visualstudio.com/iotgate-server/v1/controller"
	"wswagner.visualstudio.com/iotgate-server/v1/middleware"
	"wswagner.visualstudio.com/iotgate-server/v1/services"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("iotgate-server 1.0.0")

	r := mux.NewRouter()

	appPort := strings.TrimSpace(os.Getenv("PORT"))
	if appPort == "" {
		appPort = "80"
	}
	log.Println("configured port (PORT env): " + appPort)

	storageLocation := strings.TrimSpace(os.Getenv("STORAGE_PATH"))
	if storageLocation == "" {
		storageLocation = "./data"
	}
	log.Println("configured stroage location (STORAGE_PATH env): " + storageLocation)

	pk_err := services.TestPrivateKey()
	if pk_err != nil {
		log.Fatalln("erronous X509 RSA private key configuration (SIGNATURE_PRIVATE_KEYFILE env): " + pk_err.Error())
	}
	log.Println("signature private key OK")

	pk_err = services.TestPublicKey()
	if pk_err != nil {
		log.Fatalln("erronous X509 RSA public key configuration (SIGNATURE_PUBLIC_KEYFILE env): " + pk_err.Error())
	}
	log.Println("signature public key OK")

	mdfs := services.ManagedDocumentFileService{Basepath: storageLocation, StoragePattern: "%s.%s"}
	dds := services.MFDeviceDataService{FileService: mdfs}
	fws := services.MFFirmwareService{FileService: mdfs}

	dc := controller.DevicesController{Dds: dds}
	fc := controller.FirmwareController{Dds: dds, Fws: fws}
	pc := controller.PublicKeyController{}

	devices := r.PathPrefix("/devices/").Subrouter()
	firmware := r.PathPrefix("/firmware/").Subrouter()
	publickey := r.PathPrefix("/publickey").Subrouter()

	dc.Hook(devices)
	fc.Hook(firmware)
	pc.Hook(publickey)

	r.Use(middleware.LoggingMiddleware)

	log.Println("starting server at port :" + appPort)
	err := http.ListenAndServe(":"+appPort, r)
	log.Println("server stopped. " + err.Error())
}
