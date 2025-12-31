package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" binding:"required"`
	Email     string         `json:"email" gorm:"unique" binding:"required,email"`
	Password  string         `json:"password" binding:"required,min=6"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Tambahan Baru: Fungsi untuk Mengacak Password sebelum Save
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// Jika password sudah di-hash (panjangnya > 50), skip saja
    // Ini jaga-jaga biar gak kena hash dua kali
	if len(u.Password) > 50 { 
		return
	}
    
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}