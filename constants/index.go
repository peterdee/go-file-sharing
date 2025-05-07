package constants

const DEFAULT_DATABASE_NAME string = "fileshare"

const DEFAULT_PORT string = "9090"

const DEFAULT_REDIS_HOST string = "localhost:6379"

const DEFAULT_UPLOADS_DIRECTORY_NAME string = "uploads"

var ENV_NAMES = EnvNames{
	DatabaseConnectionString: "DATABASE_CONNECTION_STRING",
	DatabaseName:             "DATABASE_NAME",
	MaxFileSizeBytes:         "MAX_FILE_SIZE_BYTES",
	Port:                     "PORT",
	RedisHost:                "REDIS_HOST",
	RedisPassword:            "REDIS_PASSWORD",
	UplaodsDirectoryName:     "UPLOADS_DIRECTORY_NAME",
}

var RESPONSE_INFO = ResponseInfo{
	InternalServerError:   "INTERNAL_SERVER_ERROR",
	NotFound:              "NOT_FOUND",
	Ok:                    "OK",
	RequestEntityTooLarge: "REQUEST_ENTITY_TOO_LARGE",
}
