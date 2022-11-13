import { Request, Response, NextFunction } from 'express';
import HttpException from '../exceptions/HttpException';

const errorHandler = (err: HttpException, req: Request, res: Response, next: NextFunction) => {
  console.error(err);
  res.status(err.statusCode ?? 500).json({ error: err.message });
  next();
};

export default errorHandler;
