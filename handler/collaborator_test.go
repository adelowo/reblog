package handler

import (
	"bytes"
	"errors"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"github.com/pressly/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

//Tests for the signup process here

func TestAnInvalidTokenCannotBeUsedToSignUpAsAWriter(t *testing.T) {
	data := []byte(`{"name" : "Lanre Adelowo", "moniker" : "hades", "password" : "yetanotherbadpassword"}`)

	db := new(mocks.DataStore)

	token := "invalidtokenhere"

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/signup/"+token, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	db.On("FindCollaboratorByToken", token).
		Return(models.Collaborator{}, errors.New("Collaborator not found"))

	r.Post("/signup/:token", PostSignUp(h))

	r.ServeHTTP(rr, req)

	//http.HandlerFunc(PostSignUp(h)).
	//	ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatal(status)
	}

	expectedText := "Not Found\n"

	assert.Equal(t, expectedText, rr.Body.String())
}

func TestAnExpiredTokenCannotBeUsedToSignUpAsAUser(t *testing.T) {

	data := []byte(`{"name" : "Lanre Adelowo", "moniker" : "hades", "password" : "yetanotherbadpassword"}`)

	db := new(mocks.DataStore)

	token := "expiredtoken"

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/signup/"+token, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	db.On("FindCollaboratorByToken", token).
		Return(models.Collaborator{2, token, "me@lanre.com", time.Now().Add(-21 * time.Minute)}, nil)

	r.Post("/signup/:token", PostSignUp(h))

	r.ServeHTTP(rr, req)

	//http.HandlerFunc(PostSignUp(h)).
	//	ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expectedText := string(`{"status" : false, "message" : "Token is expired, Please contact the admin to resend a new token", "errors":{"moniker":"", "full_name":"", "password":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())

}

func TestCannotSignUserUpWithInvalidData(t *testing.T) {

	data := []byte(`{"name" : "Lan", "moniker" : "ha", "password" : "y"}`)

	db := new(mocks.DataStore)

	token := "token"

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/signup/"+token, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	db.On("FindCollaboratorByToken", token).
		Return(models.Collaborator{2, token, "me@lanre.com", time.Now().Add(15 * time.Minute)}, nil)

	r.Post("/signup/:token", PostSignUp(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expectedText := string(`{"status" : false, "message" : "Validation failed", "errors":{"moniker":"Your moniker should not be lesser than 4 characters", "full_name":"Your name should not be lesser than 6 characters. E.g Lanre Adelowo", "password":"Your password should have a length greater than 10"}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())

}

func TestCollaboratorIsSuccessfullyCreated(t *testing.T) {

	data := []byte(`{"name" : "Lanre Adelowo", "moniker" : "hades", "password" : "yetanotherbadpassword"}`)

	db := new(mocks.DataStore)

	token := "token"

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/signup/"+token, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	c := models.Collaborator{2, token, "me@lanre.com", time.Now().Add(15 * time.Minute)}

	db.On("DeleteCollaborator", c).
		Return(nil)

	db.On("FindCollaboratorByToken", token).
		Return(c, nil)

	db.On("CreateUser", &models.User{Moniker: "hades", Email: c.Email, Name: "Lanre Adelowo", Password: "yetanotherbadpassword"}).
		Return(nil)

	r.Post("/signup/:token", PostSignUp(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatal(status)
	}

	expectedText := string(`{"status" : true, "message" : "You have been added as a contributor to Reblog. Please login in other to get started", "errors":{"moniker":"", "full_name":"", "password":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())

}

func TestAnUnknownErrorOccurredWhileCollaboratorTriedSigningUp(t *testing.T) {

	data := []byte(`{"name" : "Lanre Adelowo", "moniker" : "hades", "password" : "yetanotherbadpassword"}`)

	db := new(mocks.DataStore)

	token := "token"

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("POST", "/signup/"+token, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	c := models.Collaborator{2, token, "me@lanre.com", time.Now().Add(15 * time.Minute)}

	db.On("DeleteCollaborator", c).
		Return(nil)

	db.On("FindCollaboratorByToken", token).
		Return(c, nil)

	db.On("CreateUser", &models.User{Moniker: "hades", Email: c.Email, Name: "Lanre Adelowo", Password: "yetanotherbadpassword"}).
		Return(errors.New("An error occured"))

	r.Post("/signup/:token", PostSignUp(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expectedText := string(`{"status" : false, "message" : "An error occured while we tried adding you as a collaborator to Reblog. Please try again", "errors":{"moniker":"", "full_name":"", "password":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())

}

func TestCanDeleteACollaborator(t *testing.T) {

	data := []byte(`{"email" : "assholeuser@app.live"}`)

	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	u := models.User{Moniker: "asshole", Type: 0, Email: "assholeuser@app.live"}

	db.On("FindByEmail", "assholeuser@app.live").Return(u, nil)

	db.On("DeleteUser", u).Return(nil)

	req, err := http.NewRequest("POST", "/reblog/collaborator/delete", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(DeleteCollaborator(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatal(status)
	}

	expected := string(`{"status":true,"message":"User was successfully deleted"}`)

	assert.JSONEq(t, expected, rr.Body.String(), "The response body differs")

}

func TestCannotDeleteANonExistentUser(t *testing.T) {

	data := []byte(`{"email" : "unknownuser@app.live"}`)

	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	db.On("FindByEmail", "unknownuser@app.live").Return(models.User{}, errors.New("User doesn't exist"))

	req, err := http.NewRequest("POST", "/reblog/collaborator/delete", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(DeleteCollaborator(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expected := string(`{"status":false,"message":"Could not delete non-existent user"}`)

	assert.JSONEq(t, expected, rr.Body.String(), "The response body differs")

}
