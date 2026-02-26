# An API to retrieve bookcovers from the internet.

This is a simple API that fetches book cover images from Goodreads. You can search for covers using either a book's title and author, or its ISBN number. It returns a direct URL to the cover image that you can use in your applications.

## Documentation

### GET /bookcover

Search for a book cover by title/author or by ISBN-13.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `book_title` | string | Yes* | The title of the book |
| `author_name` | string | Yes* | The name of the book's author |
| `isbn` | string | Yes* | The ISBN-13 number of the book |
| `image_size` | string | No | Size of the cover image: `small`, `medium`, `large` (default) |

\* Provide either `book_title` + `author_name`, or `isbn`.

**Example Requests:**
```bash
# Search by title and author
curl -X GET "https://bookcover.longitood.com/bookcover?book_title=The+Pale+Blue+Dot&author_name=Carl+Sagan"

# Search by title and author with image size
curl -X GET "https://bookcover.longitood.com/bookcover?book_title=The+Pale+Blue+Dot&author_name=Carl+Sagan&image_size=small"

# Search by ISBN
curl -X GET "https://bookcover.longitood.com/bookcover?isbn=978-0345376596"
```

**Example Response:**
```json
{
  "url": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1388620656i/55030.jpg"
}
```

### GET /bookcover/:isbn (deprecated)

The path-based ISBN lookup is still supported for backwards compatibility.

**Example Request:**
```bash
curl -X GET "https://bookcover.longitood.com/bookcover/978-0345376596"
```

## How It Works

The API fetches book cover images from Goodreads using two different approaches:

1. **Search by Title and Author**
   - Takes the book title and author name as input
   - Searches Goodreads and finds the matching book
   - Extracts the high-quality cover image URL
   - Caches the result for faster future requests

2. **Search by ISBN-13**
   - Accepts a 13-digit ISBN number
   - Performs a direct lookup on Goodreads
   - Returns the book cover URL
   - Also caches successful results


The API provides clear error messages in JSON format:
- 400 Bad Request: Missing parameters or invalid ISBN
- 404 Not Found: No matching book cover found
- 429 Too Many Requests: Rate limiting quotas were met
- All responses include appropriate CORS headers

