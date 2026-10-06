import { sql } from "drizzle-orm";
import { bigserial, bigint, check, index, integer, pgTable, timestamp, varchar } from "drizzle-orm/pg-core";

const audit = {
  createdAt: timestamp("created_at", { withTimezone: true }).notNull().defaultNow(),
  createdBy: varchar("created_by", { length: 100 }).notNull(),
  modifiedAt: timestamp("modified_at", { withTimezone: true }).notNull().defaultNow(),
  modifiedBy: varchar("modified_by", { length: 100 }).notNull(),
};

export const users = pgTable("users", {
  id: bigserial({ mode: "number" }).primaryKey(),
  username: varchar({ length: 100 }).notNull().unique(),
  password: varchar({ length: 255 }).notNull(),
  ...audit,
});

export const categories = pgTable("categories", {
  id: bigserial({ mode: "number" }).primaryKey(),
  name: varchar({ length: 100 }).notNull().unique(),
  ...audit,
});

export const books = pgTable(
  "books",
  {
    id: bigserial({ mode: "number" }).primaryKey(),
    title: varchar({ length: 255 }).notNull(),
    description: varchar({ length: 1000 }).notNull(),
    imageUrl: varchar("image_url", { length: 1000 }).notNull(),
    releaseYear: integer("release_year").notNull(),
    price: integer().notNull(),
    totalPage: integer("total_page").notNull(),
    thickness: varchar({ length: 10 }).notNull(),
    categoryId: bigint("category_id", { mode: "number" })
      .notNull()
      .references(() => categories.id, { onDelete: "restrict" }),
    ...audit,
  },
  (t) => [
    index("idx_books_category_id").on(t.categoryId),
    check("books_release_year_check", sql`${t.releaseYear} BETWEEN 1980 AND 2024`),
    check("books_price_check", sql`${t.price} >= 0`),
    check("books_total_page_check", sql`${t.totalPage} > 0`),
    check("books_thickness_check", sql`${t.thickness} IN ('tipis', 'tebal')`),
  ],
);
