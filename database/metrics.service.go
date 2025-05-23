package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const metricsCollectionName string = "metrics"

type metricsService struct {
	client     *mongo.Client
	collection *mongo.Collection
	database   *mongo.Database
}

type MetricsModel struct {
	CreatedAt      int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt      int64  `json:"deletedAt" bson:"deletedAt"`
	Downloads      int64  `json:"downloads" bson:"downloads"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
	LastDownloaded int64  `json:"lastDownloaded" bson:"lastDownloaded"`
	LastViewed     int64  `json:"lastViewed" bson:"lastViewed"`
	Uid            string `json:"uid" bson:"uid"`
	UpdatedAt      int64  `json:"updatedAt" bson:"updatedAt"`
	Views          int64  `json:"views" bson:"views"`
}

var MetricsService metricsService

func (service *metricsService) New(client *mongo.Client, database *mongo.Database) {
	service.client = client
	service.database = database
	service.collection = service.database.Collection(metricsCollectionName)
}

func (service *metricsService) DeleteOne(
	operationContext context.Context,
	filter map[string]any,
) error {
	_, queryError := service.collection.DeleteOne(operationContext, filter)
	return queryError
}

func (service *metricsService) DeleteOneByUid(
	operationContext context.Context,
	uid string,
) error {
	_, queryError := service.collection.DeleteOne(operationContext, bson.M{"uid": uid})
	return queryError
}

func (service *metricsService) FindAll(
	operationContext context.Context,
	filter map[string]any,
	destination *[]MetricsModel,
) error {
	cursor, cursorError := service.collection.Find(operationContext, filter)
	if cursorError != nil {
		return cursorError
	}
	defer cursor.Close(context.Background())
	return cursor.All(context.Background(), &destination)
}

func (service *metricsService) FindOne(
	operationContext context.Context,
	filter map[string]any,
	destination *MetricsModel,
) error {
	return service.collection.FindOne(operationContext, filter).Decode(destination)
}

func (service *metricsService) FindOneAndUpdate(
	operationContext context.Context,
	filter map[string]any,
	update map[string]any,
	destination *MetricsModel,
) error {
	return service.collection.FindOneAndUpdate(
		operationContext,
		filter,
		update,
	).Decode(destination)
}

func (service *metricsService) FindOneByUid(
	operationContext context.Context,
	uid string,
	destination *MetricsModel,
) error {
	return service.collection.FindOne(operationContext, bson.M{"uid": uid}).Decode(destination)
}

func (service *metricsService) InsertOne(
	operationContext context.Context,
	document MetricsModel,
) error {
	_, queryError := service.collection.InsertOne(operationContext, document)
	return queryError
}

func (service *metricsService) UpdateMany(
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

func (service *metricsService) UpdateOne(
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
