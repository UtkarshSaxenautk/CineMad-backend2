package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Content Types
const (
	ContentTypeJSON = "application/json"
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJsonHttpErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	errMessage := ErrorMessage{Code: statusCode, Message: err.Error()}

	// convert resp to json
	var response []byte
	response, err = json.Marshal(&errMessage)

	err = WriteHTTPResponse(w, response, ContentTypeJSON, statusCode)
	// failed to
	// write response maybe writer is closed
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func WriteJsonHttpResponse(w http.ResponseWriter, statusCode int, serviceResponse interface{}) {

	if serviceResponse != nil {
		response, err := json.Marshal(serviceResponse)
		if err != nil {
			WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, errors.New("Marshal error"))
			return
		}
		err = WriteHTTPResponse(w, response, ContentTypeJSON, statusCode)
		if err != nil {
			WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, errors.New("Write error"))
			return
		}
	} else {
		w.WriteHeader(statusCode)
	}
}

// WriteHTTPResponse would write the content type, default headers, status code, and body to the response.
// Returns error if failed. This function does not write an HTTP error so that there are no surprises
func WriteHTTPResponse(w http.ResponseWriter, data []byte, contentType string, statusCode int) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	// Write to response
	_, err := w.Write(data)
	return err
}
