package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	resp, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		statusCode = http.StatusInternalServerError
		resp = []byte(`{"message": "could not generate response"}`)
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(resp)
	return
}
