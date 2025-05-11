package constants

type EnvNames struct {
	DatabaseConnectionString string
	DatabaseName             string
	MaxFileSizeBytes         string
	Port                     string
	RedisHost                string
	RedisPassword            string
	UplaodsDirectoryName     string
}

type ResponseInfo struct {
	FileNotAvailable      string
	InternalServerError   string
	NotFound              string
	Ok                    string
	RequestEntityTooLarge string
}
