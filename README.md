# An API to retrieve bookcovers from the web.

Since I found really hard to find comprehensive APIs (tried both Goodreads' and Google's) to find and retrieve book cover images for a personal project, I decided to create this API (using Node.JS) that does exactly that. Right now it gets a .jpg image from Goodreads server by doing searches on the web. I plan to add more features on it over time.

Hope this is helpful for anyone wanting to create web apps that involves dealing with book covers in one way or another.

## Documentation

Right now it only has one method.

### getBookCover

It accepts two paramaters:

- bookTitle
- authorName

Example of a http request:

`http://bookcoverapi.herokuapp.com/getBookCover?bookTitle=Cosmos&authorName=Carl+Sagan`

Response:

```
{
"status":"success",
"delay":"5.08 seconds",
"method":"getBookCover",
"bookCoverUrl":"https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1388620656l/55030.jpg"
}
```
