package middleware

import (
	"github.com/pressly/chi/render"
	"net/http"
	"strings"
)

//Middleware to prevent people who are already authenticated to use the /login route
//This simply checkd if the request contains a jwt token. If it contains, we'd consider the user as logged in,
//Else, we keep him/her out of the app
//All this is checking is the token presence.
//The real verification is in the verifier middleware attached to the API routes.
//If that fails, obviously the client should know it is time to get rid of the token and re-authenticate.
func Guest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer := r.Header.Get("Authorization")

		//BEARER is for consistency with the verifier middleware
		if len(bearer) > 7 && (strings.ToUpper(bearer[0:6]) == "BEARER" || bearer[0:6] == "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)

			d := make(map[string]interface{}, 4)
			d["status"] = false
			d["message"] = "You have been authenticated already."

			render.JSON(w, r, d)
			return
		}

		next.ServeHTTP(w, r)

	})
}
