package controllers

import (
    "net/http"
    "task-tracker/config"
    "task-tracker/models"

    "github.com/gin-gonic/gin"
)

// GET /tasks
func FindTasks(c *gin.Context) {
    var tasks []models.Task
    config.DB.Find(&tasks) // Select * from tasks

    c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// POST /tasks
func CreateTask(c *gin.Context) {
    var input models.Task
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Simpan ke database
    task := models.Task{Title: input.Title, Description: input.Description, Status: input.Status}
    config.DB.Create(&task)

    c.JSON(http.StatusOK, gin.H{"data": task})
}

// GET /tasks/:id
func FindTaskByID(c *gin.Context) {
    var task models.Task
    
    // Cek apakah ID ada di database
    if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": task})
}

// ... fungsi FindTasks, CreateTask, FindTaskByID ada di atas sini ...

// UPDATE /tasks/:id
func UpdateTask(c *gin.Context) {
    var task models.Task
    // 1. Cari dulu datanya, ada gak?
    if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found!"})
        return
    }

    // 2. Validasi input JSON dari user
    var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // 3. Update data
    // GORM akan update kolom yang dikirim saja (Title, Description, Status)
    config.DB.Model(&task).Updates(input)

    c.JSON(http.StatusOK, gin.H{"data": task})
}

// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
    var task models.Task
    // 1. Cari dulu datanya
    if err := config.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found!"})
        return
    }

    // 2. Hapus data (Soft Delete kalau pakai gorm.Model, Hard Delete kalau struct biasa)
    config.DB.Delete(&task)

    c.JSON(http.StatusOK, gin.H{"data": true, "message": "Task berhasil dihapus!"})
}