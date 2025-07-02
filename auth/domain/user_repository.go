package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error

	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)

	FindByEmail(ctx context.Context, email Email) (*User, error)

	FindByApplication(ctx context.Context, applicationID string) ([]*User, error)

	Approve(ctx context.Context, id primitive.ObjectID, approverID string) error

	UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus UserStatus, actorID string) error

	Delete(ctx context.Context, id primitive.ObjectID, actorID primitive.ObjectID) error

	Update(ctx context.Context, u *User, actorID string) error
}
