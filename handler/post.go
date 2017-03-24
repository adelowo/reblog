package handler

import (
	"encoding/json"
	"errors"
	"github.com/adelowo/gotils/bag"
	"github.com/adelowo/reblog/middleware"
	"github.com/adelowo/reblog/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"net/http"
	"strconv"
)

const (
	UNPUBLISHED = iota
	PUBLISHED
)

func CreatePost(h *Handler) func(w http.ResponseWriter, r *http.Request) {

	type d struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	type errorMessages struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	type res struct {
		Status  bool          `json:"status"`
		Message string        `json:"message"`
		Errors  errorMessages `json:"errors"`
	}

	//If the author of the post isn't the admin, mark the post as unpublished
	return func(w http.ResponseWriter, r *http.Request) {
		var data d

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Post could not be created", errorMessages{}})
			return
		}

		//Validate
		errorBag := bag.NewValidatorErrorBag()

		if len(data.Title) < 10 {
			errorBag.Add("title", "An article's title should be more than 10 characters")
		}

		if len(data.Content) < 100 {
			errorBag.Add("content", "The content of the article is too small. Should be at least 100 characters in length")
		}

		if errorBag.Count() != 0 {
			titleErr, _ := errorBag.Get("title")
			contentError, _ := errorBag.Get("content")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Post could not be created due to invalid data", errorMessages{titleErr, contentError}})
			return

		}

		_, err := h.DB.FindPostByTitle(data.Title)

		if err == nil {
			//post exists
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Could not create post as that would lead to duplicates", errorMessages{Title: "Post with title, " + data.Title + " already exists"}})
			return
		}

		userType, err := getUserType(r)

		if err != nil {
			//this shouldn't happen though, just paranoia
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		var status int

		if userType == middleware.ADMIN {
			status = PUBLISHED
		} else {
			status = UNPUBLISHED
		}

		slug := h.Slug.Generate(data.Title)

		p := models.Post{Title: data.Title, Content: data.Content, Slug: slug, Status: status}

		if err = h.DB.CreatePost(p, userType); err == nil {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, &res{true, "Post was successfully created", errorMessages{}})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, &res{false, "An error occurred while trying to create the post", errorMessages{}})
	}
}

func DeletePost(h *Handler) func(w http.ResponseWriter, r *http.Request) {

	type res struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Errors  struct {
			PostID string `json:"post_id"`
		} `json:"errors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Invalid request", struct {
				PostID string `json:"post_id"`
			}{"Invalid post id"}})
			return
		}

		p, err := h.DB.FindPostByID(id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Post does not exist", struct {
				PostID string `json:"post_id"`
			}{"Post with the specified id could not be found"}})
			return
		}

		if err = h.DB.DeletePost(p); err == nil {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, &res{true, "Post was deleted", struct {
				PostID string `json:"post_id"`
			}{""}})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, &res{false, "An error occurred while trying to delete post", struct {
			PostID string `json:"post_id"`
		}{}})
	}
}

func UnpublishPost(h *Handler) func (w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func getUserType(r *http.Request) (int, error) {

	ctx := r.Context()

	jwtToken, ok := ctx.Value("jwt").(*jwt.Token)

	if !ok || jwtToken == nil || !jwtToken.Valid {
		return 0, errors.New("Could not fetch  user's id")
	}

	claims := jwtToken.Claims

	return int(claims["type"].(float64)), nil
}
