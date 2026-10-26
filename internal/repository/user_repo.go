package repository

import (
	"context"
	"errors"
	"go_vibe/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser menyimpan user baru ke database
func (r *UserRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (*model.User, error) {
	query := `
		INSERT INTO public.users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at;
	`
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Buat target nutrisi default untuk user baru
	initTargetQuery := `
		INSERT INTO public.target_nutrisi (user_id, calories_target, proteins_target, carbohydrate_target, fat_target)
		VALUES ($1, 2000, 50, 250, 70)
		ON CONFLICT (user_id) DO NOTHING;
	`
	_, _ = r.db.Exec(ctx, initTargetQuery, user.ID)

	return user, nil
}

// GetUserByEmail mencari user berdasarkan email (termasuk password_hash untuk verifikasi login)
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM public.users
		WHERE LOWER(email) = LOWER($1);
	`
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID mencari user berdasarkan ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM public.users
		WHERE id = $1;
	`
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

