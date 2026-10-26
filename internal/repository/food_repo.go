package repository

import (
	"context"
	"errors"
	"go_vibe/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FoodRepository struct {
	db *pgxpool.Pool
}

func NewFoodRepository(db *pgxpool.Pool) *FoodRepository {
	return &FoodRepository{db: db}
}

// SearchFoods mencari makanan berdasarkan nama, diurutkan berdasarkan kecocokan
func (r *FoodRepository) SearchFoods(ctx context.Context, keyword string, limit int) ([]model.Food, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var rows pgx.Rows
	var err error

	if keyword == "" {
		query := `
			SELECT id, name, calories, proteins, fat, carbohydrate, image
			FROM public.makanan
			ORDER BY name ASC
			LIMIT $1;
		`
		rows, err = r.db.Query(ctx, query, limit)
	} else {
		query := `
			SELECT id, name, calories, proteins, fat, carbohydrate, image
			FROM public.makanan
			WHERE name ILIKE '%' || $1 || '%'
			ORDER BY 
				CASE WHEN name ILIKE $1 || '%' THEN 1 ELSE 2 END,
				name ASC
			LIMIT $2;
		`
		rows, err = r.db.Query(ctx, query, keyword, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []model.Food
	for rows.Next() {
		var f model.Food
		err := rows.Scan(
			&f.ID,
			&f.Name,
			&f.Calories,
			&f.Proteins,
			&f.Fat,
			&f.Carbohydrate,
			&f.Image,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, f)
	}

	return foods, nil
}

// GetFoodByID mengambil detail satu makanan berdasarkan ID
func (r *FoodRepository) GetFoodByID(ctx context.Context, id int) (*model.Food, error) {
	query := `
		SELECT id, name, calories, proteins, fat, carbohydrate, image
		FROM public.makanan
		WHERE id = $1;
	`
	f := &model.Food{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.Name,
		&f.Calories,
		&f.Proteins,
		&f.Fat,
		&f.Carbohydrate,
		&f.Image,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

