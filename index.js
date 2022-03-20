const http = require('http')
const functions = require('./lib/functions.js')
const PORT = process.env.PORT || 8000

const server = http.createServer((req, res)=>{
    try{
        functions.getBookcoverUrl(req, res);
    }
    catch(error){
        res.status = 500;
        res.end(JSON.stringify({status: 'failed', error: e.message}));
    }
})

server.listen(PORT, ()=>{
    console.log(`Server listening at port ${PORT}!`)
})
