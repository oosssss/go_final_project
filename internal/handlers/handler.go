package handlers

import (
	"go_final_project/internal/repository"
)

const DateFormat = "20060102"

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
