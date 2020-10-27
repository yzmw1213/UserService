package model

// Auth 認証情報の構造体
type Auth struct {
	Token     string
	UserID    uint32
	Authority uint32
}
