package main

import (
	"github.com/adelowo/reblog/handler"
	m "github.com/adelowo/reblog/middleware"
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/utils"
	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"log"
	"net/http"
	"time"
)

const DATABASE_NAME = "reblog.db"

func main() {

	db := models.MustNewDB(DATABASE_NAME)

	jwtGenerator := utils.NewJWTGenerator()

	h := &handler.Handler{DB: db, JWT: jwtGenerator}

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Heartbeat("/pingoflife"))
	router.Use(middleware.CloseNotify)
	router.Use(middleware.Timeout(time.Second * 60))

	router.Group(func(r chi.Router) {
		r.Use(m.Guest)
		r.Post("/login", handler.PostLogin(h))
		r.Post("/signup/:token", handler.PostSignUp(h))

	})

	router.Group(func(r chi.Router) {

		r.Route("/reblog", func(ro chi.Router) {

			ro.Use(jwtGenerator.Verifier)
			ro.Use(jwtauth.Authenticator)

			ro.Route("/collaborator", func(roo chi.Router) {

				roo.Use(m.Admin)

				roo.Post("/create", handler.CreateCollaborator(h))
			})
		})

	})

	log.Println("Starting app")
	http.ListenAndServe(":3000", router)
}
