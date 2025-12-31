package helpers

import (
	// "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 1. Struktur Standar Response
type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

// 2. Fungsi Pembungkus JSON (Wrapper)
func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

// 3. Fungsi Translate Error Validasi
func FormatValidationError(err error) []string {
	var errors []string

	// Loop setiap error dan ubah jadi pesan simpel
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error()) 
        // Note: Nanti bisa kita custom pesannya di sini jadi bahasa Indonesia
	}

	return errors
}