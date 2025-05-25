package database

import (
	"context"
	"file-sharing/utilities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const fileCollectionName string = "file"

type fileService struct {
	client     *mongo.Client
	collection *mongo.Collection
	database   *mongo.Database
}

type FileModel struct {
	CreatedAt    int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt    int64  `json:"deletedAt" bson:"deletedAt"`
	IsDeleted    bool   `json:"isDeleted" bson:"isDeleted"`
	OriginalName string `json:"originalName" bson:"originalName"`
	Size         int64  `json:"size" bson:"size"`
	Uid          string `json:"uid" bson:"uid"`
	UpdatedAt    int64  `json:"updatedAt" bson:"updatedAt"`
}

var FileService fileService

func (service *fileService) New(client *mongo.Client, database *mongo.Database) {
	service.client = client
	service.database = database
	service.collection = service.database.Collection(fileCollectionName)
}

func (service *fileService) DeleteOne(
	operationContext context.Context,
	filter map[string]any,
) error {
	_, queryError := service.collection.DeleteOne(operationContext, filter)
	return queryError
}

func (service *fileService) DeleteOneByUid(
	operationContext context.Context,
	uid string,
) error {
	_, queryError := service.collection.DeleteOne(operationContext, bson.M{"uid": uid})
	return queryError
}

func (service *fileService) FindOne(
	operationContext context.Context,
	filter map[string]any,
	destination *FileModel,
) error {
	return service.collection.FindOne(operationContext, filter).Decode(destination)
}

func (service *fileService) FindOneAndUpdate(
	operationContext context.Context,
	filter map[string]any,
	update map[string]any,
	destination *FileModel,
) error {
	return service.collection.FindOneAndUpdate(
		operationContext,
		filter,
		update,
	).Decode(destination)
}

func (service *fileService) FindOneByUid(
	operationContext context.Context,
	uid string,
	destination *FileModel,
) error {
	return service.collection.FindOne(operationContext, bson.M{"uid": uid}).Decode(destination)
}

func (service *fileService) FindPaginated(
	operationContext context.Context,
	filter map[string]any,
	pagination utilities.PaginationData,
	destination *[]FileModel,
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

func (service *fileService) InsertOne(
	operationContext context.Context,
	document FileModel,
) error {
	_, queryError := service.collection.InsertOne(operationContext, document)
	return queryError
}

func (service *fileService) UpdateMany(
	operationContext context.Context,
	filter map[string]any,
	update map[string]any,
) error {
	_, queryError := service.collection.UpdateMany(
		operationContext,
		filter,
		update,
	)
	return queryError
}

func (service *fileService) UpdateOne(
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
