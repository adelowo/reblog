package middleware

import (
	"bytes"
	"github.com/adelowo/reblog/handler"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuest(t *testing.T) {
	db := &mocks.DataStore{}
	h := &handler.Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/login", nil)
	req.Header.Add("Authorization", "Bearer abc123")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Guest(http.HandlerFunc(handler.PostLogin(h))).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatal(status)
	} else {
		t.Log("One test passed")
	}

	req, err = http.NewRequest("POST", "/login", nil)

	rr = httptest.NewRecorder()

	Guest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You have been pawned"))
	})).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatal(status)
	} else {
		t.Log("Yet another test passed")
	}

	if !bytes.Equal([]byte("You have been pawned"), rr.Body.Bytes()) {
		t.Fail()
	} else {
		t.Log("Yup, another one passed")
	}
}
