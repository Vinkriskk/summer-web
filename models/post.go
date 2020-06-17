package models

// Post schema for Post table
type Post struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Caption string `json:"caption"`
	UserID  int64  `json:"user_id"`
}
