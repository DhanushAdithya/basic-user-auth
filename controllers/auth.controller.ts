import dotenv from "dotenv";
import type { Request, Response, NextFunction } from "express";
import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";
import User, { type UserType } from "../models/user.model";

dotenv.config();

const signJWT = (user: UserType) => {
	const expiresIn = Number(process.env.EXPIRES_IN) || 24 * 60 * 60;
	try {
		const token = jwt.sign({ email: user.email }, process.env.JWT_SECRET!, {
			expiresIn,
			issuer: process.env.JWT_ISSUER,
		});
		return { error: null, token };
	} catch (err) {
		if (err instanceof Error) {
			return { error: err.message, token: null };
		} else {
			return { error: "Unknown Error", token: null };
		}
	}
};

const register = async (req: Request, res: Response, _next: NextFunction) => {
	const { email, password } = req.body;

	try {
		const hash = await bcrypt.hash(password, 10);
		const user = await User.create({ email, password: hash });
		return res.status(201).json({ user });
	} catch (err) {
		if (err instanceof Error) {
			return res.status(500).json({
				message: err.message,
				err,
			});
		}
	}
};

const login = async (req: Request, res: Response, _next: NextFunction) => {
	const { email, password } = req.body;

	try {
		const user = await User.findOne({ email });
		if (!user) {
			return res.status(401).json({
				message: "Unauthorized",
			});
		}
		const passMatch = await bcrypt.compare(password, user.password);
		if (!passMatch) {
			return res.status(401).json({
				message: "Wrong Password",
			});
		}
		const { token, error } = signJWT(user);
		if (error) {
			return res.status(500).json({
				message: error,
			});
		}
		return res.status(200).json({ token, user });
	} catch (err) {
		if (err instanceof Error) {
			return res.status(500).json({
				message: err.message,
				err,
			});
		}
	}
};

const confidential = (_req: Request, res: Response, _next: NextFunction) => {
	return res.status(200).json({
		message: "boom",
	});
};

export { register, login, confidential };
