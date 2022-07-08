const PORT = process.env.PORT || 8000;
const app = require('express')();
require('dotenv').config();

app.use((req, res, next) => {
    res.setHeader('Content-Type', 'application/json');
    next();
})

app.use('/bookcover', require('./routes/bookcover'));

app.get('*', (req, res) => {
    res.setHeader('Content-Type', 'application/json');
    res.status(400).json({status: 'failed', error: 'Method not suported yet.'});
});

app.listen(PORT, ()=>{
    console.log(`Server listening at port ${PORT}!`);
})
