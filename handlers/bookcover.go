package handlers

import (
	"io"
	"net/http"
	"strings"
)

const BOOK_TITLE = "book_title"
const AUTHOR_NAME = "author_name"
const START_PATTERN_IMAGE_SEARCH = "https://images-na.ssl-images-amazon.com/images"
const START_PATTERN_GOODREADS_SEARCH = "https://www.goodreads.com/book/show/"
const END_PATTERN_GOODREADS_SEARCH = "&"
const END_PATTERN_IMAGE_SEARCH = "\""

func Bookcover(w http.ResponseWriter, r *http.Request) {
  bookTitle := strings.ReplaceAll(r.URL.Query().Get(BOOK_TITLE), " ", "+")
  authorName := strings.ReplaceAll(r.URL.Query().Get(AUTHOR_NAME), " ", "+")
  q := bookTitle + "+" + authorName + "site:goodreads.com/book/show"
  query := "https://www.google.com/search?q=" + q + "&sourceid=chrome&ie=UTF-8"

  goodreadUrl := getUrl(getBody(w, query), START_PATTERN_GOODREADS_SEARCH, END_PATTERN_GOODREADS_SEARCH)
  imageUrl := getUrl(getBody(w, goodreadUrl), START_PATTERN_IMAGE_SEARCH, END_PATTERN_IMAGE_SEARCH)

  w.Write(BuildSuccessResponse(w, imageUrl))
}

func getBody(w http.ResponseWriter, url string) string {
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

func getUrl(data string, startPattern string, endPattern string) string {
  init := strings.Index(data, startPattern)
  end := strings.Index(data[init:], endPattern)

  return data[init:init+end]
}
