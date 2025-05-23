package database

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/julyskies/gohelpers"
	"github.com/nrednav/cuid2"

	"file-sharing/constants"
	"file-sharing/utilities"
)

func seeding() {
	rootEmail := utilities.GetEnv(constants.ENV_NAMES.RootEmail)
	if rootEmail == "" {
		log.Fatalf("missing required %s environment variable", constants.ENV_NAMES.RootEmail)
	}

	var rootUser UserModel
	queryError := UserService.FindOne(
		context.Background(),
		map[string]any{"email": rootEmail},
		&rootUser,
	)
	if queryError != nil {
		if errors.Is(queryError, ErrNoDocuments) {
			timestamp := gohelpers.MakeTimestampSeconds()
			queryError := UserService.InsertOne(
				context.Background(),
				UserModel{
					CreatedAt:      timestamp,
					DeletedAt:      0,
					Email:          strings.ToLower(strings.Trim(rootEmail, " ")),
					IsDeleted:      false,
					PasswordHash:   "",
					Role:           constants.ROLES.Root,
					SetUpCompleted: false,
					Uid:            cuid2.Generate(),
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
