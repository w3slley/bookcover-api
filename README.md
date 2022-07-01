# An API to retrieve bookcovers from the internet.

Since I discovered to be quite hard to find comprehensive APIs (tried both Goodreads' and Google's) to retrieve book cover images for a personal project, I decided to create this API which does exactly that. I plan to add more features on it as times goes on.

Hope this is helpful for anyone wanting to create web apps that involves dealing with book covers in one way or another.

## Documentation

Right now it only has one method.

### GET /bookcover

It accepts the following two required paramaters:

- bookTitle (string)
- authorName (string)

Example of an http request:

`https://bookcoverapi.herokuapp.com/getBookCover?bookTitle=The+Pale+Blue+Dot&authorName=Carl+Sagan`

Response:

```
{
  "status":"success",
  "bookTitle":"The Pale Blue Dot",
  "authorName":"Carl Sagan",
  "bookCoverUrl":"https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1500191671l/61663._SY475_.jpg"
}
```
