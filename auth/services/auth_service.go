package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rupesh-sengar/auth/domain"
	"github.com/rupesh-sengar/auth/infra/mongo_config"
	"github.com/rupesh-sengar/auth/utils/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Auth0TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token,omitempty"`
}

func Auth0Login(username, password string) (*Auth0TokenResponse, error) {
	payload := fmt.Sprintf(`{
		"grant_type": "password",
		"username": "%s",
		"password": "%s",
		"audience": "%s",
		"client_id": "%s",
		"client_secret": "%s",
		"scope": "openid profile email"
	}`,
		username,
		password,
		os.Getenv("AUTH0_AUDIENCE"),
		os.Getenv("AUTH0_CLIENT_ID"),
		os.Getenv("AUTH0_CLIENT_SECRET"),
	)

	resp, err := http.Post(
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/oauth/token",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Auth0 error: %s", body)
	}

	var token Auth0TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func Auth0Signup(req types.SignupRequest, password string) (*http.Response, error) {
	auth0Payload :=
		fmt.Sprintf(`{
	  		"client_id": "%s",
	 	   	"email": "%s",
		   	"password": "%s",
			"connection": "Username-Password-Authentication",
			"user_metadata":{
				"first_name": "%s",
				"last_name": "%s",
				"role":"member",
				"application": "%s"
				}
	  }`,
			os.Getenv("AUTH0_CLIENT_ID"),
			req.Email,
			password,
			req.FirstName,
			req.LastName,
			req.Application,
		)

	fmt.Println("Auth0 signup payload: ", auth0Payload)
	resp, err := http.Post(
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/dbconnections/signup",
		"application/json",
		strings.NewReader(auth0Payload))

	fmt.Println("Auth0 signup response: ", resp)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth0 signup failed with status %d: %s", resp.StatusCode, body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println("mongo connect error: %v", err)
	}
	defer client.Disconnect(ctx)
	userCollection := client.Database("User-Management").Collection("Users")

	userRepo := mongo_config.NewUserRepository(userCollection)
	mongo_config.EnsureUserIndexes(ctx, userCollection)
	domainRes, err := domain.NewUser(req.FirstName, req.LastName, req.Email, req.Password, req.CreatorID, req.Application)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return nil, err
	}

	createErr := userRepo.Create(ctx, domainRes)

	if createErr != nil {
		fmt.Println("Error creating user in MongoDB:", createErr)
		return nil, createErr
	}

	return nil, nil
}

func AuthApprovalService(req types.AuthApprovalRequest) (*http.Response, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println("mongo connect error: %v", err)
	}
	defer client.Disconnect(ctx)
	userCollection := client.Database("User-Management").Collection("Users")

	userRepo := mongo_config.NewUserRepository(userCollection)
	mongo_config.EnsureUserIndexes(ctx, userCollection)

	err = userRepo.Approve(ctx, req.Id, req.ApproverID)
	if err != nil {
		fmt.Println("Error approving user:", err)
		return nil, err
	}
	return nil, nil
}
