package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const GOOGLE_BOOKS_API_KEY = "GOOGLE_BOOKS_API_KEY"
const BOOK_TITLE = "book_title"
const AUTHOR_NAME = "author_name"
const START_PATTERN_IMAGE_SEARCH = "https://images-na.ssl-images-amazon.com/images"
const START_PATTERN_GOODREADS_SEARCH = "https://www.goodreads.com/book/show/"
const END_PATTERN_GOODREADS_SEARCH = "&"
const END_PATTERN_IMAGE_SEARCH = "\""

func BookcoverSearch(w http.ResponseWriter, r *http.Request) {
  bookTitle := strings.ReplaceAll(r.URL.Query().Get(BOOK_TITLE), " ", "+")
  authorName := strings.ReplaceAll(r.URL.Query().Get(AUTHOR_NAME), " ", "+")
  q := bookTitle + "+" + authorName + "site:goodreads.com/book/show"
  query := "https://www.google.com/search?q=" + q + "&sourceid=chrome&ie=UTF-8"

  goodreadUrl := getUrl(getBody(w, query), START_PATTERN_GOODREADS_SEARCH, END_PATTERN_GOODREADS_SEARCH)
  imageUrl := getUrl(getBody(w, goodreadUrl), START_PATTERN_IMAGE_SEARCH, END_PATTERN_IMAGE_SEARCH)

  w.Write(BuildSuccessResponse(w, imageUrl))
}

func BookcoverByIsbn(w http.ResponseWriter, r *http.Request) {
  isbn := strings.ReplaceAll(r.PathValue("isbn"), "-", "")
  if(len(isbn) != 13) {
    w.Write(BuildErrorResponse(w, HttpException{ statusCode: http.StatusBadRequest, message: INVALID_ISBN }))
    return
  }

  query := "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn + "&key" + os.Getenv(GOOGLE_BOOKS_API_KEY)
  res := getBody(w, query)
  var googleBook GoogleBook = GoogleBook{} 
  if err := json.Unmarshal(res, &googleBook); err != nil {
    fmt.Println("Error while parsing JSON body")
    return
  }

  if googleBook.TotalItems == 0 {
    w.Write(BuildErrorResponse(w, HttpException{ 
      statusCode: http.StatusBadRequest,
      message: BOOKCOVER_NOT_FOUND,
    }))
    return
  }

  responseData := BuildSuccessResponse(w, googleBook.Items[0].VolumeInfo.ImageLinks.Thumbnail)
  w.Write(responseData)
}

func getBody(w http.ResponseWriter, url string) []byte {
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

  return body
}

func getUrl(data []byte, startPattern string, endPattern string) string {
  body := string(data)
  init := strings.Index(body, startPattern)
  end := strings.Index(body[init:], endPattern)

  return body[init:init+end]
}
