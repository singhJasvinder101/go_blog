package types

import "time"

// db tag for pgxpool, sqlx packages (DB -> Struct mapping)
// json tag for json marshal/unmarshal (Struct -> JSON mapping)
// binding tag for gin binding (validation)

type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"  binding:"required"`
	Email        string    `json:"email" db:"email"  binding:"required,email"`
	// json:"-" → Don't expose password hash
	PasswordHash string    `json:"-" db:"password_hash"` 
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Post struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"  binding:"required"`
	Description string    `json:"description" db:"description"  binding:"required"`
	UserID      int       `json:"user_id" db:"user_id"  binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"  binding:"required"`
	UserID    int       `json:"user_id" db:"user_id"  binding:"required"`
	PostID    int       `json:"post_id" db:"post_id"  binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Like struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"  binding:"required"`
	PostID    int       `json:"post_id" db:"post_id"  binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"  binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at" `
}

type PostTag struct {
	ID        int       `json:"id" db:"id" `
	PostID    int       `json:"post_id" db:"post_id"  binding:"required"`
	TagID     int       `json:"tag_id" db:"tag_id"  binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at" `
}



type APIResponse struct {
    Status string      `json:"status"`
    Data   interface{} `json:"data,omitempty"`
    Error  string      `json:"error,omitempty"`
}

