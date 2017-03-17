package handler

import (
	"bytes"
	"errors"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidEmailAddress(t *testing.T) {
	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	data := []byte(`{"email" : "adelowo"}`)

	req, err := http.NewRequest("POST", "/reblog/collaborator/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreateCollaborator(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expected := string(`{"status":false,"message":"Please provide a valid email address","data":{"email":"Please provide a valid email address"}}`)

	assert.JSONEq(t, expected, rr.Body.String(), "The response body differs")
}

func TestCannotAddACollaboratorTwice(t *testing.T) {

	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	db.On("FindByEmail", "me@lanre.me").
		Return(models.User{ID: 1, Moniker: "adelowo", Type: 0}, nil)

	data := []byte(`{"email" : "me@lanre.me"}`)

	req, err := http.NewRequest("POST", "/reblog/collaborator/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreateCollaborator(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expected := string(`{"status":false,"message":"Collaborator exists","data":{"email":"Email already identifies a collaborator"}}`)

	assert.JSONEq(t, expected, rr.Body.String(), "The response body differs")

}

func TestAnErrorOccurredWhileAddingACollaborator(t *testing.T) {
	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	db.On("FindByEmail", "me@lanre.me").
		Return(models.User{}, errors.New("could not find user"))

	db.On("CreateCollaborator", "me@lanre.me").
		Return(errors.New("Something bad happened"))

	data := []byte(`{"email" : "me@lanre.me"}`)

	req, err := http.NewRequest("POST", "/reblog/collaborator/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreateCollaborator(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatal(status)
	}

	expected := string(`{"status":false,"message":"An error occured while we tried adding a new collaborator","data":{"email":""}}`)

	assert.JSONEq(t, expected, rr.Body.String(), "The response body differs")

}
