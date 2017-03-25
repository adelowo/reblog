package handler

import (
	"encoding/json"
	"github.com/adelowo/gotils/bag"
	"github.com/adelowo/gotils/hasher"
	"github.com/adelowo/reblog/utils"
	"github.com/pressly/chi/render"
	"net/http"
)

type errorMessages struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authError struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Errors  errorMessages `json:"errors"`
}

//Admin and collaborators login
func PostLogin(h *Handler) func(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var data login
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r,
				&authError{false, "Could not log you in",
					errorMessages{Email: "Please provide your email address", Password: "Please provide your password"}})
			return
		}

		errorBag := bag.NewValidatorErrorBag()

		if !utils.IsEmail(data.Email) {
			errorBag.Add("email", "Please provide a valid email address")
		}

		if len(data.Password) == 0 {
			errorBag.Add("password", "Please provide a password")
		}

		if errorBag.Count() != 0 {

			e, _ := errorBag.Get("email")
			p, _ := errorBag.Get("password")

			w.WriteHeader(http.StatusUnauthorized)

			render.JSON(w, r, &authError{false, "Authentication failed",
				errorMessages{Email: e, Password: p}})

			return
		}

		//Check if the user exists in the database

		user, err := h.DB.FindByEmail(data.Email)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			render.JSON(w, r,
				&authError{false, "Authentication failed",
					errorMessages{Email: "Invalid username/password"}})

			return
		}

		if valid := hasher.NewBcryptHasher(12).Verify(user.Password, data.Password); valid {
			claims := make(map[string]interface{}, 4)

			claims["userID"] = user.ID
			claims["moniker"] = user.Moniker
			claims["type"] = user.Type

			h.JWT.Claims(claims)

			token, err := h.JWT.Generate()

			if err != nil {

				w.Header().Set("Content-Type", "application/json; charset=utf-8")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				return
			}

			w.WriteHeader(200)

			type d struct {
				Status  bool   `json:"status"`
				Message string `json:"message"`
				Data    struct {
					Token string `json:"token"`
				} `json:"data"`
			}

			render.JSON(w, r, &d{true, "You have been authenticated", struct {
				Token string `json:"token"`
			}{token}})

			return

		}

		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, &authError{false, "Authentication failed", errorMessages{Email: "Invalid email/password"}})
	}
}
