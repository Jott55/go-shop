package routes

import (
	"crypto"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var public_key crypto.PublicKey
var private_key crypto.PrivateKey

func decryptTokenString(token_string string) map[string]any {
	if public_key == nil {
		encoded_pub_key, _ := os.ReadFile("public-auth-ed25519.pem")
		var err error
		public_key, err = jwt.ParseEdPublicKeyFromPEM(encoded_pub_key)

		checkError(err)
	}

	token, err := jwt.Parse(token_string, func(t *jwt.Token) (any, error) {
		return public_key, nil
	})

	if checkError(err) {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims
	}
	return nil
}

func createTokenString(username string) string {
	if private_key == nil {
		encoded_private_key, _ := os.ReadFile("private-auth-ed25519.pem")
		var err error
		private_key, err = jwt.ParseEdPrivateKeyFromPEM(encoded_private_key)

		checkError(err)
	}

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"sub": username,
		"nbf": time.Now().Unix(),
	})

	token_string, err := token.SignedString(private_key)
	checkError(err)
	return token_string
}

func getTokenFromHeader(header string) (string, error) {
	tokens := strings.Split(header, " ")
	if len(tokens) < 2 {
		return "", fmt.Errorf("wrong header format")
	}
	return tokens[1], nil
}
