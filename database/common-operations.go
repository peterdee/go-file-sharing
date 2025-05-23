package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Operations CommonOperations

func (operations *CommonOperations) deleteOne(
	collection *mongo.Collection,
	filter bson.M,
	operationContext context.Context,
) error {
	_, queryError := collection.DeleteOne(operationContext, filter)
	return queryError
}

func (operations *CommonOperations) getOne(
	collection *mongo.Collection,
	filter bson.M,
	operationContext context.Context,
) *mongo.SingleResult {
	return collection.FindOne(operationContext, filter)
}

func (operations *CommonOperations) getOneAndUpdate(
	collection *mongo.Collection,
	filter bson.M,
	update bson.M,
	operationContext context.Context,
) *mongo.SingleResult {
	return collection.FindOneAndUpdate(
		operationContext,
		filter,
		update,
	)
}

func (operations *CommonOperations) insertOne(
	collection *mongo.Collection,
	document any,
	operationContext context.Context,
) error {
	_, queryError := collection.InsertOne(operationContext, document)
	return queryError
}

func (operations *CommonOperations) updateOne(
	collection *mongo.Collection,
	filter bson.M,
	update bson.M,
	operationContext context.Context,
) error {
	_, queryError := collection.UpdateOne(operationContext, filter, update)
	return queryError
}

func (operations *CommonOperations) DeleteFile(
	filter bson.M,
	requestContext context.Context,
) error {
	return operations.deleteOne(FilesCollection, filter, requestContext)
}

func (operations *CommonOperations) DeleteMetrics(
	filter bson.M,
	requestContext context.Context,
) error {
	return operations.deleteOne(MetricsCollection, filter, requestContext)
}

func (operations *CommonOperations) GetFile(
	filter bson.M,
	file *Files,
	requestContext context.Context,
) error {
	return operations.getOne(FilesCollection, filter, requestContext).Decode(file)
}

func (operations *CommonOperations) GetMetrics(
	filter bson.M,
	metrics *Metrics,
	requestContext context.Context,
) error {
	return operations.getOne(MetricsCollection, filter, requestContext).Decode(metrics)
}

func (operations *CommonOperations) GetMetricsAndUpdate(
	filter bson.M,
	update bson.M,
	metrics *Metrics,
	requestContext context.Context,
) error {
	return operations.getOneAndUpdate(
		MetricsCollection,
		filter,
		update,
		requestContext,
	).Decode(metrics)
}

func (operations *CommonOperations) GetUser(
	filter bson.M,
	user *Users,
	requestContext context.Context,
) error {
	return operations.getOne(UsersCollection, filter, requestContext).Decode(user)
}

func (operations *CommonOperations) GetUserAndUpdate(
	filter bson.M,
	update bson.M,
	user *Users,
	requestContext context.Context,
) error {
	return operations.getOneAndUpdate(
		UsersCollection,
		filter,
		update,
		requestContext,
	).Decode(user)
}

func (operations *CommonOperations) InsertFile(
	file any,
	requestContext context.Context,
) error {
	return operations.insertOne(FilesCollection, file, requestContext)
}

func (operations *CommonOperations) InsertMetrics(
	metrics any,
	requestContext context.Context,
) error {
	return operations.insertOne(MetricsCollection, metrics, requestContext)
}

func (operations *CommonOperations) UpdateUser(
	filter bson.M,
	update bson.M,
	requestContext context.Context,
) error {
	return operations.updateOne(UsersCollection, filter, update, requestContext)
}
