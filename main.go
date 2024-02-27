package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const baseUriPath = "/api/v1"

func main() {
	InitDB()

	r := gin.Default()
	r.POST(fmt.Sprintf("%s/auth/register", baseUriPath), RegisterUser)
	r.POST(fmt.Sprintf("%s/auth/login", baseUriPath), LoginUser)
	r.GET(fmt.Sprintf("%s/auth/currentUser", baseUriPath), VerifyToken, GetCurrentUserDetails)

	r.Run()
}
