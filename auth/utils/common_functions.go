package utils

import (
	"context"
	"fmt"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/rupesh-sengar/golang-collection/auth/domain"
	"github.com/rupesh-sengar/golang-collection/auth/infra/mongo_config"
)

func CheckUserStatus(email string) (domain.UserStatus, error) {
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
	user, err := userRepo.FindByEmail(ctx, domain.Email(email))
	if err != nil {
		fmt.Println("Error finding user by email:", err)
		return "", err
	}
	return user.Status, nil
}