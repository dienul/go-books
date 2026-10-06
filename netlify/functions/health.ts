import type { Config } from "@netlify/functions";
import { sql } from "drizzle-orm";
import { db } from "../../db/index.js";
import { json } from "../lib/http.js";

export default async () => {
  try {
    await db.execute(sql`SELECT 1`);
    return json({ status: "ok", message: "API dan database terhubung" });
  } catch (err) {
    console.error(err);
    return json({ status: "error", message: "database tidak tersedia" }, 503);
  }
};

export const config: Config = {
  path: "/health",
  method: "GET",
};
