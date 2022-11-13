import { Request, Response, NextFunction } from 'express';
import { getLinkGoogle, getLinkGoodreads } from '../helpers/bookcover';
import { BOOKCOVER_NOT_FOUND, INVALID_ISBN } from '../helpers/messages';
import Axios from 'axios';
import Url from 'url';
import HttpException from '../exceptions/HttpException';

type BookcoverResponse = {
  status: string,
  url: string
  bookTitle?: string,
  authorName?: string,
  isbn?: string,
}

export const getBookcoverUrl = async (req: Request, res: Response, next: NextFunction) => {
  try {
    const query = new URLSearchParams(Url.parse(req.url).query);
    const bookTitle = query.get('book_title');
    const authorName = query.get('author_name');
    const googleQuery = `${bookTitle} ${authorName} site:goodreads.com/book/show`;
    const googleResponse = await Axios.get(`https://www.google.com/search?q=${googleQuery}&sourceid=chrome&ie=UTF-8`);

    const goodreadsLink = getLinkGoogle(googleResponse.data);
    if (!goodreadsLink) {
      throw new HttpException(404, BOOKCOVER_NOT_FOUND);
    }

    const goodreadsResponse = await Axios.get(goodreadsLink);
    return res.json({ url: getLinkGoodreads(goodreadsResponse.data) });
  } catch (error) {
    next(error);
  }
}

export const getBookcoverFromISBN = async (req: Request, res: Response, next: NextFunction) => {
  try {
    const isbn = req.params.id.replace(/-/g, '');
    if (isbn.length !== 13) {
      throw new HttpException(400, INVALID_ISBN);
    }

    const response = await Axios.get(
      `https://www.googleapis.com/books/v1/volumes?q=isbn:${isbn}&key=${process.env.GOOGLE_BOOKS_API_KEY}`
    );
    if (!response.data.totalItems) {
      return res.json([]);
    }

    return res.json({ url: response.data.items[0].volumeInfo.imageLinks.thumbnail });
  } catch (error) {
    next(error);
  }
};