package mongo_config

import (
	"context"
	"errors"
	"time"

	"github.com/rupesh-sengar/golang-collection/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	coll *mongo.Collection
}

func EnsureUserIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "application", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "application", Value: 1}},
			Options: options.Index().SetName("application_idx"),
		},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}

func NewUserRepository(coll *mongo.Collection) domain.UserRepository {
	return &userRepo{coll: coll}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.coll.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("email already in use for this application")
	}
	return err
}

func (r *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "audit.deleted": false}).Decode(&u)
	return &u, err
}

func (r *userRepo) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	var u domain.User
	filter := bson.M{"email": email, "audit.deleted": false}
	err := r.coll.FindOne(ctx, filter).Decode(&u)
	return &u, err
}

func (r *userRepo) FindByApplication(ctx context.Context, applicationID string) ([]*domain.User, error) {
	filter := bson.M{"application": applicationID, "audit.deleted": false}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var u domain.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, cursor.Err()
}

func (r *userRepo) Approve(ctx context.Context, id primitive.ObjectID, approverID string) error {
	update := bson.M{
		"$set": bson.M{
			"approvedBy":      approverID,
			"status":          domain.StatusApproved,
			"audit.updatedAt": time.Now().UTC(),
			"audit.updatedBy": approverID,
		},
		"$inc": bson.M{"audit.version": 1},
	}
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "audit.deleted": false, "status": domain.StatusPending},
		update,
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found or not pending approval")
	}
	return nil
}

func (r *userRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus domain.UserStatus, actorID string) error {
	update := bson.M{
		"$set": bson.M{
			"status":          newStatus,
			"audit.updatedAt": time.Now().UTC(),
			"audit.updatedBy": actorID,
		},
		"$inc": bson.M{"audit.version": 1},
	}
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "audit.deleted": false},
		update,
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id primitive.ObjectID, actorID primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"audit.deleted":   true,
			"audit.updatedAt": time.Now().UTC(),
			"audit.updatedBy": actorID,
		},
		"$inc": bson.M{"audit.version": 1},
	}
	_, err := r.coll.UpdateByID(ctx, id, update)
	return err
}

func (r *userRepo) Update(ctx context.Context, u *domain.User, actorID string) error {
	u.Audit.UpdatedAt = time.Now().UTC()
	u.Audit.UpdatedBy = actorID
	u.Audit.Version++

	if err := u.Validate(); err != nil {
		return err
	}

	res, err := r.coll.ReplaceOne(ctx,
		bson.M{"_id": u.ID, "audit.deleted": false, "audit.version": u.Audit.Version - 1},
		u,
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found or version mismatch")
	}
	return nil
}
