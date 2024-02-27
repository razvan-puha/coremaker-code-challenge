package main

import (
	"net/mail"

	"github.com/gin-gonic/gin"
)

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDetails struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user
// @Accept  json
// @Produce  json
// @Param userRegistration body UserRegistration true "User registration details"
// @Success 201 {string} string "User registered"
// @Failure 400 {string} string "Invalid email address"
// @Failure 500 {string} string "Unable to parse request body"
// @Router /auth/register [post]
func RegisterUser(c *gin.Context) {
	var userReg UserRegistration

	err := c.Bind(&userReg)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Unable to parse request body",
		})
		return
	}
	if !validateEmail(userReg.Email) {
		c.JSON(400, gin.H{
			"message": "Invalid email address",
		})
	}

	AddUser(userReg.Email, userReg.Password, userReg.Name)
	c.JSON(201, gin.H{
		"message": "User registered",
	})
}

// LoginUser godoc
// @Summary Login a user
// @Description Login a user
// @Accept  json
// @Produce  json
// @Param userLogin body UserLogin true "User login details"
// @Success 200 {string} string "User logged in"
// @Failure 400 {string} string
// @Router /auth/login [post]
func LoginUser(c *gin.Context) {
	var loginDetails UserLogin

	err := c.Bind(&loginDetails)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Unable to parse request body",
		})
		return
	}
	if !validateEmail(loginDetails.Email) {
		c.JSON(400, gin.H{
			"message": "Invalid email address",
		})
	}

	token, err := Login(loginDetails.Email, loginDetails.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"token": token,
		})
	}
}

// GetCurrentUserDetails godoc
// @Summary Get current user details
// @Description Get current user details
// @Accept  json
// @Produce  json
// @Success 200 {object} UserDetails
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/currentUser [get]
func GetCurrentUserDetails(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth != "" {
		token := auth[7:]
		user, err := GetLoggedUserByToken(token)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, &UserDetails{
			Email: user.Email,
			Name:  user.Name,
		})
	}
}

func VerifyToken(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	token := auth[7:]

	loggedUser, err := GetLoggedUserByToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	if loggedUser == nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	} else {
		c.Next()
		return
	}
}

func validateEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}
