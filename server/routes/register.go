package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/argon2"
)

type Argon2idConfig struct {
	password []byte
	salt     []byte
	time     uint32
	memory   uint32
	threads  uint8
	keyLen   uint32
	hash     []byte
}

func argonConfigure(pass string) *Argon2idConfig {
	var a Argon2idConfig = Argon2idConfig{
		password: []byte(pass),
		salt:     make([]byte, 32),
		time:     1,
		memory:   32,
		threads:  8,
		keyLen:   32,
		hash:     nil,
	}
	_, err := rand.Read(a.salt)
	checkError(err)
	return &a
}

func createPassword(pass string) string {
	ac := argonConfigure(pass)

	ac.hash = argon2.IDKey(ac.password, ac.salt, ac.time, ac.memory, ac.threads, ac.keyLen)

	b64_hash := base64.RawStdEncoding.EncodeToString(ac.hash)
	b64_salt := base64.RawStdEncoding.EncodeToString(ac.salt)

	argon2format := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, ac.memory, ac.time, ac.threads, b64_salt, b64_hash)
	return argon2format
}

func Register(router *chi.Mux) {
	router.Post("/register/user", func(w http.ResponseWriter, r *http.Request) {
		request, err := serverio.GetStructFromRequestBody[types.UserRequest](r)

		// check if request is valid
		if checkError(err) ||
			len(request.User.Email) < 10 ||
			len(request.User.Name) < 4 ||
			len(request.User.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create user password
		request.User.Password = createPassword(request.User.Password)

		// Create user of request
		ser.User.Insert(request.User)
		// Get user id
		user_id, err := ser.User.GetIdByName(request.User.Name)
		checkError(err)

		// Create user cart
		ser.Cart.Insert(&types.CartNoId{User_id: user_id})

		// Get signed token adding username to it
		token := createTokenString(request.User.Name)

		// send token
		w.Write([]byte(token))
	})

	// test your key here
	router.Get("/test/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		claism := decryptTokenString(key)
		debug(claism)
	})
}
