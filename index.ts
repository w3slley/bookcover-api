import express, { Request, Response, NextFunction } from "express";
import HttpException from "./exceptions/HttpException";
import { METHOD_NOT_SUPPORTED } from "./helpers/messages";
import bookcoverRoute from './routes/bookcover';
import errorHandler from "./middlewares/errorHandler";
import jsonHeader from "./middlewares/headers";
import dotenv from "dotenv";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 2000;

app.use(jsonHeader);
app.use("/bookcover", bookcoverRoute);

app.get("/*", (req: Request, res: Response, next: NextFunction) => {
  next(new HttpException(400, METHOD_NOT_SUPPORTED));
});

app.use(errorHandler);

app.listen(PORT, () => {
  console.log(`Server listening at port ${PORT}!`);
})
