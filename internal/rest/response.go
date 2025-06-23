package rest

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, code int, message string, data any) {
	writeJSON(w, code, JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string, errorDetails any) {
	writeJSON(w, code, JSONResponse{
		Status:  "error",
		Message: message,
		Error:   errorDetails,
	})
}

// * Fungsi internal untuk proses encoding JSON dan penulisan header.
func writeJSON(w http.ResponseWriter, code int, payload any) {
	// * Marshal payload ke JSON dari any/struct ke json
	response, err := json.Marshal(payload)
	if err != nil {
		// * Jika ada error saat marshalling, kirim error internal server
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	// * Set header dan tulis respons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
