package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"file-sharing/utilities"
)

const userCollectionName string = "user"

type userService struct {
	client     *mongo.Client
	collection *mongo.Collection
	database   *mongo.Database
}

type UserModel struct {
	CreatedAt      int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt      int64  `json:"deletedAt" bson:"deletedAt"`
	Email          string `json:"email" bson:"email"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
	PasswordHash   string `json:"-" bson:"passwordHash"`
	Role           string `json:"role" bson:"role"`
	SetUpCompleted bool   `json:"setUpCompleted" bson:"setUpCompleted"`
	Uid            string `json:"uid" bson:"uid"`
	UpdatedAt      int64  `json:"updatedAt" bson:"updatedAt"`
}

var UserService userService

func (service *userService) New(client *mongo.Client, database *mongo.Database) {
	service.client = client
	service.database = database
	service.collection = service.database.Collection(userCollectionName)
}

func (service *userService) DeleteOne(
	operationContext context.Context,
	filter map[string]any,
) error {
	_, queryError := service.collection.DeleteOne(operationContext, filter)
	return queryError
}

func (service *userService) FindOne(
	operationContext context.Context,
	filter map[string]any,
	destination *UserModel,
) error {
	return service.collection.FindOne(operationContext, filter).Decode(destination)
}

func (service *userService) FindOneAndUpdate(
	operationContext context.Context,
	filter map[string]any,
	update map[string]any,
	destination *UserModel,
) error {
	return service.collection.FindOneAndUpdate(
		operationContext,
		filter,
		update,
	).Decode(destination)
}

func (service *userService) FindOneByUid(
	operationContext context.Context,
	uid string,
	destination *UserModel,
) error {
	return service.collection.FindOne(operationContext, bson.M{"uid": uid}).Decode(destination)
}

func (service *userService) FindPaginated(
	operationContext context.Context,
	filter map[string]any,
	pagination utilities.PaginationData,
	destination *[]UserModel,
) (int64, error) {
	count, queryError := service.collection.CountDocuments(operationContext, filter)
	if queryError != nil {
		return 0, queryError
	}

	queryOptions := options.Find()
	queryOptions.SetLimit(int64(pagination.Limit))
	queryOptions.SetSkip(int64(pagination.Offset))
	cursor, cursorError := service.collection.Find(operationContext, filter, queryOptions)
	if cursorError != nil {
		return 0, cursorError
	}

	defer cursor.Close(context.Background())
	return count, cursor.All(context.Background(), destination)
}

func (service *userService) InsertOne(
	operationContext context.Context,
	document UserModel,
) error {
	_, queryError := service.collection.InsertOne(operationContext, document)
	return queryError
}

func (service *userService) UpdateOne(
	operationContext context.Context,
	filter map[string]any,
	update map[string]any,
) error {
	_, queryError := service.collection.UpdateOne(
		operationContext,
		filter,
		update,
	)
	return queryError
}
