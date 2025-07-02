package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"time"
)

type Email string

func (e Email) Validate() error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(string(e)) {
		return errors.New("invalid email format")
	}
	return nil
}

type Name struct {
	First string `validate:"required,alpha,min=2,max=50" bson:"first" json:"first"`
	Last  string `validate:"required,alpha,min=2,max=50" bson:"last" json:"last"`
}

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
	RoleGuest  UserRole = "guest"
)

func RoleValidator(fl validator.FieldLevel) bool {
	role := UserRole(fl.Field().String())
	switch role {
	case RoleAdmin, RoleMember, RoleGuest:
		return true
	default:
		return false
	}
}

type UserStatus string

const (
	StatusPending   UserStatus = "pending"
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusApproved  UserStatus = "approved"
)

// StatusValidator ensures the status field is one of the predefined states.
func StatusValidator(fl validator.FieldLevel) bool {
	status := UserStatus(fl.Field().String())
	switch status {
	case StatusPending, StatusActive, StatusSuspended, StatusApproved:
		return true
	default:
		return false
	}
}

// —— Audit Metadata ——

type Audit struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy" json:"updatedBy"`
	Version   int64     `bson:"version" json:"version"`
	Deleted   bool      `bson:"deleted" json:"deleted"`
}

// —— User Entity ——

type User struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Name        Name                   `bson:"name" json:"name" validate:"required"`
	Email       Email                  `bson:"email" json:"email" validate:"required,emailVO"`
	Password    string                 `bson:"passwordHash" json:"-" validate:"required,min=8"`
	Roles       []UserRole             `bson:"roles" json:"roles" validate:"required,dive,role"`
	Application string                 `bson:"application" json:"application" validate:"required"`
	ApprovedBy  string                 `bson:"approvedBy,omitempty" json:"approvedBy,omitempty"`
	Status      UserStatus             `bson:"status" json:"status" validate:"required,status"`
	Meta        map[string]interface{} `bson:"meta,omitempty" json:"meta,omitempty"`
	Audit       Audit                  `bson:"audit" json:"audit"`
}

func NewUser(first, last, rawEmail, plainPassword string, creatorID string, applicationID string) (*User, error) {
	u := &User{
		ID:          primitive.NewObjectID(),
		Name:        Name{First: first, Last: last},
		Email:       Email(rawEmail),
		Password:    HashPassword(plainPassword),
		Roles:       []UserRole{RoleMember},
		Application: applicationID,
		ApprovedBy:  "",
		Status:      StatusPending,
		Meta:        make(map[string]interface{}),
		Audit: Audit{
			CreatedAt: time.Now().UTC(),
			CreatedBy: creatorID,
			UpdatedAt: time.Now().UTC(),
			UpdatedBy: creatorID,
			Version:   1,
			Deleted:   false,
		},
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Validate() error {
	v := validator.New()

	v.RegisterValidation("emailVO", func(fl validator.FieldLevel) bool {
		return Email(fl.Field().String()).Validate() == nil
	})

	v.RegisterValidation("role", RoleValidator)

	v.RegisterValidation("status", StatusValidator)
	return v.Struct(u)
}

func HashPassword(pw string) string {
	return pw
}
