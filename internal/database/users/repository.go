package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `INSERT INTO users (id, username, email, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, city, created_at, updated_at`
	u.ID = uuid.New()
	u.Username = req.Username
	u.Password = req.Password
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	row := r.userDB.QueryRow(ctx, query, u.ID, u.Username, u.Password, now, now)
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return database.User{}, err
	}

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, username, email, city, created_at, updated_at FROM users WHERE id = $1`
	row := r.userDB.QueryRow(ctx, query, userID)
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, nil
		}
		return database.User{}, err
	}

	return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, username, email, city, created_at, updated_at FROM users WHERE username = $1`
	row := r.userDB.QueryRow(ctx, query, username)
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return database.User{}, nil
		}
		return database.User{}, err
	}

	return u, nil
}
