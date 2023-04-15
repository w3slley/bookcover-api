# An API to retrieve bookcovers from the internet.

Since I discovered to be quite hard to find comprehensive APIs (tried both Goodreads' and Google's) to retrieve book cover images for a personal project, I decided to create this API which does exactly that. I plan to add more features on it as times goes on.

Hope this is helpful for anyone wanting to create web apps that involves dealing with book covers in one way or another.

## Documentation

### GET /bookcover

It accepts the following two required paramaters:

- book_title (string)
- author_name (string)

Example of an http request:

```
https://api.bookcover.longitood.com/bookcover?book_title=The+Pale+Blue+Dot&author_name=Carl+Sagan
```

Response:

```
{
  "url":"https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1500191671l/61663._SY475_.jpg"
}
```

### GET /bookcover/:isbn
Search books by ISBN-13.

Example of an http request:

```
https://api.bookcover.longitood.com/bookcover/978-0345376596
```

Response:

```
{
  "url":"http://books.google.com/books/content?id=peSODQAAQBAJ&printsec=frontcover&img=1&zoom=1&source=gbs_api"
}
```
