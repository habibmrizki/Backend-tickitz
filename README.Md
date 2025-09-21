# Backend Tickitz App With Golang 🎬

Tickitz App adalah aplikasi pemesanan tiket bioskop yang memungkinkan pengguna untuk:

- Menjelajahi film yang sedang tayang dan yang akan datang
- Memilih tempat duduk sesuai preferensi
- Melakukan pemesanan tiket secara online dengan mudah

Repositori ini berisi source code **Backend RESTful API** Tickitz App yang dibangun menggunakan bahasa **Go (Golang)** dengan pendekatan clean architecture.

## 🔧 Tech Stack

- **Go** (Golang)
- **Gin** – HTTP Web Framework
- **PostgreSQL** – Database relational
- **Swagger** – Dokumentasi API
- **JWT** – Autentikasi berbasis token
- **Go-Migrate** – Manajemen migrasi database
- **Clean Architecture** – Pemisahan tanggung jawab antara layer (handler, routes, repository)

## 📂 Struktur Direktori

```bash
.
├── cmd/                 # Entry point aplikasi
├── db/                  # Skrip migrasi database
├── docs/                # Dokumentasi Swagger
├── internal/            # Implementasi utama
│   ├── handler/         # Layer controller/handler
│   ├── repository/      # Layer repository untuk database
│   ├── middleware/      # Layer untuk memproses request
│   ├── models/          # Struct model
│   └── routes/          # Routing endpoint
├── .env                 # Konfigurasi environment
├── go.mod               # File module Go
└── README.md            # Dokumentasi proyek ini
```

## 🔐 Autentikasi

Gunakan **Bearer Token (JWT)** untuk endpoint yang membutuhkan autentikasi. Token didapat setelah login/register.

Contoh penggunaan di header:

```
Authorization: Bearer <your_token_here>
```

## 📬 API Endpoint (Contoh)

### Movies Routes

- `GET /movies` – Ambil semua daftar film yang ada
- `GET /movies/:id` – Ambil semua data detail film yang tertentu
- `GET /movies/moviespopular` – Ambil daftar film yang sedang popular (paling banyak peminat)
- `GET /movies/moviesupcoming` – Ambil daftar film yang akan tayang

### Admin Routes

- `POST /movies` – Menambahkan daftar movies sekaligus memberi jadwal
- `PUT /movies/:id` – Melakukan update data movies sekaligus memberi jadwal
- `DELETE /movies/:id` – Menghapus movies tertentu sekaligus menghapus jadwalnya
