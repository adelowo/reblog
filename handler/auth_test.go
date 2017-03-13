package handler

import (
	"testing"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"net/http"
	"net/http/httptest"
	"bytes"
	"errors"
	"fmt"
)

var _ models.DataStore = &mocks.DataStore{}

func TestPostLogin(t *testing.T) {

	db := new(mocks.DataStore)

	db.On("FindByEmail","adelowo@me.com").
		Return(models.User{},errors.New("User does not exists"))


	h := &Handler{DB : db, JWT : utils.NewJWTGenerator()}

	testInvalidPostBody(h, t)
	testDataFailsValidation(h, t)
	testInvalidUser(h, t)
	testSuccess(t)


	db.AssertExpectations(t)
}

func testInvalidPostBody(h *Handler, t *testing.T) {

	//empty POST body
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{}`)))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(PostLogin(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatal(err)
	} else {
		t.Log("Test passed :")
	}
}

func testDataFailsValidation(h *Handler, t *testing.T) {

	data := []byte(`{"email" : "adelowo", "password" : "b"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(PostLogin(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatalf("Expected %d, got %d", http.StatusUnauthorized, status)
	} else {
		t.Log("Status code check passed")
	}
}

func testInvalidUser(h *Handler, t *testing.T) {
	data := []byte(`{"email" : "adelowo@me.com", "password" : "badpass"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(PostLogin(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatalf("Expected %d, got %d", http.StatusUnauthorized, status)
	} else {
		t.Log("Passing")
	}
}


func testSuccess(t *testing.T) {
	db := new(mocks.DataStore)

	db.On("FindByEmail","adelowo@me.com").
		Return(models.User{ID:1, Password:"$2a$12$Xc6ArM465UaZVW/bbZorSec/dgkSApoC0Ac7Zfi6MajZlSnerqMAW", Moniker:"adelowo", Type: 0}, nil)


	h := &Handler{DB : db, JWT : utils.NewJWTGenerator()}

	data := []byte(`{"email" : "adelowo@me.com", "password" : "badpassword"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(PostLogin(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		fmt.Println(rr.Body.String())
		t.Fatalf("Expected %d, got %d", http.StatusOK, status)
	} else {
		t.Log("Passing")
	}
}
