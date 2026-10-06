# Go Books API

REST API untuk mengelola kategori dan data buku. Project menggunakan Go, Gin, PostgreSQL, JWT, bcrypt, dan `sql-migrate`.

## Fitur

- Login menggunakan JSON Web Token (JWT).
- Password disimpan sebagai bcrypt hash.
- CRUD kategori dan buku.
- Daftar buku berdasarkan kategori.
- Validasi `release_year` dari 1980 sampai 2024.
- Konversi otomatis `thickness`: maksimal 100 halaman adalah `tipis`, lebih dari 100 adalah `tebal`.
- Migration PostgreSQL otomatis saat aplikasi dimulai.
- Audit `created_at`, `created_by`, `modified_at`, dan `modified_by`.
- Error `404` untuk data yang tidak tersedia dan `409` jika kategori masih memiliki buku.

## Menjalankan secara lokal

Persyaratan:

- Go 1.26 atau lebih baru.
- PostgreSQL lokal atau PostgreSQL Railway yang dapat diakses melalui public URL.

Salin contoh environment dan isi nilainya:

```bash
cp .env.example .env
```

```env
DATABASE_URL=postgresql://username:password@host:port/database
JWT_SECRET=string-acak-minimal-32-karakter
ADMIN_USERNAME=admin
ADMIN_PASSWORD=password-admin-minimal-8-karakter
PORT=8080
```

Jangan commit `.env`. Jalankan aplikasi:

```bash
go mod tidy
go run ./cmd/api
```

Periksa koneksi API dan database:

```bash
curl http://localhost:8080/health
```

## Akun admin awal

Saat startup, aplikasi membuat akun dari `ADMIN_USERNAME` dan `ADMIN_PASSWORD` jika username tersebut belum tersedia. Password disimpan sebagai bcrypt hash. Mengubah `ADMIN_PASSWORD` setelah akun tercipta tidak otomatis mengubah password yang sudah tersimpan.

## Authentication

Login:

```bash
curl --request POST \
  --header "Content-Type: application/json" \
  --data '{"username":"admin","password":"PASSWORD_ADMIN"}' \
  http://localhost:8080/api/users/login
```

Gunakan token dari respons untuk endpoint kategori dan buku:

```text
Authorization: Bearer TOKEN_JWT
```

Token berlaku selama 24 jam.

## Endpoint

| Method | Path | Keterangan | Auth |
|---|---|---|---|
| GET | `/health` | Status API dan database | Tidak |
| POST | `/api/users/login` | Login dan memperoleh JWT | Tidak |
| GET | `/api/categories` | Semua kategori | JWT |
| POST | `/api/categories` | Tambah kategori | JWT |
| GET | `/api/categories/:id` | Detail kategori | JWT |
| PUT | `/api/categories/:id` | Update kategori | JWT |
| DELETE | `/api/categories/:id` | Hapus kategori | JWT |
| GET | `/api/categories/:id/books` | Buku berdasarkan kategori | JWT |
| GET | `/api/books` | Semua buku | JWT |
| POST | `/api/books` | Tambah buku | JWT |
| GET | `/api/books/:id` | Detail buku | JWT |
| PUT | `/api/books/:id` | Update buku | JWT |
| DELETE | `/api/books/:id` | Hapus buku | JWT |

`PUT` disediakan karena brief meminta validasi ketika meng-update data yang tidak tersedia.

## Format input

Kategori:

```json
{
  "name": "Pemrograman"
}
```

Buku:

```json
{
  "title": "Belajar Go",
  "description": "Dasar pemrograman Go",
  "image_url": "https://example.com/go.jpg",
  "release_year": 2024,
  "price": 125000,
  "total_page": 100,
  "category_id": 1
}
```

Jangan mengirim `thickness`; server menghitung nilainya dari `total_page`.

Contoh cURL untuk seluruh success dan error case tersedia di [docs/curl-cases.md](docs/curl-cases.md).

## Migration

Migration berada di folder `migrations/` dan diterapkan otomatis ketika aplikasi dimulai. Tabel `gorp_migrations` menyimpan migration yang sudah dijalankan.

Untuk perubahan schema berikutnya, buat file berurutan seperti:

```text
migrations/002_nama_perubahan.sql
```

Jangan mengubah migration lama yang sudah diterapkan.

## Pengujian

```bash
go test ./...
go vet ./...
go build ./...
```

Smoke test seluruh endpoint terhadap server yang sedang berjalan:

```bash
./scripts/smoke_test.sh http://localhost:8080
```

Smoke test membuat data sementara dan menghapusnya kembali setelah seluruh case berhasil.

## Deployment Railway

1. Hubungkan repository GitHub ke service Railway.
2. Tambahkan service PostgreSQL pada project yang sama.
3. Pada service aplikasi, tambahkan variables:

```env
DATABASE_URL=${{Postgres.DATABASE_URL}}
JWT_SECRET=string-acak-minimal-32-karakter
ADMIN_USERNAME=admin
ADMIN_PASSWORD=password-admin-minimal-8-karakter
GIN_MODE=release
```

4. Railway menyediakan `PORT` secara otomatis; aplikasi tetap memakai `8080` sebagai fallback lokal.
5. Deploy ulang dan periksa `/health`.

Gunakan private `DATABASE_URL` untuk aplikasi yang berjalan di Railway. `DATABASE_PUBLIC_URL` hanya diperlukan untuk aplikasi yang berjalan dari komputer lokal.

## Deployment Netlify

Project ini juga dapat berjalan sepenuhnya di Netlify. Endpoint yang sama (path, format request, respons, dan pesan error) disediakan oleh Netlify Functions di `netlify/functions/`, dengan data disimpan di Netlify Database (PostgreSQL terkelola).

- Schema database didefinisikan di `db/schema.ts` (Drizzle ORM) dan mengikuti `migrations/001_create_tables.sql`.
- Migration berada di `netlify/database/migrations/` dan diterapkan otomatis oleh Netlify setiap deploy. Untuk perubahan schema berikutnya, ubah `db/schema.ts` lalu jalankan `npx drizzle-kit generate --name nama_perubahan`.
- Koneksi database dikonfigurasi otomatis oleh Netlify; `DATABASE_URL` tidak diperlukan.

Atur environment variables berikut di **Project configuration → Environment variables**:

```env
JWT_SECRET=string-acak-minimal-32-karakter
ADMIN_USERNAME=admin
ADMIN_PASSWORD=password-admin-minimal-8-karakter
```

Akun admin dibuat otomatis saat login pertama jika username tersebut belum tersedia.

Menjalankan secara lokal:

```bash
npm install
netlify dev
```
