# Struktur Project - Project Management API

Dokumentasi ini menjelaskan fungsi dan tanggung jawab masing-masing folder dalam struktur project ini.

## 📁 Struktur Folder

```
project-management/
├── controllers/     # Menangani HTTP Request & Response
├── models/          # Definisi struktur data dan database
├── repositories/    # Operasi database langsung
├── services/        # Business logic aplikasi
└── utils/           # Fungsi helper yang reusable
```

---

## 📂 Controllers

**Fungsi**: Menangani **HTTP Request & Response** - ini adalah gerbang utama aplikasi Anda.

### Tanggung Jawab:
- Menerima request dari client (API endpoints)
- Validasi input dan parsing data dari request body/params/query
- Mengambil data user dari JWT token
- Memanggil service untuk pemrosesan data
- Mengembalikan response dalam format JSON yang konsisten

### Contoh:
```go
// controllers/board_controller.go
func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
    // 1. Parsing request body
    board := new(models.Board)
    ctx.BodyParser(board)

    // 2. Mengambil user dari JWT token
    user := ctx.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)

    // 3. Memanggil service
    c.service.Create(board)

    // 4. Mengembalikan response
    return utils.Success(ctx, "Board berhasil dibuat", board)
}
```

---

## 📦 Models

**Fungsi**: Mendefinisikan **struktur data** dan representasi database.

### Tanggung Jawab:
- Struct yang merepresentasikan tabel database
- Field properties dengan tag untuk JSON dan database
- Relasi antar entitas (Board, User, List, Card, dll)
- Tipe data kustom

### Contoh:
```go
// models/board.go
type Board struct {
    InternalID       int64      `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
    PublicID         uuid.UUID  `json:"public_id" db:"public_id"`
    Title            string     `json:"title" db:"title"`
    Description      string     `json:"description" db:"description"`
    OwnerID          int64      `json:"owner_internal_id" db:"owner_internal_id"`
    OwnerPublicID    uuid.UUID  `json:"owner_public_id" db:"owner_public_id"`
    CreatedAt        time.Time  `json:"created_at" db:"created_at"`
    DueDate          *time.Time `json:"due_date,omitempty" db:"due_date"`
}
```

### List Models yang Tersedia:
- `user.go` - Data pengguna
- `board.go` - Data board/proyek
- `list.go` - Data list dalam board
- `card.go` - Data card dalam list
- `board_member.go` - Relasi user dan board
- `list_position.go` - Posisi list dalam board
- `card_position.go` - Posisi card dalam list
- `comment.go`, `label.go`, `card_label.go`, `card_attachment.go`, `card_assignee.go`

---

## 💾 Repositories

**Fungsi**: Menangani **operasi database** langsung.

### Tanggung Jawab:
- CRUD operations (Create, Read, Update, Delete)
- Query database menggunakan GORM
- Implementasi akses data yang spesifik
- Transaksi database jika diperlukan

### Contoh:
```go
// repositories/board_repository.go
func (r *boardRepository) Create(board *models.Board) error {
    return config.DB.Create(board).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
    var board models.Board
    err := config.DB.Where("public_id = ?", publicID).First(&board).Error
    return &board, err
}
```

### List Repositories yang Tersedia:
- `user_repository.go` - Operasi database user
- `board_repository.go` - Operasi database board
- `board_member_repository.go` - Operasi database member board
- `list_repository.go` - Operasi database list
- `list_position_repository.go` - Operasi database posisi list

---

## 🔧 Services

**Fungsi**: Menangani **business logic** dan logika aplikasi.

### Tanggung Jawab:
- Koordinasi antar repositories
- Validasi bisnis yang kompleks
- Transformasi data
- Implementasi rules bisnis

### Contoh:
```go
// services/board_service.go
func (s *boardService) Create(board *models.Board) error {
    // 1. Validasi: Cek apakah owner ada
    user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
    if err != nil {
        return errors.New("owner not found")
    }

    // 2. Generate UUID baru
    board.PublicID = uuid.New()
    board.OwnerID = user.InternalID

    // 3. Simpan ke database
    return s.boardRepo.Create(board)
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDS []string) error {
    // 1. Cari board
    board, err := s.boardRepo.FindByPublicID(boardPublicID)
    if err != nil {
        return errors.New("board not found")
    }

    // 2. Validasi dan convert user public ID ke internal ID
    var userInternalIDs []uint
    for _, userPublicID := range userPublicIDS {
        user, err := s.userRepo.FindByPublicID(userPublicID)
        if err != nil {
            return errors.New("User not found: " + userPublicID)
        }
        userInternalIDs = append(userInternalIDs, uint(user.InternalID))
    }

    // 3. Cek member yang sudah ada (duplicate check)
    existingMembers, _ := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
    memberMap := make(map[uint]bool)
    for _, member := range existingMembers {
        memberMap[uint(member.InternalID)] = true
    }

    // 4. Filter hanya member baru
    var newMembersIDs []uint
    for _, userID := range userInternalIDs {
        if !memberMap[userID] {
            newMembersIDs = append(newMembersIDs, userID)
        }
    }

    // 5. Tambahkan member baru
    return s.boardRepo.AddMember(uint(board.InternalID), newMembersIDs)
}
```

### List Services yang Tersedia:
- `user_service.go` - Business logic user
- `board_service.go` - Business logic board
- `list_service.go` - Business logic list

---

## 🛠️ Utils

**Fungsi**: Menyediakan **fungsi helper** yang dapat digunakan ulang.

### Tanggung Jawab:
- Fungsi umum untuk berbagai keperluan
- Tidak terikat pada entity tertentu
- Memudahkan formatting dan operasi umum

### Contoh:
```go
// utils/response.go
func Success(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Status:        "Success",
        ResponseCode:  fiber.StatusOK,
        Message:       message,
        Data:          data,
    })
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Status:        "Created",
        ResponseCode:  fiber.StatusCreated,
        Message:       message,
        Data:          data,
    })
}
```

### List Utils yang Tersedia:
- `response.go` - Fungsi helper untuk HTTP response formatting
- `jwt.go` - Fungsi untuk JWT token generation dan validation
- `password.go` - Fungsi untuk password hashing dan comparison
- `sorting_list_position.go` - Fungsi untuk sorting logic

---

## 🔄 Alur Data (Flow)

Contoh alur data untuk operasi **Create Board**:

```
┌─────────────┐
│   Client    │  (HTTP POST /boards)
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│                  CONTROLLER                              │
│  1. Parse request body                                   │
│  2. Validate input                                       │
│  3. Extract user from JWT token                          │
│  4. Call service                                         │
└──────┬──────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│                   SERVICE                                │
│  1. Validate business logic (owner exists?)              │
│  2. Generate new UUID                                    │
│  3. Map internal ID from public ID                       │
│  4. Call repository                                      │
└──────┬──────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│                 REPOSITORY                               │
│  1. Execute SQL query (INSERT INTO boards...)            │
│  2. Return result                                        │
└──────┬──────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│                  DATABASE                                │
│  1. Store data                                           │
│  2. Return success/error                                 │
└─────────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│                  CONTROLLER                              │
│  1. Format response using utils                          │
│  2. Return JSON response to client                       │
└──────┬──────────────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│   Client    │  (JSON Response)
└─────────────┘
```

---

## ✅ Keuntungan Arsitektur Ini

1. **Pemisahan Tanggung Jawab (Separation of Concerns)**
   - Setiap layer memiliki fokus yang jelas
   - Controller: HTTP handling
   - Service: Business logic
   - Repository: Data access

2. **Mudah Diuji (Testable)**
   - Setiap layer bisa diuji secara terpisah
   - Mock repository saat testing service
   - Mock service saat testing controller

3. **Mudah Dimaintain (Maintainable)**
   - Perubahan di satu layer tidak mempengaruhi layer lain
   - Kode lebih organized dan readable

4. **Reusable Code**
   - Utils dapat digunakan di berbagai bagian
   - Service dan repository dapat digunakan di multiple controllers

5. **Scalability**
   - Mudah menambahkan fitur baru
   - Mudah mengganti database implementation
   - Mudah menambahkan caching layer

---

## 📝 Pattern yang Digunakan

Pattern yang diterapkan dalam project ini adalah **Repository Pattern** dengan layering architecture. Ini adalah pattern yang sangat umum dan recommended untuk aplikasi API modern karena:

- Clean separation of concerns
- Easy testing
- Good for medium to large scale applications
- Follows SOLID principles, especially Single Responsibility Principle

---

*Last Updated: 2026-04-11*
