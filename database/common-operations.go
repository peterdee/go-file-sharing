package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Operations CommonOperations

func (operations *CommonOperations) getOne(
	collection *mongo.Collection,
	filter bson.M,
	target any,
	operationContext context.Context,
) error {
	return collection.FindOne(
		operationContext,
		filter,
	).Decode(&target)
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

func (operations *CommonOperations) GetUser(
	filter bson.M,
	user any,
	requestContext context.Context,
) error {
	return operations.getOne(UsersCollection, filter, user, requestContext)
}

func (operations *CommonOperations) UpdateUser(
	filter bson.M,
	update bson.M,
	requestContext context.Context,
) error {
	return operations.updateOne(UsersCollection, filter, update, requestContext)
}
