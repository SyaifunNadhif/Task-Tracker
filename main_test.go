package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"task-tracker/config"
	"task-tracker/models"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
)

// Variabel Global untuk menyimpan Token dan ID Task antar test
var authToken string
var createdTaskID uint
var testEmail = "robot_tester@example.com" // Email khusus testing

// Fungsi Init untuk konek DB sebelum test jalan
func init() {
	config.ConnectDatabase()
}

// Fungsi bantu untuk membersihkan data sisa test sebelumnya (biar gak error "Email already exists")
func cleanUpData() {
	config.DB.Unscoped().Where("email = ?", testEmail).Delete(&models.User{})
	// Kita juga bisa hapus task yg dibuat test user ini kalau mau bersih total
	// Tapi hapus user biasanya sudah cukup karena relationnya
}

func TestFullWorkflow(t *testing.T) {
	// 0. Bersihkan Database dulu
	cleanUpData()

	router := SetupRouter()

	// ----------------------------------------------------------------
	// STEP 1: REGISTER USER BARU
	// ----------------------------------------------------------------
	t.Run("1. Register User", func(t *testing.T) {
		registerData := map[string]string{
			"name":     "Robot Tester",
			"email":    testEmail,
			"password": "passwordrahasia",
		}
		jsonValue, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// ----------------------------------------------------------------
	// STEP 2: LOGIN & AMBIL TOKEN
	// ----------------------------------------------------------------
	t.Run("2. Login User", func(t *testing.T) {
		loginData := map[string]string{
			"email":    testEmail,
			"password": "passwordrahasia",
		}
		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// PARSING JSON RESPONSE BUAT AMBIL TOKEN
		// Struktur response: { "data": { "token": "...", "name": "..." } }
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		dataMap := response["data"].(map[string]interface{})
		authToken = dataMap["token"].(string) // SIMPAN TOKEN KE VARIABEL GLOBAL

		fmt.Println(">> Token didapatkan:", authToken) // Debugging
	})

	// ----------------------------------------------------------------
	// STEP 3: CREATE TASK (Pakai Token)
	// ----------------------------------------------------------------
	t.Run("3. Create Task", func(t *testing.T) {
		taskData := map[string]string{
			"title":       "Testing Golang",
			"description": "Dibuat otomatis oleh robot test",
			"status":      "Pending",
		}
		jsonValue, _ := json.Marshal(taskData)
		req, _ := http.NewRequest("POST", "/tasks/", bytes.NewBuffer(jsonValue)) // Perhatikan slash di akhir kadang ngaruh di Gin group
		
		// PENTING: Pasang Token di Header
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// AMBIL ID TASK YANG BARU DIBUAT
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		dataMap := response["data"].(map[string]interface{})
		
		// JSON number di Go kadang dibaca float64, jadi perlu convert
		idFloat := dataMap["id"].(float64) 
		createdTaskID = uint(idFloat) // SIMPAN ID KE VARIABEL GLOBAL

		fmt.Println(">> Task Created ID:", createdTaskID)
	})

	// ----------------------------------------------------------------
	// STEP 4: GET ALL TASKS
	// ----------------------------------------------------------------
	t.Run("4. Get All Tasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks/", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Testing Golang") // Pastikan data tadi muncul
	})

	// ----------------------------------------------------------------
	// STEP 5: GET TASK BY ID
	// ----------------------------------------------------------------
	t.Run("5. Get Task By ID", func(t *testing.T) {
		url := "/tasks/" + strconv.Itoa(int(createdTaskID))
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Testing Golang")
	})

	// ----------------------------------------------------------------
	// STEP 6: UPDATE TASK
	// ----------------------------------------------------------------
	t.Run("6. Update Task", func(t *testing.T) {
		updateData := map[string]string{
			"title":       "Testing Golang UPDATED", // Ganti judul
			"description": "Deskripsi baru nih",
			"status":      "Completed", // Ganti status
		}
		jsonValue, _ := json.Marshal(updateData)

		url := "/tasks/" + strconv.Itoa(int(createdTaskID))
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UPDATED")
	})

	// ----------------------------------------------------------------
	// STEP 7: DELETE TASK
	// ----------------------------------------------------------------
	t.Run("7. Delete Task", func(t *testing.T) {
		url := "/tasks/" + strconv.Itoa(int(createdTaskID))
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// ----------------------------------------------------------------
	// STEP 8: VERIFIKASI DELETE (HARUS 400/404)
	// ----------------------------------------------------------------
	t.Run("8. Verify Delete", func(t *testing.T) {
		url := "/tasks/" + strconv.Itoa(int(createdTaskID))
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Karena di controller kamu FindTaskByID return 400 kalau not found
		assert.NotEqual(t, http.StatusOK, w.Code) 
	})
}