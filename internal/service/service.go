package service

type BookcoverService interface {
	GetByTitleAuthor(bookTitle, authorName string) (string, error)
	GetByISBN(isbn string) (string, error)
}
