package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type SuccessResponse struct {
	Result any `json:"result"`
}

type ErrorResponse struct {
	Error any `json:"error"`
}

func ErrResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&ErrorResponse{Error: data})
}

func OkResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&SuccessResponse{Result: data})
}

func ParseBody(r *http.Request, x interface{}) {
	reqBody, _ := io.ReadAll(r.Body)
	r.Body.Close()

	err := json.Unmarshal(reqBody, x)
	if err != nil {
		return
	}
}
