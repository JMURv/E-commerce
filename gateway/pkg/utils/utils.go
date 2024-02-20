package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var ErrWhileEncoding = "Error while encoding to JSON"

type SuccessResponse struct {
	Result any `json:"result"`
}

type ErrorResponse struct {
	Error any `json:"error"`
}

func ErrResponse(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&ErrorResponse{Error: data})
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func OkResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(&SuccessResponse{Result: data})
	if err != nil {
		log.Println(err.Error())
		ErrResponse(w, http.StatusInternalServerError, ErrWhileEncoding)
		return
	}
}

func ParseBody(r *http.Request, x interface{}) {
	rBody, _ := io.ReadAll(r.Body)
	r.Body.Close()

	err := json.Unmarshal(rBody, x)
	if err != nil {
		return
	}
}
