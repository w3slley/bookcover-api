package scraper

type Scraper interface {
	FetchByTitleAuthor(bookTitle, authorName string) (string, error)
	FetchByISBN(isbn string) (string, error)
}
