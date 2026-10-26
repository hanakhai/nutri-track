package repository

import (
	"context"
	"errors"
	"go_vibe/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TargetRepository struct {
	db *pgxpool.Pool
}

func NewTargetRepository(db *pgxpool.Pool) *TargetRepository {
	return &TargetRepository{db: db}
}

// GetTargetByUserID mengambil target nutrisi pengguna
func (r *TargetRepository) GetTargetByUserID(ctx context.Context, userID int) (*model.NutritionTarget, error) {
	query := `
		SELECT id, user_id, calories_target, proteins_target, carbohydrate_target, fat_target, updated_at
		FROM public.target_nutrisi
		WHERE user_id = $1;
	`
	t := &model.NutritionTarget{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.CaloriesTarget,
		&t.ProteinsTarget,
		&t.CarbohydrateTarget,
		&t.FatTarget,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Jika belum diset, buatkan target default
			insertQuery := `
				INSERT INTO public.target_nutrisi (user_id, calories_target, proteins_target, carbohydrate_target, fat_target)
				VALUES ($1, 2000, 50, 250, 70)
				RETURNING id, user_id, calories_target, proteins_target, carbohydrate_target, fat_target, updated_at;
			`
			err = r.db.QueryRow(ctx, insertQuery, userID).Scan(
				&t.ID,
				&t.UserID,
				&t.CaloriesTarget,
				&t.ProteinsTarget,
				&t.CarbohydrateTarget,
				&t.FatTarget,
				&t.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			return t, nil
		}
		return nil, err
	}
	return t, nil
}

// UpsertTarget menyimpan atau memperbarui target nutrisi pengguna
func (r *TargetRepository) UpsertTarget(ctx context.Context, userID int, req model.UpdateTargetRequest) (*model.NutritionTarget, error) {
	query := `
		INSERT INTO public.target_nutrisi (user_id, calories_target, proteins_target, carbohydrate_target, fat_target, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET
			calories_target = EXCLUDED.calories_target,
			proteins_target = EXCLUDED.proteins_target,
			carbohydrate_target = EXCLUDED.carbohydrate_target,
			fat_target = EXCLUDED.fat_target,
			updated_at = NOW()
		RETURNING id, user_id, calories_target, proteins_target, carbohydrate_target, fat_target, updated_at;
	`
	t := &model.NutritionTarget{}
	err := r.db.QueryRow(ctx, query, userID, req.CaloriesTarget, req.ProteinsTarget, req.CarbohydrateTarget, req.FatTarget).Scan(
		&t.ID,
		&t.UserID,
		&t.CaloriesTarget,
		&t.ProteinsTarget,
		&t.CarbohydrateTarget,
		&t.FatTarget,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

