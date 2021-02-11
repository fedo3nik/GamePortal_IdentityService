package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDTO struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	WarningCount uint               `bson:"warningCount,omitempty"`
	Nickname     string             `bson:"nickname,omitempty"`
	Password     string             `bson:"password,omitempty"`
	Email        string             `bson:"email,omitempty"`
	TokenHash    string             `bson:"tokenHash,omitempty"`
}
