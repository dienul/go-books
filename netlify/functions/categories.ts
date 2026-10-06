import type { Config, Context } from "@netlify/functions";
import { asc, eq, sql } from "drizzle-orm";
import { db } from "../../db/index.js";
import { categories } from "../../db/schema.js";
import { authenticate } from "../lib/auth.js";
import { serializeCategory } from "../lib/books.js";
import {
  invalidID,
  isForeignKeyViolation,
  isUniqueViolation,
  json,
  message,
  methodNotAllowed,
  notFound,
  parseID,
  readJSON,
  serverError,
} from "../lib/http.js";

const INVALID_NAME = "nama kategori wajib diisi, maksimal 100 karakter";

function parseName(body: Record<string, unknown> | null): string | null {
  const name = body?.name;
  if (typeof name !== "string" || name.length > 100 || name.trim() === "") return null;
  return name.trim();
}

export default async (req: Request, context: Context) => {
  try {
    const actor = await authenticate(req);
    if (actor instanceof Response) return actor;

    const rawID = context.params.id;
    if (rawID === undefined) {
      if (req.method === "GET") {
        const rows = await db.select().from(categories).orderBy(asc(categories.id));
        return json({ data: rows.map(serializeCategory) });
      }
      if (req.method === "POST") {
        const name = parseName(await readJSON(req));
        if (!name) return message(INVALID_NAME, 400);
        try {
          const [row] = await db
            .insert(categories)
            .values({ name, createdBy: actor, modifiedBy: actor })
            .returning();
          return json({ message: "kategori berhasil ditambahkan", data: serializeCategory(row) }, 201);
        } catch (err) {
          if (isUniqueViolation(err)) return message("nama kategori sudah tersedia", 409);
          throw err;
        }
      }
      return methodNotAllowed();
    }

    const id = parseID(rawID);
    if (id === null) return invalidID();

    if (req.method === "GET") {
      const [row] = await db.select().from(categories).where(eq(categories.id, id));
      return row ? json({ data: serializeCategory(row) }) : notFound();
    }

    if (req.method === "PUT") {
      const name = parseName(await readJSON(req));
      if (!name) return message(INVALID_NAME, 400);
      try {
        const [row] = await db
          .update(categories)
          .set({ name, modifiedAt: sql`NOW()`, modifiedBy: actor })
          .where(eq(categories.id, id))
          .returning();
        return row ? json({ message: "kategori berhasil diperbarui", data: serializeCategory(row) }) : notFound();
      } catch (err) {
        if (isUniqueViolation(err)) return message("nama kategori sudah tersedia", 409);
        throw err;
      }
    }

    if (req.method === "DELETE") {
      try {
        const deleted = await db.delete(categories).where(eq(categories.id, id)).returning({ id: categories.id });
        if (deleted.length === 0) return notFound("kategori tidak ditemukan");
        return message("kategori berhasil dihapus", 200);
      } catch (err) {
        if (isForeignKeyViolation(err)) {
          return message("kategori masih memiliki buku dan tidak dapat dihapus", 409);
        }
        throw err;
      }
    }

    return methodNotAllowed();
  } catch (err) {
    return serverError(err);
  }
};

export const config: Config = {
  path: ["/api/categories", "/api/categories/:id"],
  method: ["GET", "POST", "PUT", "DELETE"],
};
