package utils

import (
	"os"
	"strings"
	"time"

	repo "apna-restaurant-2.0/db/sqlc"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Token struct {
	UserId uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}

func ValidateUserRegisterOrLogin(user *repo.User, flag string) (string, bool) {
	if len(strings.TrimSpace(user.Name)) == 0 && flag != "login" {
		return "Name required", false
	} else if len(strings.TrimSpace(user.Name)) < 4 && flag != "login" {
		return "Name should be at least 3 chars", false
	} else if len(strings.TrimSpace(user.Email)) == 0 {
		return "Email required", false
	} else if !IsEmailValid(user.Email) {
		return "Invalid Email", false
	} else if len(strings.TrimSpace(user.Password)) == 0 {
		return "Password required", false
	} else if len(strings.TrimSpace(user.Password)) < 6 {
		return "Password should be at least 6 chars", false
	} else if len(strings.TrimSpace(user.PhoneNumber)) == 0 && flag != "login" {
		return "Phone num required", false
	} else if !IsPhoneValid(user.PhoneNumber) && flag != "login" {
		return "Invalid Phonenumber", false
	}
	return "Requirement passed", true
}
func GenerateToken(userId uuid.UUID, email string) string {
	tokenClaims := &Token{
		UserId: userId,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims)
	signedToken, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	return signedToken
}
