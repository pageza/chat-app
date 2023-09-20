package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		payload        interface{}
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Test OK status",
			statusCode:     http.StatusOK,
			payload:        map[string]string{"message": "OK"},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "OK"},
		},
		{
			name:           "Test Bad Request status",
			statusCode:     http.StatusBadRequest,
			payload:        map[string]string{"error": "Bad Request"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Bad Request"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			SendJSONResponse(rr, tt.statusCode, tt.payload)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			expectedResponse, _ := json.Marshal(tt.expectedBody)
			if rr.Body.String() != string(expectedResponse)+"\n" { // Adding "\n" because json.NewEncoder adds a newline
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedResponse)+"\n")
			}
		})
	}
}
