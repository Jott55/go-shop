package routes

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/argon2"
)

type LoginUserRequest struct {
	LoginUser *types.LoginUser
}

func parseArgon2ID(s string) *Argon2idConfig {
	params := strings.Split(s, "$")

	if len(params) != 6 {
		return nil
	}

	if !strings.HasPrefix(params[1], "argon2id") {
		return nil
	}

	config := &Argon2idConfig{}

	fmt.Sscanf(params[3], "m=%d,t=%d,p=%d", &config.memory, &config.time, &config.threads)

	salt, err := base64.RawStdEncoding.DecodeString(params[4])
	checkError(err)
	config.salt = salt

	hash, err := base64.RawStdEncoding.DecodeString(params[5])
	checkError(err)
	config.hash = hash

	config.keyLen = uint32(len(hash))

	return config
}

func Login(router *chi.Mux) {
	router.Post("/login/user", func(w http.ResponseWriter, r *http.Request) {
		ruser, err := serverio.GetStructFromRequestBody[LoginUserRequest](r)

		luser := ruser.LoginUser

		if checkError(err) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// get salted password
		uLogin, err := ser.User.GetUserByEmail(luser.Email)
		checkError(err)
		// parse it
		pass_config := parseArgon2ID(uLogin.Password)

		// verify authenticity
		new_hash := argon2.IDKey(
			[]byte(luser.Password),
			pass_config.salt,
			pass_config.time,
			pass_config.memory,
			pass_config.threads,
			pass_config.keyLen,
		)

		okay := subtle.ConstantTimeCompare(pass_config.hash, new_hash) == 1

		if !okay {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token := createTokenString(uLogin.Name)
		w.Write([]byte(token))
	})
}
