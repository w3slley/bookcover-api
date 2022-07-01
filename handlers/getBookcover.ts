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

type QueryParameter = {
    bookTitle?: string,
    authorName?: string,
}

export const getBookcoverUrl = (req, res) => {
    const requiredParams = [
        'bookTitle',
        'authorName',
    ];
    const missingParams = [];
    const params: QueryParameter = {};
    const query = new URLSearchParams(url.parse(req.url).query);
    requiredParams.forEach(param => {
        if (!query.has(param)) {
            missingParams.push(param);
        } else {
            params[param] = query.get(param);
        }
    });
    if (missingParams.length) {
        const missingParamsStringified = missingParams.reduce((prev, curr) => (prev + ', ' + curr));
        return res.end(JSON.stringify({status: 'failed', error: `Please insert the following required query parameters: ${missingParamsStringified}`}))
    }

    let { bookTitle } = params;
    let { authorName } = params;
    let googleQuery = `${bookTitle} ${authorName} site:goodreads.com/book/show`;
    let googleSearch = `https://www.google.com/search?q=${googleQuery}&sourceid=chrome&ie=UTF-8`;
    axios.get(googleSearch)
    .then((googleResponse) => {
        const body = googleResponse.data;
        let goodreadsLink = getLinkGoogle(body);
        if (!goodreadsLink) {
            return res.status(404).send(JSON.stringify({status: 'failed', error: 'Bookcover was not found.'}))
        }
        axios.get(goodreadsLink)
        .then((goodreadsResponse)=>{
            const body = goodreadsResponse.data;
            let bookCoverLink = getLinkGoodreads(body);
            let bookcoverResponse: BookcoverResponse = {
                status: 'success',
                bookTitle: bookTitle,
                authorName: authorName,
                bookcoverUrl: bookCoverLink
            }
            res.end(JSON.stringify(bookcoverResponse));
        })
        .catch((e: any) => {
            res.status(500).send(JSON.stringify({status: 'failed', error: e.message}));
        });
    })
    .catch((e: any) => {
        res.status(500).send(JSON.stringify({status: 'failed', error: e.message}));
    });
}