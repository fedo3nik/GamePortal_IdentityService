package mongodb

import (
	"context"
	"log"

	"github.com/fedo3nik/GamePortal_IdentityService/internal/domain/entities"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/infrastructure/dtodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(client *mongo.Client, db string) *mongo.Collection {
	database := client.Database(db)
	usersCollection := database.Collection("users")

	return usersCollection
}

func Insert(ctx context.Context, collection *mongo.Collection, user *entities.User) (*mongo.InsertOneResult, error) {
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return nil, err
	}

	return insertResult, nil
}

func DeleteAll(ctx context.Context, collection *mongo.Collection) (*mongo.DeleteResult, error) {
	deleteResult, err := collection.DeleteOne(ctx, bson.D{})
	if err != nil {
		log.Printf("Delete document error: %v", err)
		return nil, err
	}

	return deleteResult, nil
}

func UpdateWarningCountField(ctx context.Context, collection *mongo.Collection, id string, warnCount uint) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Convert to OjectId error: %v", err)
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oid}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "warningCount", Value: warnCount}}}}

	updateResult, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

func GetDocumentByID(ctx context.Context, collection *mongo.Collection, id string) (*dtodb.UserDTO, error) {
	usr := dtodb.UserDTO{Nickname: "JhonDoe"}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Convert to OjectId error: %v", err)
	}

	err = collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: oid}}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func GetDocumentByEmail(ctx context.Context, collection *mongo.Collection, email string) (*dtodb.UserDTO, error) {
	usr := dtodb.UserDTO{Nickname: "JhonDoe"}

	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "email", Value: email}}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
