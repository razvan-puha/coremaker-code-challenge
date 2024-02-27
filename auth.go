package main

import (
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

func RegisterUser(c *gin.Context) {
	var userReg UserRegistration

	err := c.Bind(&userReg)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Unable to parse request body",
		})
		return
	}

	AddUser(userReg.Email, userReg.Password, userReg.Name)
	c.JSON(201, gin.H{
		"message": "User registered",
	})
}

func LoginUser(c *gin.Context) {
	var loginDetails UserLogin

	err := c.Bind(&loginDetails)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Unable to parse request body",
		})
		return
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
