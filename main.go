package main

import (
	"task-tracker/config"
	"task-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Konek Database
	config.ConnectDatabase()

	// 2. Init Router
	r := gin.Default()

	// 3. Setup Routes
	r.GET("/tasks", controllers.FindTasks)
	r.POST("/tasks", controllers.CreateTask)
	r.GET("/tasks/:id", controllers.FindTaskByID)
	// --- TAMBAHAN BARU ---
	r.PUT("/tasks/:id", controllers.UpdateTask)    // Buat Edit
	r.DELETE("/tasks/:id", controllers.DeleteTask) // Buat Hapus

	// 4. Jalankan Server
	r.Run("localhost:8081") // Kita pakai port 8081 seperti tadi
}
