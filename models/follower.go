package models

import (
	"time"

	"gorm.io/gorm"
)

type Follower struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	FollowerID    uint           `gorm:"not null;index:idx_follower_followed" json:"followerId"` // User who follows
	FollowedID    uint           `gorm:"not null;index:idx_follower_followed" json:"followedId"` // User being followed
	FollowedSince time.Time      `gorm:"not null" json:"followedSince"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Follower User `gorm:"foreignKey:FollowerID" json:"follower"`
	Followed User `gorm:"foreignKey:FollowedID" json:"followed"`
}

// FollowRequest DTO
type FollowRequest struct {
	FollowedID uint `json:"followedId" binding:"required"`
}

// FollowerResponse for API responses
type FollowerResponse struct {
	ID            uint         `json:"id"`
	FollowerID    uint         `json:"followerId"`
	FollowedID    uint         `json:"followedId"`
	FollowedSince time.Time    `json:"followedSince"`
	Follower      UserResponse `json:"follower"`
	Followed      UserResponse `json:"followed"`
}

func (f *Follower) ToResponse() FollowerResponse {
	return FollowerResponse{
		ID:            f.ID,
		FollowerID:    f.FollowerID,
		FollowedID:    f.FollowedID,
		FollowedSince: f.FollowedSince,
		Follower:      f.Follower.ToResponse(),
		Followed:      f.Followed.ToResponse(),
	}
}
