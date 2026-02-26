package service

type BookcoverService interface {
	GetByTitleAuthor(bookTitle, authorName, imageSize string) (string, error)
	GetByISBN(isbn, imageSize string) (string, error)
}
