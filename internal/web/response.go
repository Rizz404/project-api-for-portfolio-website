package web

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
)

// * Buat yang offset base pagination
type PageInfo struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	HasPrevPage bool  `json:"has_prev_page"`
	HasNextPage bool  `json:"has_next_page"`
}

// * Buat yang cursor base pagination
type CursorInfo struct {
	NextCursor  string `json:"next_cursor"`
	HasNextPage bool   `json:"has_next_page"`
	PerPage     int    `json:"per_page"`
	Total       int64  `json:"total,omitempty"`
}

type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	// * Tergantung datanya jadi bisa gak ada
	PageInfo   *PageInfo   `json:"pagination,omitempty"`
	CursorInfo *CursorInfo `json:"cursor,omitempty"`
}

func Success(w http.ResponseWriter, code int, message string, data any) {
	writeJSON(w, code, JSONResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessWithPageInfo(w http.ResponseWriter, code int, message string, data any, total int64, perPage int, currentPage int) {
	totalPages := 0
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	pageInfo := &PageInfo{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		HasPrevPage: currentPage > 1,
		HasNextPage: currentPage < totalPages,
	}

	writeJSON(w, code, JSONResponse{
		Status:   "success",
		Message:  message,
		Data:     data,
		PageInfo: pageInfo,
	})
}

func SuccessWithCursor(w http.ResponseWriter, code int, message string, data any, nextCursor string, hasNextPage bool, perPage int, total ...int64) {
	cursorInfo := &CursorInfo{
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
		PerPage:     perPage,
	}

	if len(total) > 0 {
		cursorInfo.Total = total[0]
	}

	writeJSON(w, code, JSONResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		CursorInfo: cursorInfo,
	})
}

func Error(w http.ResponseWriter, code int, message string, errorDetails any) {
	writeJSON(w, code, JSONResponse{
		Status:  "error",
		Message: message,
		Error:   errorDetails,
	})
}

// * Otomatis handle error berdasarkan domain error (abstraksi terus sampe mampus)
func HandleError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			// * Log error aslinya dari AppError
			log.Printf("ERROR: Internal server error handled: %v", appErr.Unwrap())
		}
		Error(w, appErr.Code, appErr.Error(), nil)
	} else {
		log.Printf("ERROR: Internal server error: %v", err)
		// * Selalu kirim error generic ke client bukan detail
		Error(w, http.StatusInternalServerError, "An unexpected error occurred", err)
	}
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
