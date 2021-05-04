package usercontrollers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

//hashAndSalt: encrypt password
func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

//comparePasswords: compare hashed password with a clear one
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//getJWT generate an access token, valid during 60 min
func getJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Krissanawat"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(MySigningKey)

	if err != nil {
		return "", fmt.Errorf("Something Went Wrong: %s", err.Error())
	}

	return tokenString, nil
}

//IsAuthorized this func allows us to check if a valid token has been sent with a request
func IsAuthorized(r *http.Request) (string, error) {
	if r.Header["Authorization"] != nil {
		token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("Invalid Signing Method"))
			}
			aud := "billing.jwtgo.io"
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAudience {
				return nil, fmt.Errorf(("invalid aud"))
			}
			// verify iss claim
			iss := "jwtgo.io"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return nil, fmt.Errorf(("invalid iss"))
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			return "", err
		}

		if token.Valid {
			return r.Header["Authorization"][0], nil
		}

	}
	return "", fmt.Errorf(("no Token detected"))
}
