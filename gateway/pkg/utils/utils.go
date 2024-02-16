package utils

import (
	"encoding/json"
	"io"
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
	json.NewEncoder(w).Encode(&ErrorResponse{Error: data})
}

func OkResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(&SuccessResponse{data})
	if err != nil {
		ErrResponse(w, http.StatusInternalServerError, ErrWhileEncoding)
		return
	}
}

func ParseBody(r *http.Request, x interface{}) {
	reqBody, _ := io.ReadAll(r.Body)
	r.Body.Close()

	err := json.Unmarshal(reqBody, x)
	if err != nil {
		return
	}
}
