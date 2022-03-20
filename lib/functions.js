const fs = require('fs')
const url = require('url')
const querystring = require('querystring')
const axios = require('axios')
const http = require('http')

function find(str, term, startsBy=0){
    if(str === undefined){
        return -1
    }
    let len = 0
    let pos = null
    for(let i=startsBy; i < str.length; i++){
        if(str[i] == term[len])
            len++
        else
            len=0
        if(len == term.length){
            pos = i + 1 - term.length //gets position i-term.length but has to add 1 given that startsBy has default value 0
            break
        }
    }
    if(pos != null) return pos

    return -1
}

function getLinkGoogle(data){
    if(data === undefined){
        return null;
    }
    let init = find(data, 'https://www.goodreads.com/book/show/')
    let final = find(data, "&", init+10)
    let linkGoogle = data.slice(init, final)
    return linkGoogle
}
function getLinkGoodreads(data){
    if(data === undefined){
        return null;
    }
    let init = find(data, '<img src="https://i.gr-assets.com/images/')
    let final = find(data, '"', init+10)
    let linkGoodreads = data.slice(init+10, final)
    return linkGoodreads
}

module.exports.getBookcoverUrl = (req, res) => {
    res.setHeader('Content-Type', 'application/json')
    //instantiating date object to measure time it took to get image
    let d1 = new Date()
    let timeInit = d1.getTime()
    let urlString = url.parse(req.url)
    switch(urlString.pathname){
        case '/getBookCover':
            if(querystring.parse(urlString.query)['bookTitle'] || querystring.parse(urlString.query)['authorName']){
                let bookTitle = querystring.parse(urlString.query)['bookTitle'] || ''
                let authorName = querystring.parse(urlString.query)['authorName'] || ''
                let query = `${bookTitle.replace(' ', '+')} ${authorName.replace(' ', '+')}`
                //making request to google to get book's goodreads page
                axios.get(`https://www.google.com/search?q=${query}+goodreads&sourceid=chrome&ie=UTF-8`)
                .then((response) => {
                    const body = response.data;
                    let goodreadsLink = getLinkGoogle(body)
                    //Making request to goodreads to get the book cover image tag
                    axios.get(goodreadsLink)
                    .then((response)=>{
                        const body = response.data;
                        let bookCoverLink = getLinkGoodreads(body)
                        //instantiating new date object to get time finished
                        let d2 = new Date()
                        let timeEnd = d2.getTime()
                        let diff =  timeEnd-timeInit
                        //sending json response
                        res.end(JSON.stringify({status: 'success', bookTitle: bookTitle, authorName: authorName, delay:`${diff/1000} seconds` ,method: 'getBookCover', bookCoverUrl: bookCoverLink}))
                    })
                    .catch((e)=>{
                        res.status = 500;
                        res.end(JSON.stringify({status: 'failed', error: e.message}));
                    });
                })
                .catch( (e) => {
                    res.status = 500;
                    res.end(JSON.stringify({status: 'failed', error: e.message}));
                });
            }
            else//if no query is inserted
                res.end(JSON.stringify({status: 'failed', error: 'Please insert options for search.'}))

            break
        default:
            res.end(JSON.stringify({status: 'failed', error: 'Method not suported yet.'}))
    }
}