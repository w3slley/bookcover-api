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
    const query = new URLSearchParams(url.parse(req.url).query);
    if(!query.has('bookTitle') && !query.has('authorName')){
        return res.end(JSON.stringify({status: 'failed', error: 'Please insert options for search.'}))
    }

    let bookTitle = query.get('bookTitle');
    let authorName = query.get('authorName');
    //making request to google to get book's goodreads page
    let googleQuery = `${bookTitle} ${authorName} site:goodreads.com/book/show`;
    let googleSearch = `https://www.google.com/search?q=${googleQuery}&sourceid=chrome&ie=UTF-8`;
    axios.get(googleSearch)
    .then((googleResponse) => {
        const body = googleResponse.data;
        let goodreadsLink = getLinkGoogle(body);
        if(!goodreadsLink) {
            return res.status(404).send(JSON.stringify({status: 'failed', error: 'Bookcover was not found.'}))
        }

        //Making request to goodreads to get the book cover image tag
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