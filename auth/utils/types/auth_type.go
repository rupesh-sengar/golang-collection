package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignupRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Application string `json:"application"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CreatorID   string `json:"creator_id"`
}

type AuthApprovalRequest struct{
	Id primitive.ObjectID `json:"id"`
	ApproverID string `json:"approver_id"`
}