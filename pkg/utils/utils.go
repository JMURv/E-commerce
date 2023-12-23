package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func ParseBody(r *http.Request, x interface{}) {
	reqBody, _ := io.ReadAll(r.Body)
	r.Body.Close()

	err := json.Unmarshal(reqBody, x)
	if err != nil {
		return
	}
}
