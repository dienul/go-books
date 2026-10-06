import { books, categories } from "../../db/schema.js";

type BookRow = typeof books.$inferSelect;
type CategoryRow = typeof categories.$inferSelect;

export type BookInput = {
  title: string;
  description: string;
  imageUrl: string;
  releaseYear: number;
  price: number;
  totalPage: number;
  categoryId: number;
};

export const INVALID_BOOK_MESSAGE =
  "data buku tidak valid; release_year harus 1980-2024 dan total_page harus lebih dari 0";

export const thicknessFromPages = (totalPage: number) => (totalPage > 100 ? "tebal" : "tipis");

const isNonEmptyString = (value: unknown, max: number): value is string =>
  typeof value === "string" && value !== "" && value.length <= max;

const isInt = (value: unknown): value is number => typeof value === "number" && Number.isSafeInteger(value);

function isURL(value: string) {
  try {
    const url = new URL(value);
    return url.protocol !== "" && url.host !== "";
  } catch {
    return false;
  }
}

// parseBookInput menerapkan aturan validasi yang sama dengan versi Go.
export function parseBookInput(body: Record<string, unknown> | null): BookInput | null {
  if (!body) return null;
  const { title, description, image_url, release_year, price, total_page, category_id } = body;
  if (!isNonEmptyString(title, 255)) return null;
  if (!isNonEmptyString(description, 1000)) return null;
  if (!isNonEmptyString(image_url, 1000) || !isURL(image_url)) return null;
  if (!isInt(release_year) || release_year < 1980 || release_year > 2024) return null;
  const priceValue = price ?? 0;
  if (!isInt(priceValue) || priceValue < 0 || priceValue > 2147483647) return null;
  if (!isInt(total_page) || total_page <= 0 || total_page > 2147483647) return null;
  if (!isInt(category_id) || category_id <= 0) return null;
  return {
    title,
    description,
    imageUrl: image_url,
    releaseYear: release_year,
    price: priceValue,
    totalPage: total_page,
    categoryId: category_id,
  };
}

export const serializeBook = (book: BookRow) => ({
  id: book.id,
  title: book.title,
  description: book.description,
  image_url: book.imageUrl,
  release_year: book.releaseYear,
  price: book.price,
  total_page: book.totalPage,
  thickness: book.thickness,
  category_id: book.categoryId,
  created_at: book.createdAt,
  created_by: book.createdBy,
  modified_at: book.modifiedAt,
  modified_by: book.modifiedBy,
});

export const serializeCategory = (category: CategoryRow) => ({
  id: category.id,
  name: category.name,
  created_at: category.createdAt,
  created_by: category.createdBy,
  modified_at: category.modifiedAt,
  modified_by: category.modifiedBy,
});
