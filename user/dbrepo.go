package user

import (
	"context"
	"fmt"

	"github.com/gvre/api-sample-app/app"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DatabaseRepository implements the UserRepository interface.
type DatabaseRepository struct {
	db *pgxpool.Pool
}

// NewDatabaseRepository returns a new DatabaseRepository.
func NewDatabaseRepository(db *pgxpool.Pool) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// Ping checks if database is up.
func (repo *DatabaseRepository) Ping(ctx context.Context) error {
	_, err := repo.db.Query(ctx, "SELECT 1")

	return err
}

// FetchAll returns all users.
func (repo *DatabaseRepository) FetchAll(ctx context.Context) ([]app.User, error) {
	users := []app.User{}
	rows, err := repo.db.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, app.NewError(
			"Error while fetching users",
			fmt.Errorf("query: %w", err),
		)
	}

	for rows.Next() {
		var u app.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, app.NewError(
				"Error while reading users",
				fmt.Errorf("scan: %w", err),
			)
		}

		users = append(users, u)
	}

	return users, nil
}

// FetchByID returns the user with the provided id.
// A nil value is returned if the user does not exist.
func (repo *DatabaseRepository) FetchByID(ctx context.Context, id int) (*app.User, error) {
	var (
		uid   int
		uname string
	)

	err := repo.db.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id).Scan(&uid, &uname)
	switch err {
	case nil:
		return &app.User{ID: uid, Name: uname}, nil
	case pgx.ErrNoRows:
		return nil, app.NewError("User does not exist", app.ErrNoRecords)
	default:
		return nil, app.NewError("Error while fetching user", fmt.Errorf("fetch by id: %w", err))
	}
}

// Add inserts a new user into the database and returns its id.
func (repo *DatabaseRepository) Add(ctx context.Context, name string) (int, error) {
	var id int
	err := repo.db.QueryRow(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, app.NewError(
			"Error while adding user",
			fmt.Errorf("add: %w", err),
		)
	}

	return id, nil
}
