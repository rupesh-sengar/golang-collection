package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rupesh-sengar/golang-collection/auth/services"
	"github.com/rupesh-sengar/golang-collection/auth/utils"
	"github.com/rupesh-sengar/golang-collection/auth/utils/types"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	EncPassword string `json:"enc_password" binding:"required"`
	Application string `json:"application" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	status, err:=utils.CheckUserStatus(req.Username)
	fmt.Println("User status: ", status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user status", "detail": err.Error()})
		return
	}
	if status != "approved"{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not approved"})
		return
	}
	password, err := utils.DecryptEncPassword(req.EncPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decryption failed", "detail": err.Error()})
		return
	}

	token, err := services.Auth0Login(req.Username, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth0 login failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

func SignupHandler(c *gin.Context) {
	var req types.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	password, err := utils.DecryptEncPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decryption failed", "detail": err.Error()})
		return
	}
	resp, err := services.Auth0Signup(req, password)
	fmt.Println("Auth0 signup erro: ", err)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") || strings.Contains(err.Error(), "email already in use") {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed", "detail": err.Error()})
		return
	}

	fmt.Println("Auth0 signup response status:", resp)

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})

}

func AuthApprovalHandler(c *gin.Context) {
	var req types.AuthApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if _, err := services.AuthApprovalService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Approval failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}
