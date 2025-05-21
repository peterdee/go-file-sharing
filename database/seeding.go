package database

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/julyskies/gohelpers"
	"github.com/nrednav/cuid2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/utilities"
)

func seeding() {
	rootEmail := utilities.GetEnv(constants.ENV_NAMES.RootEmail)
	if rootEmail == "" {
		log.Fatalf("missing required %s environment variable", constants.ENV_NAMES.RootEmail)
	}

	var rootUser Users
	queryError := UsersCollection.FindOne(
		context.Background(),
		bson.M{
			"email": rootEmail,
		},
	).Decode(&rootUser)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			timestamp := gohelpers.MakeTimestampSeconds()
			_, queryError = UsersCollection.InsertOne(
				context.Background(),
				Users{
					CreatedAt:      timestamp,
					DeletedAt:      0,
					Email:          strings.ToLower(strings.Trim(rootEmail, " ")),
					IsDeleted:      false,
					PasswordHash:   "",
					Role:           constants.ROLES.Root,
					SetUpCompleted: false,
					UID:            cuid2.Generate(),
					UpdatedAt:      timestamp,
				},
			)
			if queryError != nil {
				log.Fatal(queryError)
			}
			log.Println("Created root account record")
			return
		}
		log.Fatal(queryError)
	}
}
