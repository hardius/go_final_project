package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "my_secret_key"

type Claims struct {
	Hash string
	jwt.RegisteredClaims
}

type PasswordJson struct {
	Password string `json:"password"`
}

type TokenJson struct {
	Token string `json:"token"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var p PasswordJson

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusInternalServerError)
		return
	}

	if p.Password != os.Getenv("TODO_PASSWORD") {
		err = errors.New("Wrong password.")
		fmt.Println("2")
		writeJson(w, errWrap(err), http.StatusInternalServerError)
		return
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	hash, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		writeJson(w, errWrap(err), http.StatusInternalServerError)
		return
	}

	claims := Claims{Hash: hash}

	jwtTokenResp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtTokenResp.SignedString([]byte(p.Password))
	if err != nil {
		writeJson(w, errWrap(err), http.StatusInternalServerError)
		return
	}

	tokenResponse := TokenJson{
		Token: signedToken,
	}
	writeJson(w, tokenResponse, http.StatusOK)
}

func ValidateToken(tokenString string) bool {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	hash, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return false
	}

	claims := Claims{Hash: hash}

	jwtTokenResp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtTokenResp.SignedString([]byte(os.Getenv("TODO_PASSWORD")))
	if err != nil {
		return false
	}

	return tokenString == signedToken
}
