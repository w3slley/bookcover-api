package handler

import (
	"net/http"
	"strings"

	"bookcover-api/internal/config"
	"bookcover-api/internal/service"
	"bookcover-api/pkg/response"
)

const (
	bookTitleParam  = "book_title"
	authorNameParam = "author_name"
)

type BookcoverHandler struct {
	service service.BookcoverService
}

func NewBookcoverHandler(svc service.BookcoverService) *BookcoverHandler {
	return &BookcoverHandler{
		service: svc,
	}
}

func (h *BookcoverHandler) Search(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.URL.Query().Get(bookTitleParam)
	authorName := r.URL.Query().Get(authorNameParam)

	if bookTitle == "" || authorName == "" {
		w.Write(response.Error(w, http.StatusBadRequest, config.MandidatoryParamsMissing))
		return
	}

	imageURL, err := h.service.GetByTitleAuthor(bookTitle, authorName)
	if err != nil {
		w.Write(response.Error(w, http.StatusNotFound, err.Error()))
		return
	}

	w.Write(response.Success(w, imageURL))
}

func (h *BookcoverHandler) ByISBN(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	isbn := strings.TrimPrefix(path, "/bookcover/")
	isbn = strings.ReplaceAll(isbn, "-", "")

	if len(isbn) != 13 {
		w.Write(response.Error(w, http.StatusBadRequest, config.InvalidISBN))
		return
	}

	imageURL, err := h.service.GetByISBN(isbn)
	if err != nil {
		w.Write(response.Error(w, http.StatusNotFound, err.Error()))
		return
	}

	w.Write(response.Success(w, imageURL))
}
