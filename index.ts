const PORT = process.env.PORT || 8000;
const app = require('express')();
import { getBookcoverUrl } from './handlers/getBookcover';

app.get('/bookcover', (req, res) => {
    res.setHeader('Content-Type', 'application/json');
    try{
        return getBookcoverUrl(req, res);
    }
    catch(error){
        res.send(500, JSON.stringify({status: 'failed', error: error.message}));
    }
});

app.get('*', (req, res) => {
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify({status: 'failed', error: 'Method not suported yet.'}));
});

app.listen(PORT, ()=>{
    console.log(`Server listening at port ${PORT}!`);
})
