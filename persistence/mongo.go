package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAdapter struct {
	mongoConnection *mongo.Collection
}

func NewMongoAdapter(dbClient *mongo.Client) MongoAdapter {
	return MongoAdapter{
		mongoConnection: dbClient.Database("test").Collection("base"),
	}
}

func (m MongoAdapter) Create(ctx context.Context, document BaseModel) (id string, err error) {
	if document.ID == nil {
		temp := primitive.NewObjectID()
		document.ID = &temp
	}
	document.CreatedAt = time.Now()
	res, err := m.mongoConnection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	respID, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return respID.Hex(), nil
	}

	return "", fmt.Errorf("could not insert document")
}

func (m MongoAdapter) GetByID(ctx context.Context, id string) (doc *BaseModel, err error) {
	if id == "" {
		return nil, fmt.Errorf("cannot delete with no id")
	}
	asObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid objectID to delete")
	}
	result := m.mongoConnection.FindOne(ctx, bson.M{"_id": asObjID})
	if result.Err() != nil {
		return nil, result.Err()

	}
	res := BaseModel{}
	err = result.Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m MongoAdapter) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("cannot delete with no id")
	}
	asObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid objectID to delete")
	}
	_, err = m.mongoConnection.DeleteOne(ctx, bson.M{"_id": asObjID})
	return err
}

func (m MongoAdapter) GetAllCreatedSince(ctx context.Context, date time.Time) (docs []BaseModel, err error) {
	result, err := m.mongoConnection.Find(ctx, bson.M{"created_at": bson.M{"$gt": date}})
	if err != nil {
		return nil, err
	}
	list := []BaseModel{}
	for result.Next(ctx) {
		elem := BaseModel{}
		err = result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
	}

	return list, nil
}