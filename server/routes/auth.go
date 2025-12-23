package routes

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
		return "", fmt.Errorf("wrong header format: %s", header)
	}
	return tokens[1], nil
}

func Auth(router chi.Router) {
	router.Get("/auth/create-pem-files", func(w http.ResponseWriter, r *http.Request) {
		pubKey, privKey, err := ed25519.GenerateKey(nil)
		checkError(err)

		bPriv, err := x509.MarshalPKCS8PrivateKey(privKey)
		checkError(err)
		privPem := &pem.Block{
			Type:    "PRIVATE KEY",
			Headers: nil,
			Bytes:   bPriv,
		}

		bPub, err := x509.MarshalPKIXPublicKey(pubKey)

		pubPem := &pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   bPub,
		}

		privPemFile, err := os.Create("private-auth-ed25519.pem")
		checkError(err)
		pubPemFile, err := os.Create("public-auth-ed25519.pem")
		checkError(err)

		err = pem.Encode(privPemFile, privPem)
		checkError(err)
		err = pem.Encode(pubPemFile, pubPem)
		checkError(err)

		w.Write([]byte("Auth files created!"))
	})
}
