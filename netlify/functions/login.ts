import type { Config } from "@netlify/functions";
import bcrypt from "bcryptjs";
import { eq } from "drizzle-orm";
import { db } from "../../db/index.js";
import { users } from "../../db/schema.js";
import { signToken } from "../lib/auth.js";
import { json, message, readJSON, serverError } from "../lib/http.js";

// seedAdmin membuat akun dari ADMIN_USERNAME dan ADMIN_PASSWORD jika belum tersedia.
async function seedAdmin() {
  const username = Netlify.env.get("ADMIN_USERNAME");
  const password = Netlify.env.get("ADMIN_PASSWORD");
  if (!username || !password || password.length < 8) return;

  const [existing] = await db.select({ id: users.id }).from(users).where(eq(users.username, username)).limit(1);
  if (existing) return;

  const hash = await bcrypt.hash(password, 10);
  await db
    .insert(users)
    .values({ username, password: hash, createdBy: username, modifiedBy: username })
    .onConflictDoNothing({ target: users.username });
}

export default async (req: Request) => {
  const body = await readJSON(req);
  const username = body?.username;
  const password = body?.password;
  if (typeof username !== "string" || username === "" || username.length > 100 || typeof password !== "string" || password === "") {
    return message("username dan password wajib diisi", 400);
  }

  try {
    await seedAdmin();

    const [user] = await db.select().from(users).where(eq(users.username, username)).limit(1);
    if (!user || !(await bcrypt.compare(password, user.password))) {
      return message("username atau password salah", 401);
    }

    const { token, expiresIn } = await signToken(user.username);
    return json({ message: "login berhasil", token, expires_in: expiresIn });
  } catch (err) {
    return serverError(err);
  }
};

export const config: Config = {
  path: "/api/users/login",
  method: "POST",
};
