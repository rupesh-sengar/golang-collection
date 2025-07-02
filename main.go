package main

import (
	"fmt"
	"git"
	"github.com/rupesh-sengar/golang-collection/auth"
)

func main(){
	fmt.Println("Welcome to Golang Collection")
	fmt.Println("This is the Auth Service")
	auth.StartServer()
	fmt.Println("Auth Service started successfully")
}