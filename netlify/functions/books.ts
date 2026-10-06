import type { Config, Context } from "@netlify/functions";
import { asc, eq, sql } from "drizzle-orm";
import { db } from "../../db/index.js";
import { books } from "../../db/schema.js";
import { authenticate } from "../lib/auth.js";
import { INVALID_BOOK_MESSAGE, parseBookInput, serializeBook, thicknessFromPages } from "../lib/books.js";
import {
  invalidID,
  isForeignKeyViolation,
  json,
  message,
  methodNotAllowed,
  notFound,
  parseID,
  readJSON,
  serverError,
} from "../lib/http.js";

const MISSING_CATEGORY = "category_id tidak tersedia";

export default async (req: Request, context: Context) => {
  try {
    const actor = await authenticate(req);
    if (actor instanceof Response) return actor;

    const rawID = context.params.id;
    if (rawID === undefined) {
      if (req.method === "GET") {
        const rows = await db.select().from(books).orderBy(asc(books.id));
        return json({ data: rows.map(serializeBook) });
      }
      if (req.method === "POST") {
        const input = parseBookInput(await readJSON(req));
        if (!input) return message(INVALID_BOOK_MESSAGE, 400);
        try {
          const [row] = await db
            .insert(books)
            .values({ ...input, thickness: thicknessFromPages(input.totalPage), createdBy: actor, modifiedBy: actor })
            .returning();
          return json({ message: "buku berhasil ditambahkan", data: serializeBook(row) }, 201);
        } catch (err) {
          if (isForeignKeyViolation(err)) return message(MISSING_CATEGORY, 400);
          throw err;
        }
      }
      return methodNotAllowed();
    }

    const id = parseID(rawID);
    if (id === null) return invalidID();

    if (req.method === "GET") {
      const [row] = await db.select().from(books).where(eq(books.id, id));
      return row ? json({ data: serializeBook(row) }) : notFound();
    }

    if (req.method === "PUT") {
      const input = parseBookInput(await readJSON(req));
      if (!input) return message(INVALID_BOOK_MESSAGE, 400);
      try {
        const [row] = await db
          .update(books)
          .set({ ...input, thickness: thicknessFromPages(input.totalPage), modifiedAt: sql`NOW()`, modifiedBy: actor })
          .where(eq(books.id, id))
          .returning();
        return row ? json({ message: "buku berhasil diperbarui", data: serializeBook(row) }) : notFound();
      } catch (err) {
        if (isForeignKeyViolation(err)) return message(MISSING_CATEGORY, 400);
        throw err;
      }
    }

    if (req.method === "DELETE") {
      const deleted = await db.delete(books).where(eq(books.id, id)).returning({ id: books.id });
      if (deleted.length === 0) return notFound("buku tidak ditemukan");
      return message("buku berhasil dihapus", 200);
    }

    return methodNotAllowed();
  } catch (err) {
    return serverError(err);
  }
};

export const config: Config = {
  path: ["/api/books", "/api/books/:id"],
  method: ["GET", "POST", "PUT", "DELETE"],
};
