package constants

type EnvNames struct {
	DatabaseConnectionString string
	DatabaseName             string
	Port                     string
	RedisHost                string
	RedisPassword            string
	UplaodsDirectoryName     string
}

type ResponseInfo struct {
	InternalServerError string
	NotFound            string
	Ok                  string
}
