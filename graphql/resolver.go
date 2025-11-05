package graphql

import (
	"context"
	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/models"
)

type Resolver struct{}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	result, err := database.DB.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", 
		input.FirstName, input.LastName, input.Email)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	
	var user models.User
	err = database.DB.QueryRow("SELECT id, first_name, last_name, email, created_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Health is the resolver for the health field.
func (r *queryResolver) Health(ctx context.Context) (string, error) {
	return "GraphQL API is healthy!", nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	rows, err := database.DB.Query("SELECT id, first_name, last_name, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, first_name, last_name, email, created_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAt.Format("2006-01-02 15:04:05"), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }