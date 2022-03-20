const http = require('http')
const functions = require('./lib/functions.js')
const PORT = process.env.PORT || 8000

const server = http.createServer((req, res)=>{
    functions.getBookcoverUrl(req, res);
})

server.listen(PORT, ()=>{
    console.log(`Server listening at port ${PORT}!`)
})
