package apiserver

import (
	teststore "RestApi/internal/app/store/testStore"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {

	s := NewServer(teststore.NewStore())

	testCases := []struct {
		name               string
		payload            interface{}
		expectedStatusCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: map[string]string{
				"password": "password",
			},
			expectedStatusCode: 422,
		},
		{
			name: "invalid Params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
	// rec := httptest.NewRecorder()
	// req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	// s.ServeHTTP(rec, req)
	// assert.Equal(t, rec.Code, http.StatusOK)
}
