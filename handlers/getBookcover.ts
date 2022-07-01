const url = require('url');
const axios = require('axios');
import { getLinkGoogle, getLinkGoodreads} from '../helpers/bookcover';

type BookcoverResponse = {
    status: string,
    bookcoverUrl: string
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
        return res.end(JSON.stringify({status: 'failed', error: `Please insert the following required query parameters: ${missingParamsStringified}`}))
    }
    const bookTitle = query.get('book_title');
    const authorName = query.get('author_name');
    const googleQuery = `${bookTitle} ${authorName} site:goodreads.com/book/show`;
    const googleSearch = `https://www.google.com/search?q=${googleQuery}&sourceid=chrome&ie=UTF-8`;
    axios.get(googleSearch)
    .then((googleResponse) => {
        const body = googleResponse.data;
        const goodreadsLink = getLinkGoogle(body);
        if (!goodreadsLink) {
            return res.status(404).send(JSON.stringify({status: 'failed', error: 'Bookcover was not found.'}))
        }
        axios.get(goodreadsLink)
        .then((goodreadsResponse)=>{
            const body = goodreadsResponse.data;
            const bookCoverLink = getLinkGoodreads(body);
            res.json({
                status: 'success',
                bookcoverUrl: bookCoverLink
            });
        })
        .catch((e: any) => {
            res.status(500).send(JSON.stringify({status: 'failed', error: e.message}));
        });
    })
    .catch((e: any) => {
        res.status(500).send(JSON.stringify({status: 'failed', error: e.message}));
    });
}