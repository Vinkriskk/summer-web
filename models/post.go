package models

// Post schema for Post table
type Post struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Caption string `json:"caption" gorm:"not null"`
	UserID  uint   `json:"user_id" gorm:"not null"`
}
