package main

import (
	"task-tracker/config"
	"task-tracker/controllers"
	"task-tracker/middlewares"

	// --- TAMBAHAN IMPORT SWAGGER ---
	_ "task-tracker/docs" // <--- Penting: Import folder docs yang baru dibuat
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
    // -------------------------------

	"github.com/gin-contrib/cors" // <--- Import Library CORS
	"github.com/gin-gonic/gin"
)

// @title           Task Tracker API
// @version         1.0
// @description     Aplikasi Task Tracker dengan Golang GIN
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.ConnectDatabase()
	r := gin.Default()

	// --- PASANG CORS DI SINI ---
	// Kita atur supaya SEMUA orang boleh akses (Mode Development)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	// --- ROUTE SWAGGER ---
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Penting: Izinkan Header Authorization (biar token JWT bisa lewat)
	corsConfig.AddAllowHeaders("Authorization")

	// 3. Setup Routes
	//User
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	//proteced
	protected := r.Group("/")
    protected.Use(middlewares.AuthMiddleware()) // <--- Pasang Satpam
    {
        // Masukkan semua route Task ke sini!
        protected.GET("/tasks", controllers.FindTasks)
        protected.POST("/tasks", controllers.CreateTask) // <--- Ini yang penting
		protected.GET("/tasks/:id", controllers.FindTaskByID)
		protected.PUT("/tasks/:id", controllers.UpdateTask)
		protected.DELETE("/tasks/:id", controllers.DeleteTask)

		protected.POST("/seed", controllers.SeedTasks)
	}


	// 4. Jalankan Server
	r.Run("localhost:8081") // Kita pakai port 8081 seperti tadi
}
