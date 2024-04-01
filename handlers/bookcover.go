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
  query := "https://www.google.com/search?q="+q+"&sourceid=chrome&ie=UTF-8"

  response, err := http.Get(query)
  if err != nil {
    w.Write(BuildErrorResponse(w, HttpException{  
      statusCode: http.StatusBadRequest, message: BOOKCOVER_NOT_FOUND,
    }))
    return
  }

  body, err := io.ReadAll(response.Body)
  if err != nil {
    w.Write(BuildErrorResponse(w, HttpException{  
      statusCode: http.StatusBadRequest, message: ERROR_READING_BODY,
    }))
    return
  }

  html := string(body)

  // parse html for goodread link, then do the same thing again but for the image url
  goodreadUrl := GetUrl(html)
  fmt.Println(goodreadUrl)
}


func GetUrl(data string) string {
  //strings.Index(haystack, needle)
  return ""
}
