import dotenv from "dotenv";
import express from "express";
import mongoose from "mongoose";

import authRoutes from "./routes/auth.route";

dotenv.config();

const app = express();
const port = Number(process.env.PORT) || 3000;

mongoose
	.connect(process.env.MONGODB_URI!)
	.then(() => {
		console.log("MongoDB connected");
	})
	.catch(err => {
		if (err instanceof Error) {
			console.error(err.message);
		}
	});

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.use("/auth", authRoutes);

app.listen(port, () => {
	console.log(`Server started on port ${port}`);
});
