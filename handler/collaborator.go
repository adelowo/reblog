package handler

import (
	"encoding/json"
	"github.com/adelowo/gotils/bag"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/utils"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"net/http"
	"time"
)

const tokenTTL = 20 * time.Minute

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

func DeleteCollaborator(h *Handler) func(w http.ResponseWriter, r *http.Request) {
	type d struct {
		Email string `json:"email"`
	}

	type res struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var data d

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {

			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, &res{false, "An error occured while we tried deleting the collaborator"})
			return
		}

		if user, err := h.DB.FindByEmail(data.Email); err == nil {
			defer h.DB.DeleteUser(user)

			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, &res{true, "User was successfully deleted"})

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &res{false, "Could not delete non-existent user"})

	}
}

func PostSignUp(h *Handler) func(w http.ResponseWriter, r *http.Request) {

	type d struct {
		Moniker  string `json:"moniker"`
		Name     string `json:"full_name`
		Password string `json:"password"`
	}

	type errorMessages struct {
		//just embed struct d ?
		Moniker  string `json:"moniker"`
		Name     string `json:"full_name"`
		Password string `json:"password"`
	}

	type res struct {
		Status  bool          `json:"status"`
		Message string        `json:"message"`
		Errors  errorMessages `json:"errors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var data d
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "No input was provided",
				errorMessages{"Please provide your moniker", "Please provide your name", "Please provide your password"}})

			return
		}

		token := chi.URLParam(r, "token")

		collaborator, err := h.DB.FindCollaboratorByToken(token)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if time.Now().Sub(collaborator.CreatedAt) > tokenTTL {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, &res{false, "Token is expired, Please contact the admin to resend a new token", errorMessages{}})
			return
		}

		bag := bag.NewValidatorErrorBag()

		if len(data.Moniker) < 4 {
			bag.Add("moniker", "Your moniker should not be lesser than 4 characters")
		}

		if len(data.Name) < 6 {
			bag.Add("name", "Your name should not be lesser than 6 characters. E.g Lanre Adelowo")
		}

		if len(data.Password) < 10 {
			bag.Add("password", "Your password should have a length greater than 10")
		}

		if bag.Count() != 0 {
			w.WriteHeader(http.StatusBadRequest)
			monikerErr, _ := bag.Get("moniker")
			nameErr, _ := bag.Get("name")
			passwordErr, _ := bag.Get("password")

			render.JSON(w, r, &res{false, "Validation failed",
				errorMessages{monikerErr, nameErr, passwordErr}})
			return
		}

		//ALl went successfully, we can add the user as a collaborator now

		err = h.DB.CreateUser(&models.User{Moniker: data.Moniker, Email: collaborator.Email, Name: data.Name, Password: data.Password})

		if err == nil {
			defer h.DB.DeleteCollaborator(collaborator)

			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, &res{true, "You have been added as a contributor to Reblog. Please login in other to get started", errorMessages{}})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, &res{false, "An error occured while we tried adding you as a collaborator to Reblog. Please try again", errorMessages{}})

	}
}
