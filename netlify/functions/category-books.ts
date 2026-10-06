import type { Config, Context } from "@netlify/functions";
import { asc, eq } from "drizzle-orm";
import { db } from "../../db/index.js";
import { books, categories } from "../../db/schema.js";
import { authenticate } from "../lib/auth.js";
import { serializeBook, serializeCategory } from "../lib/books.js";
import { invalidID, json, notFound, parseID, serverError } from "../lib/http.js";

export default async (req: Request, context: Context) => {
  try {
    const actor = await authenticate(req);
    if (actor instanceof Response) return actor;

    const id = parseID(context.params.id);
    if (id === null) return invalidID();

    const [category] = await db.select().from(categories).where(eq(categories.id, id));
    if (!category) return notFound();

    const rows = await db.select().from(books).where(eq(books.categoryId, id)).orderBy(asc(books.id));
    return json({ category: serializeCategory(category), data: rows.map(serializeBook) });
  } catch (err) {
    return serverError(err);
  }
};

export const config: Config = {
  path: "/api/categories/:id/books",
  method: "GET",
};
