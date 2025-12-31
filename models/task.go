package models

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status"`
	
	UserID      uint   `json:"user_id"`
	
	// Perhatikan perubahannya di sini: json:"-"
	// Artinya: "Jangan tampilkan field ini di output JSON, tapi tetap simpan di database"
	User        User   `json:"-" binding:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}