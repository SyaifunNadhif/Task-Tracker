package controllers

import (
	"net/http"
	"task-tracker/config"
	"task-tracker/helpers"
	"task-tracker/inputs" // Pastikan folder inputs diimport
	"task-tracker/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Mendaftarkan User Baru
// @Description  Membuat akun user baru
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body inputs.RegisterInput true "Data Register"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Router       /register [post]
func Register(c *gin.Context) {
	// UBAH 1: Pakai struct dari inputs, bukan models.User langsung
	var input inputs.RegisterInput

	// 1. Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helpers.APIResponse("Register gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// 2. Pindahkan data dari Input ke Model Database (Manual Mapping)
	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// 3. Simpan User ke Database
	if err := config.DB.Create(&newUser).Error; err != nil {
		response := helpers.APIResponse("Register gagal (Email mungkin sudah ada)", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 4. Kembalikan Response Sukses
	newUser.Password = "********" // Sembunyikan password

	response := helpers.APIResponse("Akun berhasil dibuat!", http.StatusOK, "success", newUser)
	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary      Login User
// @Description  Masuk untuk mendapatkan Token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body inputs.LoginInput true "Login Credential"
// @Success      200  {object}  helpers.Response
// @Failure      401  {object}  helpers.Response
// @Router       /login [post]
func Login(c *gin.Context) {
	// UBAH 2: Pakai struct inputs.LoginInput (JANGAN pakai struct manual lagi)
	var input inputs.LoginInput

	// 1. Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helpers.APIResponse("Login gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// 2. Cari User pakai Email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response := helpers.APIResponse("Email atau Password salah", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// 3. Cek Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		response := helpers.APIResponse("Email atau Password salah", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// 4. Generate Token
	token, err := helpers.GenerateToken(user.ID)
	if err != nil {
		response := helpers.APIResponse("Gagal generate token", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 5. Kirim Token
	response := helpers.APIResponse("Login berhasil", http.StatusOK, "success", gin.H{
		"token": token,
		"name":  user.Name,
	})
	c.JSON(http.StatusOK, response)
}