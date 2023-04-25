const PORT = process.env.PORT || 2000;
import express, { Request, Response, NextFunction } from 'express';
import HttpException from './exceptions/HttpException';
import { METHOD_NOT_SUPPORTED } from './helpers/messages';
import errorHandler from './middlewares/errorHandler';
import jsonHeader from './middlewares/headers';
import allowOrigin from './middlewares/cors';
const app = express();
require('dotenv').config();

app.use(allowOrigin);
app.use(jsonHeader);
app.use('/bookcover', require('./routes/bookcover'));
app.get('/*', (req: Request, res: Response, next: NextFunction) => {
  next(new HttpException(400, METHOD_NOT_SUPPORTED));
});
app.use(errorHandler);

app.listen(PORT, () => {
  console.log(`Server listening at port ${PORT}!`);
})
