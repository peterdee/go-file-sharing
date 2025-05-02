package constants

type EnvNames struct {
	DatabaseConnectionString string
	DatabaseName             string
	Port                     string
	RedisHost                string
	RedisPassword            string
}

type ResponseInfo struct {
	InternalServerError string
	Ok                  string
}
