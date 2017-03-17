package handler

import (
	"encoding/json"
	"github.com/pressly/chi/render"
	"net/http"
	//	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/utils"
)

//This is used to create a token to be sent to the user
//After which the user would be authenticated with the token.
func CreateCollaborator(h *Handler) func(w http.ResponseWriter, r *http.Request) {

	type d struct {
		Email string `json:"email"`
	}

	type res struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Email string `json:"email"`
		} `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data d

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, &res{false, "Bad request", struct {
				Email string `json:"email"`
			}{"Invalid Email"}})
			return
		}

		if !utils.IsEmail(data.Email) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Please provide a valid email address", struct {
				Email string `json:"email"`
			}{"Please provide a valid email address"}})
			return
		}

		//check if the user already exists as a user

		if _, err := h.DB.FindByEmail(data.Email); err == nil {
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, &res{false, "Collaborator exists", struct {
				Email string `json:"email"`
			}{"Email already identifies a collaborator"}})
			return
		}

		if err := h.DB.CreateCollaborator(data.Email); err == nil {
			w.WriteHeader(http.StatusOK)

			//				defer sendEmailHere()
			render.JSON(w, r, &res{true, "A email has been sent to the collaborator", struct {
				Email string `json:"email"`
			}{""}})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, &res{false, "An error occured while we tried adding a new collaborator", struct {
			Email string `json:"email"`
		}{""}})
	}
}

func PostSignUp(h *Handler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}
