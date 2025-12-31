package inputs

type TaskInput struct {
	Title       string `json:"title" binding:"required" example:"Belajar Golang"`
	Description string `json:"description" binding:"required" example:"Mencoba koneksi database"`
	Status      string `json:"status" example:"Pending"`
}