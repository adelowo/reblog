package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi/render"
	"net/http"
)

const (
	COLLABORATOR = iota
	ADMIN
)

//This middleware protects some routes from collaborators (writers)
//Users are divided into two types ; Collaborators (writers) and admin
//We obviously do no want to give collaborators super cow powers
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		jwtToken, ok := ctx.Value("jwt").(*jwt.Token)

		if !ok || jwtToken == nil || !jwtToken.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		claims := jwtToken.Claims

		userType := int(claims["type"].(float64))

		if userType == ADMIN {
			next.ServeHTTP(w, r)
			return

		}

		sendFailureResponse(w, r)
		return
	})
}

func sendFailureResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)

	d := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{
		false,
		"You do not have permission to view this resource",
	}

	render.JSON(w, r, d)

	return

}
