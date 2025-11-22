package routes

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Authorize() {
}

func decryptTokenString(token_string string) map[string]any {
	b, _ := os.ReadFile("public-auth-ed25519.pem")
	pkey, err := jwt.ParseEdPublicKeyFromPEM(b)

	if checkError(err) {
		return nil
	}

	token, err := jwt.Parse(token_string, func(t *jwt.Token) (any, error) {
		return pkey, nil
	})

	if checkError(err) {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		debug(claims["sub"], claims["nbf"])
		return claims
	}
	return nil
}

func createTokenString(username any) string {
	data, _ := os.ReadFile("private-auth-ed25519.pem")
	key, err := jwt.ParseEdPrivateKeyFromPEM(data)

	checkError(err)

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"sub": username,
		"nbf": time.Now().Unix(),
	})

	token_string, err := token.SignedString(key)
	checkError(err)
	return token_string
}
