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

func SetupRouter() *gin.Engine {
	r := gin.Default()

    // Config CORS (biarkan sama)
    r.Use(cors.Default())

	// Route Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth Routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Task Routes
	taskRoutes := r.Group("/tasks")
	taskRoutes.Use(middlewares.AuthMiddleware())
	{
		taskRoutes.GET("/", controllers.FindTasks)
		taskRoutes.POST("/", controllers.CreateTask)
		taskRoutes.GET("/:id", controllers.FindTaskByID)
		taskRoutes.PUT("/:id", controllers.UpdateTask)
		taskRoutes.DELETE("/:id", controllers.DeleteTask)
	}

    // Route Seeding
	r.POST("/seed", middlewares.AuthMiddleware(), controllers.SeedTasks)

	return r
}

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
// 2. MAIN FUNCTION JADI LEBIH BERSIH
func main() {
	config.ConnectDatabase() // Tetap connect DB di main
	
    r := SetupRouter() // Panggil fungsi setup tadi

	r.Run("localhost:8081")
}
