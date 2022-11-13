import { Request, Response, NextFunction } from 'express';
const url = require('url');

const inputValidation = (req: Request, res: Response, next: NextFunction) => {
  const requiredParams = [
    'book_title',
    'author_name',
  ];
  const missingParams = [];
  const query = new URLSearchParams(url.parse(req.url).query);
  requiredParams.forEach(param => {
    if (!query.has(param)) {
      missingParams.push(param);
    }
  });
  if (missingParams.length) {
    const missingParamsStringified = missingParams.reduce((prev, curr) => (prev + ', ' + curr));
    return res.status(400).json({ status: 'failed', error: `Please insert the following required query parameters: ${missingParamsStringified}` });
  }

  next();
}

export default inputValidation;
