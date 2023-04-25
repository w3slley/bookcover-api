import { Request, Response, NextFunction } from 'express';

const allowOrigin = (req: Request, res: Response, next: NextFunction) => {
  res.header('Access-Control-Allow-Origin', '*');
  next();
};

export default allowOrigin;
