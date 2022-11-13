import { Request, Response, NextFunction } from 'express';

const jsonHeader = (req: Request, res: Response, next: NextFunction) => {
  res.setHeader('Content-Type', 'application/json');
  next();
}
export default jsonHeader;