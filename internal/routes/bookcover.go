package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bookcover-api/internal/helpers"
)

const (
	GOOGLE_BOOKS_API_KEY                  = "GOOGLE_BOOKS_API_KEY"
	BOOK_TITLE                            = "book_title"
	AUTHOR_NAME                           = "author_name"
	START_PATTERN_GOODREADS_IMAGE_SEARCH  = "https://images-na.ssl-images-amazon.com/images"
	START_PATTERN_GOODREADS_GOOGLE_SEARCH = "https://www.goodreads.com/book/show/"
	START_PATTERN_AMAZON_GOOGLE_SEARCH    = "https://www.amazon.com/"
	START_PATTERN_AMAZON_IMAGE_SEARCH     = "https://m.media-amazon.com/images/"
	END_PATTERN_AMAZON_GOOGLE_SEARCH      = "&amp"
	END_PATTERN_GOODREADS_GOOGLE_SEARCH   = "&"
	END_PATTERN_IMAGE_SEARCH              = "\""
)

type HttpException struct {
	statusCode int
	message    string
}

func BuildSuccessResponse(w http.ResponseWriter, url string) []byte {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]string{"url": url})
	w.WriteHeader(200)

	return buffer.Bytes()
}

func BuildErrorResponse(w http.ResponseWriter, ex HttpException) []byte {
	data, err := json.Marshal(map[string]string{"error": ex.message})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ex)
	w.WriteHeader(ex.statusCode)
	return data
}

func BookcoverSearch(w http.ResponseWriter, r *http.Request) {
	bookTitle := r.URL.Query().Get(BOOK_TITLE)
	authorName := r.URL.Query().Get(AUTHOR_NAME)
	if bookTitle == "" || authorName == "" {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    helpers.MANDATORY_PARAMS_MISSING,
		}))
		return
	}
	bookTitle = strings.ReplaceAll(bookTitle, " ", "+")
	authorName = strings.ReplaceAll(authorName, " ", "+")
	q := bookTitle + "+" + authorName + "+site:goodreads.com/book/show"
	query := "https://www.google.com/search?q=" + q + "&sourceid=chrome&ie=UTF-8"

	body, err := getBody(query)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
		}))
		return
	}

	goodreadUrl, err := getUrl(body, START_PATTERN_GOODREADS_GOOGLE_SEARCH, END_PATTERN_GOODREADS_GOOGLE_SEARCH)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
		}))
		return
	}

	body, err = getBody(goodreadUrl)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
		}))
		return
	}

	imageUrl, err := getUrl(body, START_PATTERN_GOODREADS_IMAGE_SEARCH, END_PATTERN_IMAGE_SEARCH)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
		}))
		return
	}

	w.Write(BuildSuccessResponse(w, imageUrl))
}

func BookcoverByIsbn(w http.ResponseWriter, r *http.Request) {
	isbn := strings.ReplaceAll(r.PathValue("isbn"), "-", "")
	if len(isbn) != 13 {
		w.Write(BuildErrorResponse(w, HttpException{statusCode: http.StatusBadRequest, message: helpers.INVALID_ISBN}))
		return
	}

	query := "https://www.goodreads.com/search?utf8=✓&query=" + isbn
	body, err := getBody(query)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
		}))
		return
	}

	imageUrl, err := getUrl(body, START_PATTERN_GOODREADS_IMAGE_SEARCH, END_PATTERN_IMAGE_SEARCH)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
		}))
		return
	}

	w.Write(BuildSuccessResponse(w, imageUrl))
}

func getBody(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(helpers.BOOKCOVER_NOT_FOUND)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf(helpers.ERROR_READING_BODY)
	}

	return body, nil
}

func getUrl(data []byte, startPattern string, endPattern string) (string, error) {
	body := string(data)
	init := strings.Index(body, startPattern)
	if init == -1 {
		err := fmt.Errorf("Initial pattern not found")
		return "", err
	}
	end := strings.Index(body[init:], endPattern)
	if init == -1 {
		err := fmt.Errorf("End pattern not found")
		return "", err
	}
	return body[init : init+end], nil
}
