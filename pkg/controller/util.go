package controller

import (
	"encoding/json"
	"fmt"
	"money_share/pkg/dto/response"
	"net/http"
	"time"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, "Page not found", http.StatusNotFound)
}

func ResponseError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	responseBody := response.ErrorResponse{
		Code:  code,
		Message: msg,
	}
	_ = json.NewEncoder(w).Encode(responseBody)
	fmt.Printf("%s : %s\n", time.Now(), msg)
}

func ResponseJSON(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		ResponseError(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func ResponseFile(w http.ResponseWriter,file []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(file)
	if err != nil {
		ResponseError(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}