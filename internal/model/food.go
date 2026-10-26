package model

// Food merepresentasikan data makanan di tabel public.makanan
type Food struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Calories     float64  `json:"calories"`     // Kalori per 100 gram
	Proteins     float64  `json:"proteins"`     // Protein per 100 gram
	Fat          float64  `json:"fat"`          // Lemak per 100 gram
	Carbohydrate float64  `json:"carbohydrate"` // Karbohidrat per 100 gram
	Image        *string  `json:"image"`
}

