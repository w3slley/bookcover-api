const PORT = process.env.PORT || 2000;
import express, { Request, Response, NextFunction } from 'express';
const app = express();
require('dotenv').config();

app.use((req: Request, res: Response, next: NextFunction) => {
    res.setHeader('Content-Type', 'application/json');
    next();
})

app.use('/bookcover', require('./routes/bookcover'));

app.get('*', (req: Request, res: Response) => {
    res.setHeader('Content-Type', 'application/json');
    res.status(400).json({ status: 'failed', error: 'Method not suported yet.' });
});

app.listen(PORT, () => {
    console.log(`Server listening at port ${PORT}!`);
})
