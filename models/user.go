package models

import (
	"time"
)

// User schema for User table
type User struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	Username       string     `json:"username" gorm:"unique;not null"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Password       string     `json:"password,omitempty"`
	FollowerCount  int        `json:"follower_count"`
	FollowingCount int        `json:"following_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
