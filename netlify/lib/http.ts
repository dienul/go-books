export const json = (body: unknown, status = 200) => Response.json(body, { status });

export const message = (text: string, status: number) => json({ message: text }, status);

export const serverError = (err: unknown) => {
  console.error(err);
  return message("terjadi kesalahan pada server", 500);
};

export const notFound = (text = "data tidak ditemukan") => message(text, 404);

export const methodNotAllowed = () => message("method tidak didukung", 405);

// parseID mengembalikan null jika id bukan bilangan bulat positif.
export function parseID(raw: string | undefined): number | null {
  if (!raw || !/^\d+$/.test(raw)) return null;
  const id = Number(raw);
  return Number.isSafeInteger(id) && id > 0 ? id : null;
}

export const invalidID = () => message("id harus berupa bilangan bulat positif", 400);

export async function readJSON(req: Request): Promise<Record<string, unknown> | null> {
  try {
    const body = await req.json();
    return body && typeof body === "object" && !Array.isArray(body) ? body : null;
  } catch {
    return null;
  }
}

function pgCode(err: unknown): string | undefined {
  let current: any = err;
  while (current) {
    if (typeof current.code === "string") return current.code;
    current = current.cause;
  }
  return undefined;
}

export const isUniqueViolation = (err: unknown) => pgCode(err) === "23505";
export const isForeignKeyViolation = (err: unknown) => pgCode(err) === "23503";
