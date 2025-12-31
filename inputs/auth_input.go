package inputs  // <--- Pastikan nama package-nya "inputs" (bukan input, bukan main)

// Huruf 'L' harus BESAR (Capital) biar bisa dibaca dari luar
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Huruf 'R' harus BESAR
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}