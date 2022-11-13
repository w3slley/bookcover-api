import express, { Request, Response } from 'express';
import { getBookcoverFromISBN, getBookcoverUrl } from '../handlers/bookcover';
import inputValidation from '../middlewares/inputValidation';
const router = express.Router();

router.get('/', inputValidation, getBookcoverUrl);
router.get('/:id', getBookcoverFromISBN);

module.exports = router;