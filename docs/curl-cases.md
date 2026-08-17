# cURL Test Cases - Go Books API

Dokumen ini menguji deployment production berikut secara langsung:

```text
https://go-books-production-a329.up.railway.app
```

Semua contoh sengaja tidak menggunakan `export`. Ganti placeholder berikut secara manual:

- `PASSWORD_ADMIN_ASLI`: nilai `ADMIN_PASSWORD` pada service aplikasi Railway.
- `TOKEN_ANDA`: token dari respons login.
- ID `1` atau `2`: ID sebenarnya dari respons create.

> Jangan menyimpan password asli atau token JWT ke Git.

## 0. Persiapan Database dan Akun Admin

Pada service aplikasi Railway, pastikan variables berikut tersedia:

```env
DATABASE_URL=${{Postgres.DATABASE_URL}}
JWT_SECRET=string-acak-minimal-32-karakter
ADMIN_USERNAME=admin
ADMIN_PASSWORD=password-kuat-minimal-8-karakter
GIN_MODE=release
```

`Postgres` harus sesuai dengan nama service database. Klik **Deploy Changes** setelah mengubah variables. Aplikasi otomatis menjalankan migration dan membuat akun admin ketika tabel `users` masih kosong.

Periksa akun melalui Railway Database:

```sql
SELECT id, username, created_at FROM users;
```

Jika akun `admin` lama memakai password berbeda dan belum ada data user penting, hapus lalu redeploy:

```sql
DELETE FROM users WHERE username = 'admin';
```

## 1. Health Check

### API dan database terhubung

Ekspektasi: HTTP `200`.

```bash
curl --include "https://go-books-production-a329.up.railway.app/health"
```

## 2. Authentication

### Login berhasil

Ekspektasi: HTTP `200` dan respons memiliki properti `token`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"username\":\"admin\",\"password\":\"PASSWORD_ADMIN_ASLI\"}" \
  "https://go-books-production-a329.up.railway.app/api/users/login"
```

Salin nilai `token` dari respons login dan ganti setiap `TOKEN_ANDA` pada case berikutnya.

### Login dengan password salah

Ekspektasi: HTTP `401`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"username\":\"admin\",\"password\":\"password-salah\"}" \
  "https://go-books-production-a329.up.railway.app/api/users/login"
```

### Login tanpa username

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data '{"password":"password"}' \
  "https://go-books-production-a329.up.railway.app/api/users/login"
```

### Mengakses API tanpa token

Ekspektasi: HTTP `401`.

```bash
curl --include "https://go-books-production-a329.up.railway.app/api/categories"
```

### Mengakses API dengan token tidak valid

Ekspektasi: HTTP `401`.

```bash
curl --include \
  --header "Authorization: Bearer token-tidak-valid" \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

## 3. Categories

Semua endpoint kategori memerlukan header berikut:

```text
Authorization: Bearer TOKEN_ANDA
```

### Menambahkan kategori

Ekspektasi: HTTP `201`. Catat nilai `id` dari respons sebagai `CATEGORY_ID`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"name":"Pemrograman"}' \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

Pada case berikutnya, ganti ID `1` dengan `data.id` dari respons ini.

### Menambahkan kategori tanpa nama

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{}' \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

### Menambahkan kategori dengan nama yang sama

Ekspektasi: HTTP `409`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"name":"Pemrograman"}' \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

### Menampilkan seluruh kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

### Menampilkan detail kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/1"
```

### Detail kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/999999"
```

### ID kategori tidak valid

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/abc"
```

### Memperbarui kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"name":"Teknologi"}' \
  "https://go-books-production-a329.up.railway.app/api/categories/1"
```

### Memperbarui kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"name":"Tidak Ada"}' \
  "https://go-books-production-a329.up.railway.app/api/categories/999999"
```

### Menampilkan buku berdasarkan kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/1/books"
```

## 4. Books

### Menambahkan buku tipis

`thickness` tidak dikirim melalui JSON. Karena `total_page` adalah `100`, server harus menyimpannya sebagai `tipis`.

Ekspektasi: HTTP `201`. Catat `id` dari respons sebagai `BOOK_ID`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Belajar Go\",\"description\":\"Dasar pemrograman Go\",\"image_url\":\"https://example.com/go.jpg\",\"release_year\":2024,\"price\":125000,\"total_page\":100,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books"
```

Pada case berikutnya, ganti ID buku `1` dengan `data.id` dari respons ini.

### Menambahkan buku tebal

Karena `total_page` lebih dari `100`, server harus menyimpan `thickness` sebagai `tebal`.

Ekspektasi: HTTP `201`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Go Lanjutan\",\"description\":\"Pemrograman Go lanjutan\",\"image_url\":\"https://example.com/go-advanced.jpg\",\"release_year\":2023,\"price\":175000,\"total_page\":101,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Tahun terbit di bawah 1980

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Buku Lama\",\"description\":\"Tahun tidak valid\",\"image_url\":\"https://example.com/old.jpg\",\"release_year\":1979,\"price\":50000,\"total_page\":80,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Tahun terbit di atas 2024

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Buku Masa Depan\",\"description\":\"Tahun tidak valid\",\"image_url\":\"https://example.com/future.jpg\",\"release_year\":2025,\"price\":50000,\"total_page\":80,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Menambahkan buku dengan kategori tidak tersedia

Ekspektasi: HTTP `400` atau `404`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"title":"Buku Tanpa Kategori","description":"Kategori tidak ada","image_url":"https://example.com/book.jpg","release_year":2024,"price":50000,"total_page":80,"category_id":999999}' \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Menambahkan buku dengan field wajib kosong

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{}' \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Menampilkan seluruh buku

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/books"
```

### Menampilkan detail buku

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/books/1"
```

### Detail buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/books/999999"
```

### Memperbarui buku

Server harus menghitung ulang `thickness` berdasarkan `total_page` yang baru.

Ekspektasi: HTTP `200` dan `thickness` bernilai `tebal`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Belajar Go Edisi Revisi\",\"description\":\"Dasar Go yang diperbarui\",\"image_url\":\"https://example.com/go-revised.jpg\",\"release_year\":2024,\"price\":150000,\"total_page\":150,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books/1"
```

### Memperbarui buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Tidak Ada\",\"description\":\"Buku tidak ada\",\"image_url\":\"https://example.com/missing.jpg\",\"release_year\":2024,\"price\":10000,\"total_page\":50,\"category_id\":1}" \
  "https://go-books-production-a329.up.railway.app/api/books/999999"
```

### Menghapus buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/books/999999"
```

### Menghapus buku

Ekspektasi: HTTP `200` atau `204`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/books/1"
```

## 5. Category Delete Cases

Jalankan bagian ini setelah pengujian buku karena kategori masih diperlukan untuk membuat buku.

### Menghapus kategori yang masih memiliki buku

Ekspektasi: HTTP `409`. Gunakan kategori yang masih mempunyai buku.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/1"
```

### Menghapus kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/999999"
```

### Menghapus kategori kosong

Buat kategori baru yang belum memiliki buku, simpan ID-nya sebagai `EMPTY_CATEGORY_ID`, kemudian hapus.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer TOKEN_ANDA" \
  --header "Content-Type: application/json" \
  --data '{"name":"Kategori Sementara"}' \
  "https://go-books-production-a329.up.railway.app/api/categories"
```

Ganti ID `2` dengan `data.id` kategori sementara dari respons di atas.

Ekspektasi: HTTP `200` atau `204`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer TOKEN_ANDA" \
  "https://go-books-production-a329.up.railway.app/api/categories/2"
```
