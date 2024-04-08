package controllers

import (
	"context"
	"net/http"
	"time"

	repo "apna-restaurant-2.0/db/sqlc"
	"apna-restaurant-2.0/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *repo.Queries
}

func NewAuthController(db *repo.Queries) *AuthController {
	return &AuthController{db}
}

func (ac *AuthController) SignUpUser(c *gin.Context) {
	var credentials *repo.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if resp, ok := utils.ValidateUserRegisterOrLogin(credentials, ""); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	existingUserCount, err := ac.db.CheckExistingUser(context.Background(), credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if existingUserCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	user := &repo.CreateUserParams{
		Email:       credentials.Email,
		Name:        credentials.Name,
		PhoneNumber: credentials.PhoneNumber,
		Password:    string(hashedPassword),
		UpdatedAt:   time.Now(),
	}
	createdUser, err := ac.db.CreateUser(context.Background(), *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Adding user in DB"})
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "data": createdUser})
}

func (ac *AuthController) SignInUser(c *gin.Context) {
	var credentials *repo.User
	response := make(map[string]interface{})

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if resp, ok := utils.ValidateUserRegisterOrLogin(credentials, "login"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp})
		return
	}
	existingUserCount, err := ac.db.CheckExistingUser(context.Background(), credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if existingUserCount != 1 {
		c.JSON(http.StatusConflict, gin.H{"error": "User does not exist"})
		return
	}
	dbUser, err := ac.db.GetUserByEmail(context.Background(), credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	token := utils.GenerateToken(dbUser.ID, dbUser.Email)
	data := map[string]interface{}{
		"user":  dbUser,
		"token": token,
	}
	response["data"] = data
	c.JSON(http.StatusOK, response)
}
