# An API to retrieve bookcovers from the web.

Since I found really hard to find useful apis (tried both the Goodreads' and Google's and failed) to find and retrieve book cover images for a personal project, I decided to create this API that does exactly that. Right now it gets a .jpg image from Goodreads server by doing searches on the web - that's the main reason it is somewhat slow (it takes between 3 to 9 seconds to get a response). But I plan to improve that along with the API as a whole.

Hope this is helpful for anyone wanting to create web apps that in one way or another utilize book covers.

## Documentation

Right now it only has one method.

#### getBookCover

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


