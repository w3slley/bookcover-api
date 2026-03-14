package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"bookcover-api/internal/config"
	"bookcover-api/internal/service"
	"bookcover-api/pkg/response"
)

const (
	bookTitleParam  = "book_title"
	authorNameParam = "author_name"
	isbnParam       = "isbn"
	imageSizeParam  = "image_size"
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
	isbn := r.URL.Query().Get(isbnParam)
	bookTitle := r.URL.Query().Get(bookTitleParam)
	authorName := r.URL.Query().Get(authorNameParam)
	imageSize := r.URL.Query().Get(imageSizeParam)

	if isbn != "" && (bookTitle != "" || authorName != "") {
		w.Write(response.Error(w, http.StatusBadRequest, config.ConflictingParams))
		return
	}

	if isbn != "" {
		h.searchByISBN(w, isbn, imageSize)
		return
	}

	if bookTitle == "" || authorName == "" {
		w.Write(response.Error(w, http.StatusBadRequest, config.MandidatoryParamsMissing))
		return
	}

	imageURL, err := h.service.GetByTitleAuthor(bookTitle, authorName, imageSize)
	if err != nil {
		w.Write(response.Error(w, http.StatusNotFound, err.Error()))
		return
	}

	w.Write(response.Success(w, imageURL))
}

func (h *BookcoverHandler) searchByISBN(w http.ResponseWriter, isbn, imageSize string) {
	isbn = strings.ReplaceAll(isbn, "-", "")

	if len(isbn) != 13 {
		w.Write(response.Error(w, http.StatusBadRequest, config.InvalidISBN))
		return
	}

	imageURL, err := h.service.GetByISBN(isbn, imageSize)
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
	imageSize := r.URL.Query().Get(imageSizeParam)

	if len(isbn) != 13 {
		w.Write(response.Error(w, http.StatusBadRequest, config.InvalidISBN))
		return
	}

	imageURL, err := h.service.GetByISBN(isbn, imageSize)
	if err != nil {
		w.Write(response.Error(w, http.StatusNotFound, err.Error()))
		return
	}

	w.Write(response.Success(w, imageURL))
}

func CacheStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := service.GetMetricsStats()

		response := map[string]interface{}{
			"total_requests":   stats.TotalRequests,
			"cache_hits":       stats.CacheHits,
			"cache_misses":     stats.CacheMisses,
			"new_books_cached": stats.NewBooksCached,
			"scraping_errors":  stats.ScrapingErrors,
			"hit_ratio":        stats.HitRatio(),
			"miss_ratio":       stats.MissRatio(),
			"new_book_ratio":   stats.NewBookRatio(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
