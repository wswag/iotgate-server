package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// HomeHandler returns the firmware server version string
func homeHandler(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(writer, "iotgate-server 0.1")
}

func deviceIdRegisterHandler(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	fmt.Fprintf(writer, "device %s registered", v["id"])
	log.Println(fmt.Sprintf("device %s registered", v["id"]))
}

func deviceIdFirmwareHandler(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(writer, "v1.2.44")
}

func deviceIdFirmwareBlockHandler(writer http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	fn := "firmware/" + v["id"] + ".bin"

	start, _ := strconv.Atoi(r.FormValue("start"))
	len, _ := strconv.Atoi(r.FormValue("len"))
	file, err := os.Open(fn)

	if err != nil {
		fmt.Fprintf(writer, err.Error())
		return
	}
	defer file.Close()
	stat, _ := file.Stat()

	if len == 0 {
		len = int(stat.Size()) - start
	}

	if int64(start+len) > stat.Size() {
		fmt.Fprint(writer, "invalid length")
		return
	}

	buf := make([]byte, len)
	file.Seek(int64(start), 0)
	file.Read(buf)

	//aw := base64.StdEncoding.EncodeToString(buf)
	//fmt.Fprintf(writer, aw)
	writer.Write(buf)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/devices/{id}/register", deviceIdRegisterHandler).Methods("POST")
	r.HandleFunc("/devices/{id}/firmware", deviceIdFirmwareHandler)
	r.HandleFunc("/devices/{id}/firmware/block", deviceIdFirmwareBlockHandler).Queries("start", "{start:[0-9]+}", "len", "{len:[0-9]+}")

	http.ListenAndServe("127.0.0.1:8080", r)
	log.Println("server stopped.")
}
