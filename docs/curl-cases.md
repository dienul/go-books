# cURL Test Cases - Go Books API

Dokumen ini menggunakan server lokal `http://localhost:8080`. Jalankan aplikasi terlebih dahulu:

```bash
go run ./cmd/api
```

Siapkan variabel terminal berikut. Ganti nilai username dan password dengan akun admin lokal Anda.

```bash
export BASE_URL="http://localhost:8080"
export ADMIN_USERNAME="admin"
export ADMIN_PASSWORD="password-admin-anda"
```

> Jangan menyimpan password asli atau token JWT ke Git.

## 1. Health Check

### API dan database terhubung

Ekspektasi: HTTP `200`.

```bash
curl --include "$BASE_URL/health"
```

## 2. Authentication

### Login berhasil

Ekspektasi: HTTP `200` dan respons memiliki properti `token`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"$ADMIN_PASSWORD\"}" \
  "$BASE_URL/api/users/login"
```

Salin token dari respons login, lalu simpan sementara di terminal:

```bash
export TOKEN="tempel-token-jwt-di-sini"
```

### Login dengan password salah

Ekspektasi: HTTP `401`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"password-salah\"}" \
  "$BASE_URL/api/users/login"
```

### Login tanpa username

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Content-Type: application/json" \
  --data '{"password":"password"}' \
  "$BASE_URL/api/users/login"
```

### Mengakses API tanpa token

Ekspektasi: HTTP `401`.

```bash
curl --include "$BASE_URL/api/categories"
```

### Mengakses API dengan token tidak valid

Ekspektasi: HTTP `401`.

```bash
curl --include \
  --header "Authorization: Bearer token-tidak-valid" \
  "$BASE_URL/api/categories"
```

## 3. Categories

Semua endpoint kategori memerlukan header berikut:

```text
Authorization: Bearer $TOKEN
```

### Menambahkan kategori

Ekspektasi: HTTP `201`. Catat nilai `id` dari respons sebagai `CATEGORY_ID`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"Pemrograman"}' \
  "$BASE_URL/api/categories"
```

```bash
export CATEGORY_ID="1"
```

### Menambahkan kategori tanpa nama

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{}' \
  "$BASE_URL/api/categories"
```

### Menambahkan kategori dengan nama yang sama

Ekspektasi: HTTP `409`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"Pemrograman"}' \
  "$BASE_URL/api/categories"
```

### Menampilkan seluruh kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories"
```

### Menampilkan detail kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/$CATEGORY_ID"
```

### Detail kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/999999"
```

### ID kategori tidak valid

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/abc"
```

### Memperbarui kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"Teknologi"}' \
  "$BASE_URL/api/categories/$CATEGORY_ID"
```

### Memperbarui kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"Tidak Ada"}' \
  "$BASE_URL/api/categories/999999"
```

### Menampilkan buku berdasarkan kategori

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/$CATEGORY_ID/books"
```

## 4. Books

### Menambahkan buku tipis

`thickness` tidak dikirim melalui JSON. Karena `total_page` adalah `100`, server harus menyimpannya sebagai `tipis`.

Ekspektasi: HTTP `201`. Catat `id` dari respons sebagai `BOOK_ID`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Belajar Go\",\"description\":\"Dasar pemrograman Go\",\"image_url\":\"https://example.com/go.jpg\",\"release_year\":2024,\"price\":125000,\"total_page\":100,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books"
```

```bash
export BOOK_ID="1"
```

### Menambahkan buku tebal

Karena `total_page` lebih dari `100`, server harus menyimpan `thickness` sebagai `tebal`.

Ekspektasi: HTTP `201`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Go Lanjutan\",\"description\":\"Pemrograman Go lanjutan\",\"image_url\":\"https://example.com/go-advanced.jpg\",\"release_year\":2023,\"price\":175000,\"total_page\":101,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books"
```

### Tahun terbit di bawah 1980

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Buku Lama\",\"description\":\"Tahun tidak valid\",\"image_url\":\"https://example.com/old.jpg\",\"release_year\":1979,\"price\":50000,\"total_page\":80,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books"
```

### Tahun terbit di atas 2024

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Buku Masa Depan\",\"description\":\"Tahun tidak valid\",\"image_url\":\"https://example.com/future.jpg\",\"release_year\":2025,\"price\":50000,\"total_page\":80,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books"
```

### Menambahkan buku dengan kategori tidak tersedia

Ekspektasi: HTTP `400` atau `404`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"title":"Buku Tanpa Kategori","description":"Kategori tidak ada","image_url":"https://example.com/book.jpg","release_year":2024,"price":50000,"total_page":80,"category_id":999999}' \
  "$BASE_URL/api/books"
```

### Menambahkan buku dengan field wajib kosong

Ekspektasi: HTTP `400`.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{}' \
  "$BASE_URL/api/books"
```

### Menampilkan seluruh buku

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/books"
```

### Menampilkan detail buku

Ekspektasi: HTTP `200`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/books/$BOOK_ID"
```

### Detail buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/books/999999"
```

### Memperbarui buku

Server harus menghitung ulang `thickness` berdasarkan `total_page` yang baru.

Ekspektasi: HTTP `200` dan `thickness` bernilai `tebal`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Belajar Go Edisi Revisi\",\"description\":\"Dasar Go yang diperbarui\",\"image_url\":\"https://example.com/go-revised.jpg\",\"release_year\":2024,\"price\":150000,\"total_page\":150,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books/$BOOK_ID"
```

### Memperbarui buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request PUT \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data "{\"title\":\"Tidak Ada\",\"description\":\"Buku tidak ada\",\"image_url\":\"https://example.com/missing.jpg\",\"release_year\":2024,\"price\":10000,\"total_page\":50,\"category_id\":$CATEGORY_ID}" \
  "$BASE_URL/api/books/999999"
```

### Menghapus buku yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/books/999999"
```

### Menghapus buku

Ekspektasi: HTTP `200` atau `204`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/books/$BOOK_ID"
```

## 5. Category Delete Cases

Jalankan bagian ini setelah pengujian buku karena kategori masih diperlukan untuk membuat buku.

### Menghapus kategori yang masih memiliki buku

Ekspektasi: HTTP `409`. Gunakan kategori yang masih mempunyai buku.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/$CATEGORY_ID"
```

### Menghapus kategori yang tidak tersedia

Ekspektasi: HTTP `404`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/999999"
```

### Menghapus kategori kosong

Buat kategori baru yang belum memiliki buku, simpan ID-nya sebagai `EMPTY_CATEGORY_ID`, kemudian hapus.

```bash
curl --include \
  --request POST \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"Kategori Sementara"}' \
  "$BASE_URL/api/categories"
```

```bash
export EMPTY_CATEGORY_ID="2"
```

Ekspektasi: HTTP `200` atau `204`.

```bash
curl --include \
  --request DELETE \
  --header "Authorization: Bearer $TOKEN" \
  "$BASE_URL/api/categories/$EMPTY_CATEGORY_ID"
```
