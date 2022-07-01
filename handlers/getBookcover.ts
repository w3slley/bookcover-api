const url = require('url');
const axios = require('axios');
import { getLinkGoogle, getLinkGoodreads} from '../helpers/getBookCover';

type BookcoverResponse = {
    status: string,
    url: string
    bookTitle?: string,
    authorName?: string,
    isbn?: string,
}

export const getBookcoverUrl = (req, res) => {
    const requiredParams = [
        'book_title',
        'author_name',
    ];
    const missingParams = [];
    const query = new URLSearchParams(url.parse(req.url).query);
    requiredParams.forEach(param => {
        if (!query.has(param)) {
            missingParams.push(param);
        }
    });
    if (missingParams.length) {
        const missingParamsStringified = missingParams.reduce((prev, curr) => (prev + ', ' + curr));
        return res.status(400).json({status: 'failed', error: `Please insert the following required query parameters: ${missingParamsStringified}`});
    }
    const bookTitle = query.get('book_title');
    const authorName = query.get('author_name');
    const googleQuery = `${bookTitle} ${authorName} site:goodreads.com/book/show`;
    axios.get(`https://www.google.com/search?q=${googleQuery}&sourceid=chrome&ie=UTF-8`)
    .then((googleResponse) => {
        const goodreadsLink = getLinkGoogle(googleResponse.data);
        if (!goodreadsLink) {
            return res.status(404).json({status: 'failed', error: 'Bookcover was not found.'});
        }
        axios.get(goodreadsLink)
        .then((goodreadsResponse)=>{
            res.json({
                status: 'success',
                url: getLinkGoodreads(goodreadsResponse.data)
            });
        })
        .catch((e: any) => {
            res.status(500).json({status: 'failed', error: e.message});
        });
    })
    .catch((e: any) => {
        res.status(500).json({status: 'failed', error: e.message});
    });
}