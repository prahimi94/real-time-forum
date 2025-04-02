package utils

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	HttpStatus int         `json:"httpStatus,omitempty"`
}

func ReturnJson(w http.ResponseWriter, outputData Result) {
	w.Header().Set("Content-Type", "application/json")
	if outputData.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(outputData.HttpStatus)
	}

	if err := json.NewEncoder(w).Encode(outputData); err != nil {
		errorResponse := Result{
			Success: false,
			Message: "Failed to encode JSON",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(errorResponse)
	}
}
