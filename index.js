const http = require('http')
const fs = require('fs')
const url = require('url')
const querystring = require('querystring')
const request = require('request')
const functions = require('./lib/functions.js')
const PORT = process.env.PORT || 8000

const server = http.createServer((req, res)=>{
    try{
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
                    let query = `${bookTitle.replace(' ', '+')} ${authorName.replace('', '+')}`
                    //making request to google to get book's goodreads page
                    request(`https://www.google.com/search?q=${query}+goodreads&sourceid=chrome&ie=UTF-8`, (err, response, body)=>{
                        if(err){
                            res.end(JSON.stringify({error: err}))
                        }
                        let goodreadsLink = functions.getLinkGoogle(body)
                        //Making request to goodreads to get the book cover image tag
                        request(goodreadsLink, (err, response, body)=>{
                            if(err)
                                res.end(JSON.stringify({error: err}))
                            let bookCoverLink = functions.getLinkGoodreads(body)
                            //instantiating new date object to get time finished
                            let d2 = new Date()
                            let timeEnd = d2.getTime()
                            let diff =  timeEnd-timeInit
                            //sending json response
                            res.end(JSON.stringify({status: 'success', bookTitle: bookTitle, authorName: authorName, delay:`${diff/1000} seconds` ,method: 'getBookCover', bookCoverUrl: bookCoverLink}))
                        })
                    })
                }
                else//if no query is inserted
                    res.end(JSON.stringify({status: 'failed', error: 'Please insert options for search.'}))

                break
            default:
                res.end(JSON.stringify({status: 'failed', error: 'Method not suported yet.'}))
        }
    }
    catch(e){
        res.end(JSON.stringify({status: 'failed', error: e.message}))
    }
})

server.listen(PORT, ()=>{
    console.log(`Server listening at port ${PORT}!`)
})
