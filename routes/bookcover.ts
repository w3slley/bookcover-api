import express, { Request, Response } from 'express';
import { getBookcoverFromISBN, getBookcoverUrl } from '../handlers/bookcover';
const router = express.Router();

router.get('/', (req: Request, res: Response) => {
    try {
        return getBookcoverUrl(req, res);
    }
    catch(error){
        res.status(500).json({status: 'failed', error: error.message});
    }
});

router.get('/:id', (req: Request, res: Response) => {
    try {
        return getBookcoverFromISBN(req, res);
    } catch (error) {
        res.status(500).json({status: 'failed', error: error.message});
    }
});

module.exports = router;