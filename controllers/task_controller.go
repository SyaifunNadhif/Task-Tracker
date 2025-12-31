package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"task-tracker/config"
	"task-tracker/helpers"
	_ "task-tracker/inputs" // <--- PENTING: Import ini wajib ada biar Swagger jalan
	"task-tracker/models"

	"github.com/gin-gonic/gin"
)

// FindTasks godoc
// @Summary      Melihat Daftar Tugas
// @Description  Mengambil semua tugas milik user yang sedang login (pagination)
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        page   query     int  false  "Halaman ke berapa"
// @Param        limit  query     int  false  "Jumlah data per halaman"
// @Success      200    {object}  helpers.Response
// @Failure      401    {object}  helpers.Response
// @Router       /tasks [get]
func FindTasks(c *gin.Context) {
	userID, _ := c.Get("currentUser")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	offset := (page - 1) * limit

	var tasks []models.Task
	var total int64

	query := config.DB.Model(&models.Task{}).Where("user_id = ?", userID)

	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&tasks)

	data := gin.H{
		"tasks": tasks,
		"paging": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_items":  total,
			"total_pages":  math.Ceil(float64(total) / float64(limit)),
		},
	}

	response := helpers.APIResponse("List of tasks", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// CreateTask godoc
// @Summary      Membuat Tugas Baru
// @Description  Menambahkan data tugas baru ke database
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body inputs.TaskInput true "Data Task"
// @Success      200     {object}  helpers.Response
// @Failure      401     {object}  helpers.Response
// @Failure      422     {object}  helpers.Response
// @Router       /tasks [post]
func CreateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helpers.APIResponse("Input tidak valid", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		response := helpers.APIResponse("User tidak terautentikasi", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		UserID:      currentUser.(uint),
	}

	if err := config.DB.Create(&task).Error; err != nil {
		response := helpers.APIResponse("Gagal menyimpan ke Database", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Task berhasil dibuat", http.StatusOK, "success", task)
	c.JSON(http.StatusOK, response)
}

// FindTaskByID godoc
// @Summary      Detail Tugas
// @Description  Melihat detail satu tugas berdasarkan ID
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "ID Tugas"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Router       /tasks/{id} [get]
func FindTaskByID(c *gin.Context) {
	var task models.Task
	if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		response := helpers.APIResponse("Task tidak ditemukan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Detail task", http.StatusOK, "success", task)
	c.JSON(http.StatusOK, response)
}

// UpdateTask godoc
// @Summary      Update Tugas
// @Description  Mengubah data tugas (Judul, Deskripsi, Status)
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "ID Tugas"
// @Param        request  body      inputs.TaskInput  true  "Data Update"
// @Success      200      {object}  helpers.Response
// @Failure      404      {object}  helpers.Response
// @Failure      422      {object}  helpers.Response
// @Router       /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	userID, _ := c.Get("currentUser")

	var task models.Task
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		response := helpers.APIResponse("Task tidak ditemukan atau bukan milik Anda", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helpers.APIResponse("Input tidak valid", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	config.DB.Model(&task).Updates(input)

	response := helpers.APIResponse("Task berhasil diupdate", http.StatusOK, "success", task)
	c.JSON(http.StatusOK, response)
}

// DeleteTask godoc
// @Summary      Hapus Tugas
// @Description  Menghapus tugas berdasarkan ID secara permanen
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "ID Tugas"
// @Success      200  {object}  helpers.Response
// @Failure      404  {object}  helpers.Response
// @Router       /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	userID, _ := c.Get("currentUser")

	var task models.Task
	if err := config.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		response := helpers.APIResponse("Task tidak ditemukan atau bukan milik Anda", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	config.DB.Delete(&task)

	response := helpers.APIResponse("Task berhasil dihapus", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// SeedTasks godoc
// @Summary      Generate Dummy Data
// @Description  Membuat 50 data tugas palsu secara otomatis
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  helpers.Response
// @Router       /seed [post]
func SeedTasks(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	uid := userID.(uint)

	for i := 1; i <= 50; i++ {
		title := fmt.Sprintf("Tugas Penting ke-%d", i)
		desc := fmt.Sprintf("Ini adalah deskripsi untuk tugas nomor %d. Dibuat otomatis oleh robot.", i)

		status := "Pending"
		if i%2 == 0 {
			status = "Completed"
		}

		task := models.Task{
			Title:       title,
			Description: desc,
			Status:      status,
			UserID:      uid,
		}

		config.DB.Create(&task)
	}

	response := helpers.APIResponse("Berhasil menanam 50 data dummy!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}










// GET All Tasks
// func FindTasks(c *gin.Context) {
// 	var tasks []models.Task
// 	config.DB.Find(&tasks)

// 	// Pakai Helper APIResponse
// 	response := helpers.APIResponse("List of tasks", http.StatusOK, "success", tasks)
// 	c.JSON(http.StatusOK, response)
// }

// GET All Tasks 2
// func FindTasks(c *gin.Context) {
// 	// 1. Ambil ID User yang login
// 	userID, _ := c.Get("currentUser")

// 	var tasks []models.Task

// 	// 2. Filter: SELECT * FROM tasks WHERE user_id = ...
// 	config.DB.Where("user_id = ?", userID).Find(&tasks)

// 	response := helpers.APIResponse("List of tasks", http.StatusOK, "success", tasks)
// 	c.JSON(http.StatusOK, response)
// }

// CREATE Task
// func CreateTask(c *gin.Context) {
// 	var input models.Task
//     // Jika JSON salah format
// 	if err := c.ShouldBindJSON(&input); err != nil {

//         // Format Error Validasi biar rapi
// 		errors := helpers.FormatValidationError(err)
// 		errorMessage := gin.H{"errors": errors}

// 		response := helpers.APIResponse("Gagal membuat task", http.StatusUnprocessableEntity, "error", errorMessage)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}

// 	task := models.Task{Title: input.Title, Description: input.Description, Status: input.Status}
// 	config.DB.Create(&task)

// 	response := helpers.APIResponse("Task berhasil dibuat", http.StatusOK, "success", task)
// 	c.JSON(http.StatusOK, response)
// }



// UPDATE Task
// func UpdateTask(c *gin.Context) {
// 	var task models.Task
// 	if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
// 		response := helpers.APIResponse("Task tidak ditemukan", http.StatusBadRequest, "error", nil)
// 		c.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	var input models.Task
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		errors := helpers.FormatValidationError(err)
// 		errorMessage := gin.H{"errors": errors}

// 		response := helpers.APIResponse("Gagal update task", http.StatusUnprocessableEntity, "error", errorMessage)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}

// 	config.DB.Model(&task).Updates(input)
// 	response := helpers.APIResponse("Task berhasil diupdate", http.StatusOK, "success", task)
// 	c.JSON(http.StatusOK, response)
// }

// DELETE Task
// func DeleteTask(c *gin.Context) {
// 	var task models.Task
// 	if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
// 		response := helpers.APIResponse("Task tidak ditemukan", http.StatusBadRequest, "error", nil)
// 		c.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	config.DB.Delete(&task)
// 	response := helpers.APIResponse("Task berhasil dihapus", http.StatusOK, "success", nil)
// 	c.JSON(http.StatusOK, response)
// }

