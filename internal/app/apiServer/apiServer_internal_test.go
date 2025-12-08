package apiserver

import (
	"RestApi/internal/app/model"
	teststore "RestApi/internal/app/store/testStore"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {

	s := NewServer(teststore.NewStore(), sessions.NewCookieStore([]byte("secret")))

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
}

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.NewStore()
	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name               string
		cookieValue        map[interface{}]interface{}
		expectedStatusCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "not authenticated",
			cookieValue:        nil,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	var secretKey = "secret"
	s := NewServer(store, sessions.NewCookieStore([]byte(secretKey)))
	sc := securecookie.New([]byte(secretKey), nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenicateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}

func TestServer_HandleSessionCreate(t *testing.T) {
	u := model.TestUser(t)
	store := teststore.NewStore()
	store.User().Create(u)
	s := NewServer(store, sessions.NewCookieStore([]byte("secret")))
	testCases := []struct {
		name               string
		payload            interface{}
		expectedStatusCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid payload",
			payload:            "invalid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "",
				"password": u.Password,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}
