package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"file-sharing/constants"
	"file-sharing/utilities"
)

var Client *mongo.Client

var Database *mongo.Database

var ErrNoDocuments = mongo.ErrNoDocuments

func Connect() {
	connectionString := utilities.GetEnv(constants.ENV_NAMES.DatabaseConnectionString)
	if connectionString == "" {
		log.Fatal("Database connection string was not provided via environment variable")
	}
	databaseName := utilities.GetEnv(
		constants.ENV_NAMES.DatabaseName,
		constants.DEFAULT_DATABASE_NAME,
	)

	for i := 1; i <= 6; i += 1 {
		client, connectionError := mongo.Connect(options.Client().ApplyURI(connectionString))
		if connectionError != nil {
			log.Printf("MongoDB connection failed, retry in %d seconds", i)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		if i == 6 {
			log.Fatal(connectionError)
		}
		Client = client
		break
	}

	ctx := context.Background()
	for i := 1; i <= 6; i += 1 {
		pingError := Client.Ping(ctx, readpref.Primary())
		if pingError == nil {
			break
		}
		log.Printf("MongoDB ping failed, retry in %d seconds", i)
		time.Sleep(time.Duration(i) * time.Second)
	}

	Database = Client.Database(databaseName)

	FileService.New(Client, Database)
	MetricsService.New(Client, Database)
	UserService.New(Client, Database)

	seeding()

	log.Println("MongoDB connection is ready")
}
