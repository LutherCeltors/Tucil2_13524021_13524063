# Voxelizer — Tucil2 IF2211 Strategi Algoritma 2025/2026

> Konversi model 3D `.obj` menjadi voxel (kubus-kubus kecil seragam) menggunakan **Octree** dan algoritma **Divide and Conquer**, ditulis dalam bahasa **Go**.

---

## Daftar Isi

- [Deskripsi Program](#deskripsi-program)
- [Fitur](#fitur)
- [Requirement & Instalasi](#requirement--instalasi)
- [Kompilasi](#kompilasi)
- [Cara Menjalankan](#cara-menjalankan)
- [Format Input `.obj`](#format-input-obj)
- [Output Program](#output-program)
- [Rekomendasi `max_depth`](#rekomendasi-max_depth)
- [Struktur Repository](#struktur-repository)
- [Author](#author)

---

## Deskripsi Program

Program ini mengkonversi sebuah model 3D dalam format `.obj` menjadi representasi voxel — susunan kubus-kubus kecil seragam seperti di Minecraft. Proses konversi memanfaatkan struktur data **Octree** dengan algoritma **Divide and Conquer**:

1. **Parse** file `.obj` — baca semua verteks dan face, triangulasi polygon jika perlu.
2. **Hitung root bounding cube** — kotak kubik terkecil yang mencakup seluruh model.
3. **Divide & Conquer (Octree)** — bagi kotak secara rekursif menjadi 8 oktant. Setiap oktant yang tidak mengandung segitiga langsung dipangkas (*pruned*).
4. **Uji interseksi SAT** — gunakan Separating Axis Theorem untuk menentukan apakah segitiga menyentuh suatu kotak.
5. **Bangkitkan voxel** — setiap leaf node pada kedalaman maksimum yang berpotongan dengan permukaan model menjadi satu voxel (8 verteks + 6 quad face).
6. **Tulis output** — simpan hasil sebagai file `.obj` baru dan hasilkan viewer 3D interaktif.

---

## Fitur

| Fitur | Keterangan |
|-------|-----------|
| ✅ Voxelization | Konversi `.obj` → voxel seragam berbasis Octree D&C |
| ✅ Validasi input | Deteksi format `.obj` tidak valid, indeks out-of-range, dll |
| ✅ CLI Report | Statistik lengkap: jumlah voxel, node octree, pruning, waktu |
| ✅ Concurrency | Goroutine paralel pada depth ≤ 4, semaphore cap 512 goroutine |
| ✅ Cross-platform | Windows, Linux, macOS |

---

## Requirement & Instalasi

**Satu-satunya requirement: Go 1.21 atau lebih baru.**

Tidak ada library eksternal. Semua fitur (termasuk viewer) menggunakan standard library Go.

### Install Go

Download dari **https://go.dev/dl/** dan ikuti instruksi untuk OS Anda.

Verifikasi instalasi:
```bash
go version
# Output contoh: go version go1.21.0 linux/amd64
```

---

## Kompilasi

### Linux / macOS
```bash
cd src
go build -o ../bin/voxelizer .
```

### Windows
```cmd
cd src
go build -o ..\bin\voxelizer.exe .
```

Executable akan tersimpan di folder `bin/`.

---

## Cara Menjalankan

```
./bin/voxelizer <input.obj> <max_depth> [--view]
```

### Parameter

| Parameter | Wajib | Keterangan |
|-----------|-------|-----------|
| `input.obj` | ✅ | Path ke file `.obj` yang ingin dikonversi |
| `max_depth` | ✅ | Kedalaman maksimum octree (integer positif, misal: `5`) |
| `--view` | ❌ | Buka viewer 3D di browser secara otomatis setelah konversi |

### Contoh Penggunaan

```bash
# Konversi dasar
./bin/voxelizer test/cube.obj 4

# Konversi dan langsung buka viewer 3D
./bin/voxelizer test/cube.obj 4 --view

# Model kompleks dengan detail lebih tinggi
./bin/voxelizer models/pumpkin.obj 7 --view
```

### Contoh Output CLI

```
Loading model: test/cube.obj
Loaded 8 vertices and 12 faces.
Voxelizing with max depth 4...

=== Voxelization Report ===
Voxel count        : 56
Vertex count       : 448
Face count         : 336
Octree depth       : 4

Octree nodes per depth:
  1 : 8
  2 : 48
  3 : 192
  4 : 512

Pruned (skipped) nodes per depth:
  1 : 0
  2 : 16
  3 : 136
  4 : 400

Elapsed time       : 2.11ms
Output saved to    : test/cube-voxelized.obj
Viewer saved to    : test/cube-voxelized-viewer.html
Tip: run with --view to open the 3D viewer automatically.
```


## Format Input `.obj`

Program memproses dua jenis baris:

- `v x y z` — verteks dengan koordinat float (contoh: `v 1.0 -0.5 2.3`)
- `f i j k ...` — face dengan indeks verteks 1-based (contoh: `f 1 2 3`). Polygon dengan lebih dari 3 sudut otomatis ditriangulasi.

Semua baris lain (`vt`, `vn`, `s`, `g`, `mtllib`, `usemtl`, `#`, dll.) **diabaikan**.

### Contoh File `.obj` Valid

```
# Komentar diabaikan
v 1.0  0.0  0.0
v 0.0  1.0  0.0
v 0.0  0.0  1.0
v 0.0  0.0  0.0
f 1 2 3
f 1 2 4
f 1 3 4
f 2 3 4
```

### Validasi yang Dilakukan

Program menolak input dan menampilkan pesan error jika:
- File tidak ditemukan atau tidak bisa dibaca.
- Baris `v` tidak memiliki tepat 3 nilai float.
- Baris `f` memiliki indeks di luar rentang `[1, jumlah_verteks]`.
- File tidak memiliki verteks atau face sama sekali.
- `max_depth` bukan integer positif.

---

## Output Program

Untuk input `path/ke/model.obj`, program menghasilkan:

| File | Keterangan |
|------|-----------|
| `path/ke/model-voxelized.obj` | Model voxel dalam format `.obj` |

---

## Struktur Repository

```
Tucil2_NIM1_NIM2/
├── src/
│   ├── main.go       # Entry point, CLI, argument parsing
│   ├── obj.go        # Parser & writer file .obj
│   ├── geometry.go   # Vec3, AABB, SAT triangle-box intersection
│   ├── octree.go     # Octree, algoritma D&C, concurrency
│   └── go.mod        # Go module definition
├── bin/              # Executable hasil kompilasi
├── test/             # File .obj untuk pengujian & hasil konversi
├── doc/              # Laporan tugas dalam format PDF
└── README.md
```

---

## Author

| Nama | NIM |
|------|-----|
| Natanael I. Manurung | 13524021 | 
| Marcel L. Sitorus | 13524 | 
