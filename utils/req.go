package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson(data interface{}, w http.ResponseWriter, code int) {
	var respBody []byte
	if data != nil {
		var err error
		respBody, err = json.Marshal(data)
		if err != nil {
			log.Printf("Error parsing response json: %v", err)
		}
	}
	w.WriteHeader(code)
	w.Write(respBody)

}

func WriteError(err error, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write(([]byte)(err.Error()))
}