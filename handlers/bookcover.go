package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BOOK_TITLE = "book_title"
const AUTHOR_NAME = "author_name"
const GOODREAD_IMAGE_URL_PATTERN = "https://images-na.ssl-images-amazon.com/images"
const GOODREAD_URL = "https://www.goodreads.com/book/show/"

func Bookcover(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // TODO: move to middleware

	bookTitle := strings.ReplaceAll(r.URL.Query().Get(BOOK_TITLE), " ", "+")
	authorName := strings.ReplaceAll(r.URL.Query().Get(AUTHOR_NAME), " ", "+")
	q := bookTitle + "+" + authorName + "site:goodreads.com/book/show"
	query := "https://www.google.com/search?q=" + q + "&sourceid=chrome&ie=UTF-8"

	goodreadUrl := GetUrl(GetBody(w, query), GOODREAD_URL, "&")
	imageUrl := GetUrl(GetBody(w, goodreadUrl), GOODREAD_IMAGE_URL_PATTERN, "\"")
	fmt.Println(imageUrl)

}

func GetBody(w http.ResponseWriter, url string) string {
	response, err := http.Get(url)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest, message: BOOKCOVER_NOT_FOUND,
		}))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		w.Write(BuildErrorResponse(w, HttpException{
			statusCode: http.StatusBadRequest, message: ERROR_READING_BODY,
		}))
	}

	return string(body)
}

func GetUrl(data string, startPattern string, endPattern string) string {
	init := strings.Index(data, startPattern)
	end := strings.Index(data[init:], endPattern)

	return data[init : end+init]
}
