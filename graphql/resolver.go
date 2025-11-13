package graphql

import (
	"context"
	"strconv"
	"strings"

	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/models"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

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

// Field resolvers for User type
func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}

func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}
