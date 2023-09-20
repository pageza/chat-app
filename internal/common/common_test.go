// common_test.go

package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewAPIError tests the NewAPIError function
func TestNewAPIError(t *testing.T) {
	status := http.StatusBadRequest
	message := "Bad Request"

	apiErr := NewAPIError(status, message)

	assert.Equal(t, status, apiErr.Status, "Status should match")
	assert.Equal(t, message, apiErr.Message, "Message should match")
}

// TestRespondWithError tests the RespondWithError function
func TestRespondWithError(t *testing.T) {
	status := http.StatusNotFound
	message := "Not Found"

	apiErr := NewAPIError(status, message)

	_, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	RespondWithError(rr, apiErr)

	assert.Equal(t, status, rr.Code, "Status code should match")

	var response APIError
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, response.Status, "Response status should match")
	assert.Equal(t, message, response.Message, "Response message should match")
}
