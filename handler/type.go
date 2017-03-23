package handler

import (
	"github.com/adelowo/reblog/models"
	"github.com/adelowo/reblog/utils"
)

type Handler struct {
	DB   models.DataStore
	JWT  *utils.JWTTokenGenerator
	Slug utils.Slug
}
