
<div align="right">
  <details>
    <summary >🌐 Language</summary>
    <div>
      <div align="center">
        <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=en">English</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=zh-CN">简体中文</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=zh-TW">繁體中文</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=ja">日本語</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=ko">한국어</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=hi">हिन्दी</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=th">ไทย</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=fr">Français</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=de">Deutsch</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=es">Español</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=it">Italiano</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=ru">Русский</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=pt">Português</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=nl">Nederlands</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=pl">Polski</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=ar">العربية</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=fa">فارسی</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=tr">Türkçe</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=vi">Tiếng Việt</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=id">Bahasa Indonesia</a>
        | <a href="https://openaitx.github.io/view.html?user=w3slley&project=bookcover-api&lang=as">অসমীয়া</
      </div>
    </div>
  </details>
</div>

# An API to retrieve bookcovers from the internet.

This is a simple API that fetches book cover images from Goodreads. You can search for covers using either a book's title and author, or its ISBN number. It returns a direct URL to the cover image that you can use in your applications.

## Documentation

### GET /bookcover

Search for a book cover using the book title and author name.

**Required Parameters:**
- `book_title` (string): The title of the book
- `author_name` (string): The name of the book's author

**Example Request:**
```bash
curl -X GET "https://bookcover.longitood.com/bookcover?book_title=The+Pale+Blue+Dot&author_name=Carl+Sagan"
```

**Example Response:**
```json
{
  "url": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1388620656i/55030.jpg"
}
```

### GET /bookcover/:isbn

Search for a book cover using its ISBN-13 number.

**Example Request:**
```bash
curl -X GET "https://bookcover.longitood.com/bookcover/978-0345376596"
```

**Example Response:**
```json
{
  "url": "https://images-na.ssl-images-amazon.com/images/S/compressed.photo.goodreads.com/books/1500191671i/61663.jpg"
}
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
- All responses include appropriate CORS headers
