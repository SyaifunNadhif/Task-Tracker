package models

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
    // Perhatikan bagian binding:"..."
	Title       string `json:"title" binding:"required,min=3"` 
	Description string `json:"description" binding:"required"`
	Status      string `json:"status"` // Status gak wajib, boleh kosong
}