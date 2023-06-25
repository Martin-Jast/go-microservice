package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/Martin-Jast/go-microservice/utils"
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


// MongoBaseModel small extension of the generic baseModel to accomodate mongoID
type MongoBaseModel struct {
	ID *primitive.ObjectID `bson:"_id"`
	*BaseModel `bson:"inline"`
}


func (m MongoAdapter) Create(ctx context.Context, document BaseModel) (id string, err error) {
	mBase := MongoBaseModel{}

	if document.ID == nil {
		temp := primitive.NewObjectID()
		mBase.ID = &temp
	}
	mBase.BaseModel = &document
	if mBase.CreatedAt == nil {
		temp := time.Now()
		mBase.CreatedAt = &temp
	}
	res, err := m.mongoConnection.InsertOne(ctx, mBase)
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
		return nil, fmt.Errorf("cannot GetByID with no id")
	}
	asObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid objectID to find")
	}
	result := m.mongoConnection.FindOne(ctx, bson.M{"_id": asObjID})
	if result.Err() != nil {
		return nil, result.Err()

	}
	res := MongoBaseModel{}
	err = result.Decode(&res)
	if err != nil {
		return nil, err
	}
	res.BaseModel.ID = utils.StrPnt(res.ID.Hex())
	return res.BaseModel, nil
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
	result, err := m.mongoConnection.Find(ctx, bson.M{"created_at": bson.M{"$gt": primitive.NewDateTimeFromTime(date)}})
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)

	list := []BaseModel{}
	for result.Next(ctx) {
		elem := MongoBaseModel{}
		err = result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		elem.BaseModel.ID = utils.StrPnt(elem.ID.Hex())
		list = append(list, *elem.BaseModel)
	}

	return list, nil
}

func (m MongoAdapter) DeleteAll(ctx context.Context) (error) {
	_, err := m.mongoConnection.DeleteMany(ctx, bson.M{})
	return err
}
