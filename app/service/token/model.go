package token

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func CreateJWT() (string, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Expiry"] = time.Now().Add(15 * time.Minute)
	tokenSecret := []byte(os.Getenv("tokenSecret"))

	tokenStr, err := token.SignedString(tokenSecret)
	if err != nil {
		fmt.Println("Error signing token")
		fmt.Println(err)
		return "", err
	}

	fmt.Println(tokenStr)
	return tokenStr, nil
}
