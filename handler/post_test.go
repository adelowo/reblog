package handler

import (
	"bytes"
	"github.com/adelowo/reblog/middleware"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/models/mocks"
	"github.com/adelowo/reblog/utils"
	"github.com/pkg/errors"
	"github.com/pressly/chi"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ models.DataStore = new(mocks.DataStore)

func TestCannotCreatePostDueToInvalidData(t *testing.T) {

	db := new(mocks.DataStore)

	data := []byte(`{"title" : "go", "content" : "go code"}`)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator(), Slug: utils.NewSlugGenerator()}

	req, err := http.NewRequest("POST", "/reblog/posts/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreatePost(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatal(status)
	}

	expectedText := string(`{"status": false, "message" : "Post could not be created due to invalid data", "errors" : {"title" : "An article's title should be more than 10 characters", "content" : "The content of the article is too small. Should be at least 100 characters in length"}}`)

	assert.JSONEq(t, expectedText, rr.Body.String(), "The response body differs")
}

func TestPostCouldNotBeCreatedBecauseItAlreadyExists(t *testing.T) {
	//s := strings.Repeat("Go is awesome", 90)

	data := []byte(`{"title" : "Go is awesome", "content": "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome"}`)

	req, err := http.NewRequest("POST", "/reblog/post/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	db := new(mocks.DataStore)

	db.On("FindPostByTitle", "Go is awesome").
		Once().
		Return(models.Post{}, nil)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator(), Slug: utils.NewSlugGenerator()}

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreatePost(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected %d. Got %d instead", http.StatusBadRequest, status)
	}

	expectedText := string(`{"status":false, "message":"Could not create post as that would lead to duplicates", "errors":{"title" : "Post with title, Go is awesome already exists","content" :""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())
}

func TestACollaboratorCanCreateAPost(t *testing.T) {
	data := []byte(`{"title" : "Go is awesome", "content": "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome"}`)

	req, err := http.NewRequest("POST", "/reblog/post/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	db := new(mocks.DataStore)

	db.On("FindPostByTitle", "Go is awesome").
		Once().
		Return(models.Post{}, errors.New("Post could not be found")) //like seriously ?

	p := models.Post{Title: "Go is awesome", Content: "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome", Slug: "Go-is-awesome", Status: UNPUBLISHED}

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator(), Slug: utils.NewSlugGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 51
	claims["moniker"] = "collab"
	claims["type"] = middleware.COLLABORATOR

	db.On("CreatePost", p, claims["type"]).
		Return(nil)

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreatePost(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected %d. Got %d instead", http.StatusBadRequest, status)
	}

	expectedText := string(`{"status" : true, "message" : "Post was successfully created", "errors":{"title":"", "content":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())
}

func TestAdminCanCreateAPost(t *testing.T) {
	data := []byte(`{"title" : "Go is awesome", "content": "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome"}`)

	req, err := http.NewRequest("POST", "/reblog/post/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	db := new(mocks.DataStore)

	db.On("FindPostByTitle", "Go is awesome").
		Once().
		Return(models.Post{}, errors.New("Post could not be found"))

	p := models.Post{Title: "Go is awesome", Content: "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome", Slug: "Go-is-awesome", Status: PUBLISHED}

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator(), Slug: utils.NewSlugGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 51
	claims["moniker"] = "collab"
	claims["type"] = middleware.ADMIN

	db.On("CreatePost", p, claims["type"]).
		Return(nil)

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreatePost(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected %d. Got %d instead", http.StatusBadRequest, status)
	}

	expectedText := string(`{"status" : true, "message" : "Post was successfully created", "errors":{"title":"", "content":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())
}

func TestAnErrorOccurredWhileTryingToCreateAPost(t *testing.T) {
	data := []byte(`{"title" : "Go is awesome", "content": "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome"}`)

	req, err := http.NewRequest("POST", "/reblog/post/create", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	db := new(mocks.DataStore)

	db.On("FindPostByTitle", "Go is awesome").
		Once().
		Return(models.Post{}, errors.New("Post could not be found"))

	p := models.Post{Title: "Go is awesome", Content: "Go is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesomeGo is awesome", Slug: "Go-is-awesome", Status: PUBLISHED}

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator(), Slug: utils.NewSlugGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 51
	claims["moniker"] = "collab"
	claims["type"] = middleware.ADMIN

	db.On("CreatePost", p, claims["type"]).
		Return(errors.New("Could not create post"))

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	http.HandlerFunc(CreatePost(h)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("Expected %d. Got %d instead", http.StatusInternalServerError, status)
	}

	expectedText := string(`{"status" : false, "message" : "An error occurred while trying to create the post", "errors":{"title":"", "content":""}}`)

	assert.JSONEq(t, expectedText, rr.Body.String())
}

func TestAnInvalidRequestCannotBeUsedToDeleteAPost(t *testing.T) {

	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("DELETE", "/reblog/posts/eighty", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Delete("/reblog/posts/:id", DeletePost(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected %d. Got %d", http.StatusBadRequest, status)
	}

	expectedText := string(`{"status" : false, "message" : "Invalid request", "errors": {"post_id" : "Invalid post id"}}`)

	assert.JSONEq(t, expectedText, rr.Body.String(), "THe response body differs")
}

func TestACollaboratorCannotDeleteAPost(t *testing.T) {

	db := new(mocks.DataStore)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("DELETE", "/reblog/posts/10", nil)

	if err != nil {
		t.Fatal(err)
	}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.COLLABORATOR

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Handle("/reblog/posts/:id", middleware.Admin(http.HandlerFunc(DeletePost(h))))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatalf("Expected %d. Got %d", http.StatusUnauthorized, status)
	}

	expected := string(`{"status" : false, "message" : "You do not have permission to view this resource"}`) //The admin middleware prevents the real handler from being called since the user is a collaborator.

	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAnAdminCanDeleteAPost(t *testing.T) {

	db := new(mocks.DataStore)

	p := models.Post{ID: 10, Status: PUBLISHED, Title: "Testing is key"}

	db.On("FindPostByID", 10).Once().Return(p, nil)

	db.On("DeletePost", p).Once().Return(nil)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("DELETE", "/reblog/posts/10", nil)

	if err != nil {
		t.Fatal(err)
	}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.ADMIN

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Delete("/reblog/posts/:id", DeletePost(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected %d. Got %d", http.StatusOK, status)
	}

	expected := string(`{"status" : true, "message" : "Post was deleted", "errors" : {"post_id" : ""}}`)

	assert.JSONEq(t, expected, rr.Body.String())

}

func TestNonExistentPostCannotBeDeleted(t *testing.T) {

	db := new(mocks.DataStore)

	p := models.Post{}

	db.On("FindPostByID", 10).Once().Return(p, errors.New("Post does not exist"))

	db.On("DeletePost", p).Once().Return(nil)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	req, err := http.NewRequest("DELETE", "/reblog/posts/10", nil)

	if err != nil {
		t.Fatal(err)
	}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.ADMIN

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	r := chi.NewRouter()

	r.Delete("/reblog/posts/:id", DeletePost(h))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected %d. Got %d", http.StatusBadRequest, status)
	}

	expected := string(`{"status" :false, "message" : "Post does not exist", "errors" : {"post_id" : "Post with the specified id could not be found"}}`)

	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminOnlyCanMarkAPostAsUnpublished(t *testing.T) {

	req, err := http.NewRequest("PUT", "/reblog/posts/10", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := &Handler{DB: new(mocks.DataStore), JWT: utils.NewJWTGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.COLLABORATOR

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	middleware.Admin(http.HandlerFunc(UnpublishPost(h))).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Fatalf("Expected %d. Got %d", http.StatusUnauthorized, status)
	}

	expected := string(`{"status" : false, "message" : "You do not have permission to view this resource"}`)

	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminCanMarkAPostAsUnpublished(t *testing.T) {

	req, err := http.NewRequest("PUT", "/reblog/posts/80", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := new(mocks.DataStore)

	p := models.Post{ID: 80}

	db.On("FindPostByID", 80).Once().Return(p, nil)

	db.On("UnpublishPost", p).Once().Return(nil)

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.ADMIN

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	r := chi.NewRouter()

	r.Handle("/reblog/posts/:id", middleware.Admin(http.HandlerFunc(UnpublishPost(h))))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected %d. Got %d", http.StatusOK, status)
	}

	expected := string(`{"status":true,"message":"Post was updated","errors":{"post_id":""}}`)

	assert.JSONEq(t, expected, rr.Body.String())
}

func TestCannotDeleteNonExistentPost(t *testing.T) {

	req, err := http.NewRequest("PUT", "/reblog/posts/80", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := new(mocks.DataStore)

	db.On("FindPostByID", 80).Once().Return(models.Post{}, errors.New("Post does not exist"))

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.ADMIN

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	r := chi.NewRouter()

	r.Handle("/reblog/posts/:id", middleware.Admin(http.HandlerFunc(UnpublishPost(h))))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("Expected %d. Got %d", http.StatusNotFound, status)
	}

	expected := string(`{"status":false,"message":"Post does not exist","errors":{"post_id":"Post with the specified id does not exist"}}`)

	assert.JSONEq(t, expected, rr.Body.String())

}

func TestAnErrorOccurredWhileTryingToUnpublishPost(t *testing.T) {

	req, err := http.NewRequest("PUT", "/reblog/posts/80", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	p := models.Post{ID: 80}

	db := new(mocks.DataStore)

	db.On("FindPostByID", 80).Once().Return(p, nil)

	db.On("UnpublishPost", p).Once().Return(errors.New("Could not unpublish Post"))

	h := &Handler{DB: db, JWT: utils.NewJWTGenerator()}

	claims := make(map[string]interface{}, 4)

	claims["userID"] = 15
	claims["moniker"] = "horus"
	claims["type"] = middleware.ADMIN

	h.JWT.Claims(claims)

	token, err := h.JWT.Generate()

	if err != nil {
		t.Fatal(err)
	}

	to, err := h.JWT.Decode(token)

	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()

	ctx = context.WithValue(ctx, "jwt", to)

	req = req.WithContext(ctx)

	r := chi.NewRouter()

	r.Handle("/reblog/posts/:id", middleware.Admin(http.HandlerFunc(UnpublishPost(h))))

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("Expected %d. Got %d", http.StatusInternalServerError, status)
	}

	expected := string(`{"status":false,"message":"An error occurred while trying to unpublish the post","errors":{"post_id":"Post could not be unpublished"}}`)

	assert.JSONEq(t, expected, rr.Body.String())

}
