package repository

import (
	"context"
	"errors"
	"fmt"
	"go_vibe/internal/model"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogRepository struct {
	db         *pgxpool.Pool
	targetRepo *TargetRepository
}

func NewLogRepository(db *pgxpool.Pool, targetRepo *TargetRepository) *LogRepository {
	return &LogRepository{db: db, targetRepo: targetRepo}
}

func round(val float64) float64 {
	return math.Round(val*10) / 10
}

// GetSessionsByDate mengambil seluruh sesi makan pengguna pada tanggal tertentu
func (r *LogRepository) GetSessionsByDate(ctx context.Context, userID int, dateStr string) ([]model.MealSession, error) {
	sessionQuery := `
		SELECT id, user_id, name, TO_CHAR(meal_date, 'YYYY-MM-DD') as meal_date, meal_time, note, created_at
		FROM public.sesi_makan
		WHERE user_id = $1 AND meal_date = $2::date
		ORDER BY id ASC;
	`
	rows, err := r.db.Query(ctx, sessionQuery, userID, dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.MealSession
	for rows.Next() {
		var s model.MealSession
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Name,
			&s.MealDate,
			&s.MealTime,
			&s.Note,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.Items = []model.MealItem{}
		sessions = append(sessions, s)
	}

	// Ambil item untuk setiap sesi
	for i := range sessions {
		itemQuery := `
			SELECT 
				ip.id, 
				ip.sesi_makan_id, 
				ip.makanan_id, 
				m.name, 
				ip.gram, 
				COALESCE(ip.serving_unit, 'Porsi'), 
				ip.serving_note,
				(ip.gram / 100.0 * m.calories) as calories,
				(ip.gram / 100.0 * m.proteins) as proteins,
				(ip.gram / 100.0 * m.fat) as fat,
				(ip.gram / 100.0 * m.carbohydrate) as carbohydrate
			FROM public.isi_piring ip
			JOIN public.makanan m ON ip.makanan_id = m.id
			WHERE ip.sesi_makan_id = $1
			ORDER BY ip.id ASC;
		`
		itemRows, err := r.db.Query(ctx, itemQuery, sessions[i].ID)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {
			var item model.MealItem
			err := itemRows.Scan(
				&item.ID,
				&item.SesiMakanID,
				&item.MakananID,
				&item.MakananName,
				&item.Gram,
				&item.ServingUnit,
				&item.ServingNote,
				&item.Calories,
				&item.Proteins,
				&item.Fat,
				&item.Carbohydrate,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}

			// Bulatkan nilai nutrisi
			item.Calories = round(item.Calories)
			item.Proteins = round(item.Proteins)
			item.Fat = round(item.Fat)
			item.Carbohydrate = round(item.Carbohydrate)

			sessions[i].SessionCalories += item.Calories
			sessions[i].SessionProteins += item.Proteins
			sessions[i].SessionFat += item.Fat
			sessions[i].SessionCarbohydrate += item.Carbohydrate

			sessions[i].Items = append(sessions[i].Items, item)
		}
		itemRows.Close()

		sessions[i].SessionCalories = round(sessions[i].SessionCalories)
		sessions[i].SessionProteins = round(sessions[i].SessionProteins)
		sessions[i].SessionFat = round(sessions[i].SessionFat)
		sessions[i].SessionCarbohydrate = round(sessions[i].SessionCarbohydrate)
	}

	return sessions, nil
}

// CreateMealSession membuat sesi makan baru beserta daftar makanannya dalam satu transaksi
func (r *LogRepository) CreateMealSession(ctx context.Context, userID int, req model.CreateMealSessionRequest) (*model.MealSession, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	sessionName := req.Name
	if sessionName == "" {
		sessionName = fmt.Sprintf("Makan %s", req.MealTime)
	}

	sessionQuery := `
		INSERT INTO public.sesi_makan (user_id, name, meal_date, meal_time, note)
		VALUES ($1, $2, $3::date, $4, $5)
		RETURNING id, user_id, name, TO_CHAR(meal_date, 'YYYY-MM-DD'), meal_time, note, created_at;
	`
	session := &model.MealSession{}
	err = tx.QueryRow(ctx, sessionQuery, userID, sessionName, req.MealDate, req.MealTime, req.Note).Scan(
		&session.ID,
		&session.UserID,
		&session.Name,
		&session.MealDate,
		&session.MealTime,
		&session.Note,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat sesi makan: %w", err)
	}

	// Masukkan item ke isi_piring
	for _, item := range req.Items {
		unit := item.ServingUnit
		if unit == "" {
			unit = "Porsi"
		}

		insertItemQuery := `
			INSERT INTO public.isi_piring (sesi_makan_id, makanan_id, gram, serving_unit, serving_note)
			VALUES ($1, $2, $3, $4, $5);
		`
		_, err = tx.Exec(ctx, insertItemQuery, session.ID, item.MakananID, item.Gram, unit, item.ServingNote)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan item makanan: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Muat kembali sesi dengan rincian lengkap
	sessions, err := r.GetSessionsByDate(ctx, userID, req.MealDate)
	if err != nil {
		return session, nil
	}
	for _, s := range sessions {
		if s.ID == session.ID {
			return &s, nil
		}
	}

	return session, nil
}

// AddMealItem menambahkan satu makanan ke sesi makan yang sudah ada
func (r *LogRepository) AddMealItem(ctx context.Context, userID int, sessionID int, req model.AddMealItemRequest) error {
	// Verifikasi kepemilikan sesi makan
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM public.sesi_makan WHERE id = $1 AND user_id = $2);`, sessionID, userID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("sesi makan tidak ditemukan atau Anda tidak memiliki akses")
	}

	unit := req.ServingUnit
	if unit == "" {
		unit = "Porsi"
	}

	query := `
		INSERT INTO public.isi_piring (sesi_makan_id, makanan_id, gram, serving_unit, serving_note)
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err = r.db.Exec(ctx, query, sessionID, req.MakananID, req.Gram, unit, req.ServingNote)
	return err
}

// DeleteMealSession menghapus satu sesi makan milik pengguna
func (r *LogRepository) DeleteMealSession(ctx context.Context, userID int, sessionID int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM public.sesi_makan WHERE id = $1 AND user_id = $2;`, sessionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("sesi makan tidak ditemukan atau Anda tidak memiliki akses")
	}
	return nil
}

// DeleteMealItem menghapus satu item makanan dari isi piring milik pengguna
func (r *LogRepository) DeleteMealItem(ctx context.Context, userID int, itemID int) error {
	query := `
		DELETE FROM public.isi_piring
		WHERE id = $1 
		  AND sesi_makan_id IN (SELECT id FROM public.sesi_makan WHERE user_id = $2);
	`
	tag, err := r.db.Exec(ctx, query, itemID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("item makanan tidak ditemukan atau Anda tidak memiliki akses")
	}
	return nil
}

// GetDailySummary menghitung total nutrisi dan membandingkannya dengan target harian
func (r *LogRepository) GetDailySummary(ctx context.Context, userID int, dateStr string) (*model.DailyNutritionSummary, error) {
	sessions, err := r.GetSessionsByDate(ctx, userID, dateStr)
	if err != nil {
		return nil, err
	}

	summary := &model.DailyNutritionSummary{
		Date:     dateStr,
		Sessions: sessions,
	}

	for _, s := range sessions {
		summary.TotalCalories += s.SessionCalories
		summary.TotalProteins += s.SessionProteins
		summary.TotalFat += s.SessionFat
		summary.TotalCarb += s.SessionCarbohydrate
	}

	summary.TotalCalories = round(summary.TotalCalories)
	summary.TotalProteins = round(summary.TotalProteins)
	summary.TotalFat = round(summary.TotalFat)
	summary.TotalCarb = round(summary.TotalCarb)

	// Ambil target pengguna
	target, err := r.targetRepo.GetTargetByUserID(ctx, userID)
	if err == nil {
		summary.Target = target
	}

	return summary, nil
}
