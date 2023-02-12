import dotenv from "dotenv";
import jwt from "jsonwebtoken";
import type { Request, Response, NextFunction } from "express";

dotenv.config();

const verifyJWT = (req: Request, res: Response, next: NextFunction) => {
	const token = req.headers.authorization?.split(" ")[1];

	if (token) {
		try {
			const decodedToken = jwt.verify(token, process.env.JWT_SECRET!);
			if (decodedToken) {
				next();
			} else {
				return res.status(401).json({
					message: "Unauthorized",
				});
			}
		} catch (err) {
			if (err instanceof Error) {
				return res.status(500).json({
					message: err.message,
					err,
				});
			}
		}
	} else {
		return res.status(401).json({
			message: "Unauthorized",
		});
	}
};

export { verifyJWT };
