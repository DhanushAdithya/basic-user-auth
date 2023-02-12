import { Router } from "express";
import { confidential, login, register } from "../controllers/auth.controller";
import { verifyJWT } from "../middlewares/auth.middleware";

const router = Router();

router.get("/confidential", verifyJWT, confidential);
router.post("/register", register);
router.post("/login", login);

export default router;
