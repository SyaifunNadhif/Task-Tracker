package middlewares

import (
	"net/http"
	"strings"
	"task-tracker/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Header Authorization
		authHeader := c.GetHeader("Authorization")

		// Cek format: "Bearer <token>"
		if !strings.Contains(authHeader, "Bearer") {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// 2. Ambil string tokennya saja (buang kata "Bearer "-nya)
		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		// 3. Validasi Token
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return helpers.SECRET_KEY, nil // Menggunakan kunci rahasia yang sama dari helpers
		})

		// 4. Jika token valid, lanjut ke Controller
		if token != nil && token.Valid {
			// AMBIL DATA DARI TOKEN (CLAIMS)
			claims := token.Claims.(jwt.MapClaims)
			
			// Ambil user_id sebagai float64 dulu (bawaan JWT)
			userIDFloat := claims["user_id"].(float64) 

			// Ubah jadi uint sebelum disimpan ke Context
			userIDUint := uint(userIDFloat)

			// Simpan yang versi UINT ke context
			c.Set("currentUser", userIDUint)

			c.Next() // Silakan lewat
		} else {
			// Jika token expired atau palsu
			response := helpers.APIResponse("Unauthorized: Token Invalid", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
