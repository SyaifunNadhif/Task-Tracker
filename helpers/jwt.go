package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ambil SECRET_KEY dari .env, kalau gak ada pakai default (buat dev aja)
var SECRET_KEY = []byte(getSecretKey())

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "rahasia_negara_api_ini_sangat_aman" // Default key (JANGAN DIPAKAI DI PRODUCTION)
	}
	return secret
}

func GenerateToken(userID uint) (string, error) {
	// 1. Tentukan isi token (Claims)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	// 2. Buat token dengan algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Tanda tangani token dengan Secret Key
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}