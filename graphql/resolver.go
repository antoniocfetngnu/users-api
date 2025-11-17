package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"strconv"
	"strings"

	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/models"
)

type Resolver struct{}

// Follower field resolvers
func (r *followerResolver) ID(ctx context.Context, obj *models.Follower) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *followerResolver) FollowerID(ctx context.Context, obj *models.Follower) (string, error) {
	return strconv.FormatUint(uint64(obj.FollowerID), 10), nil
}

func (r *followerResolver) FollowedID(ctx context.Context, obj *models.Follower) (string, error) {
	return strconv.FormatUint(uint64(obj.FollowedID), 10), nil
}

func (r *followerResolver) FollowedSince(ctx context.Context, obj *models.Follower) (string, error) {
	return obj.FollowedSince.Format("2006-01-02T15:04:05Z07:00"), nil
}

// Users resolver
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// User by ID resolver
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// User by username resolver
func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// User by email resolver
func (r *queryResolver) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Search users resolver
func (r *queryResolver) SearchUsers(ctx context.Context, query string) ([]*models.User, error) {
	var users []*models.User
	searchQuery := "%" + strings.ToLower(query) + "%"

	if err := database.DB.Where(
		"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(username) LIKE ?",
		searchQuery, searchQuery, searchQuery,
	).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Get all users that a specific user follows
func (r *queryResolver) Following(ctx context.Context, userID string) ([]*models.Follower, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	var followers []*models.Follower
	if err := database.DB.
		Preload("Follower").
		Preload("Followed").
		Where("follower_id = ?", id).
		Find(&followers).Error; err != nil {
		return nil, err
	}

	return followers, nil
}

// Get all followers of a specific user
func (r *queryResolver) Followers(ctx context.Context, userID string) ([]*models.Follower, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	var followers []*models.Follower
	if err := database.DB.
		Preload("Follower").
		Preload("Followed").
		Where("followed_id = ?", id).
		Find(&followers).Error; err != nil {
		return nil, err
	}

	return followers, nil
}

// Check if userA follows userB
func (r *queryResolver) IsFollowing(ctx context.Context, followerID string, followedID string) (bool, error) {
	fID, err := strconv.ParseUint(followerID, 10, 32)
	if err != nil {
		return false, err
	}

	fdID, err := strconv.ParseUint(followedID, 10, 32)
	if err != nil {
		return false, err
	}

	var count int64
	if err := database.DB.Model(&models.Follower{}).
		Where("follower_id = ? AND followed_id = ?", fID, fdID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Get follower relationship details
func (r *queryResolver) FollowerRelationship(ctx context.Context, followerID string, followedID string) (*models.Follower, error) {
	fID, err := strconv.ParseUint(followerID, 10, 32)
	if err != nil {
		return nil, err
	}

	fdID, err := strconv.ParseUint(followedID, 10, 32)
	if err != nil {
		return nil, err
	}

	var follower models.Follower
	if err := database.DB.
		Preload("Follower").
		Preload("Followed").
		Where("follower_id = ? AND followed_id = ?", fID, fdID).
		First(&follower).Error; err != nil {
		return nil, err
	}

	return &follower, nil
}

// Get follower count for a user
func (r *queryResolver) FollowerCount(ctx context.Context, userID string) (int, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := database.DB.Model(&models.Follower{}).
		Where("followed_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// Get following count for a user
func (r *queryResolver) FollowingCount(ctx context.Context, userID string) (int, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := database.DB.Model(&models.Follower{}).
		Where("follower_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// Field resolvers for User type
func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}

// Follower returns FollowerResolver implementation.
func (r *Resolver) Follower() FollowerResolver { return &followerResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type followerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
