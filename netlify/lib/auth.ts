import { SignJWT, jwtVerify } from "jose";
import { message } from "./http.js";

const TOKEN_TTL_SECONDS = 24 * 60 * 60;

function secretKey(): Uint8Array {
  const secret = Netlify.env.get("JWT_SECRET") ?? "";
  if (secret.length < 32) {
    throw new Error("JWT_SECRET wajib diatur dan minimal 32 karakter");
  }
  return new TextEncoder().encode(secret);
}

export async function signToken(username: string) {
  const token = await new SignJWT({ username })
    .setProtectedHeader({ alg: "HS256" })
    .setSubject(username)
    .setIssuedAt()
    .setExpirationTime(`${TOKEN_TTL_SECONDS}s`)
    .sign(secretKey());
  return { token, expiresIn: TOKEN_TTL_SECONDS };
}

// authenticate mengembalikan username dari token Bearer, atau Response 401.
export async function authenticate(req: Request): Promise<string | Response> {
  const parts = (req.headers.get("authorization") ?? "").trim().split(/\s+/);
  if (parts.length !== 2 || parts[0].toLowerCase() !== "bearer") {
    return message("token autentikasi diperlukan", 401);
  }
  const key = secretKey();
  try {
    const { payload } = await jwtVerify(parts[1], key, { algorithms: ["HS256"] });
    if (typeof payload.username !== "string" || payload.username === "") throw new Error("username kosong");
    return payload.username;
  } catch {
    return message("token tidak valid atau kedaluwarsa", 401);
  }
}
