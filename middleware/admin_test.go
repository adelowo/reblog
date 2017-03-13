package middleware

import (
	"bytes"
	"github.com/adelowo/reblog/handler"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdmin(t *testing.T) {
	db := &mocks.DataStore{}

	user := models.User{ID: 1, Moniker: "hades", Type: 1}

	h := &handler.Handler{DB: db, JWT: utils.NewJWTGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = user.ID
	claims["moniker"] = user.Moniker
	claims["type"] = user.Type

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/reblog/collaborator/create", nil)

	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.JWT.Verifier(Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Created a new user"))
	}))).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatal(status)
	}

	if !bytes.Equal([]byte("Created a new user"), rr.Body.Bytes()) {
		t.Fatal("Response body differ")
	} else {
		t.Log("Response body were the same")
	}

	testPreventsCollaboratorsFromAccessingThisEndpoint(h, t)
}

func testPreventsCollaboratorsFromAccessingThisEndpoint(h *handler.Handler, t *testing.T) {
	user := models.User{ID: 4, Moniker: "alcheme", Type: 0}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = user.ID
	claims["moniker"] = user.Moniker
	claims["type"] = user.Type

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/reblog/collaborator/create", nil)

	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.JWT.Verifier(Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Created a new user"))
	}))).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatal(status)
	}

	expected := string(`{"status":false,"message":"You do not have permission to view this resource"}`)

	assert.JSONEq(t, expected, rr.Body.String(), "Expected json to be equal")

}
