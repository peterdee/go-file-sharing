package constants

type EnvNames struct {
	DatabaseConnectionString string
	DatabaseName             string
	IsDockerImage            string
	JwtExpirationSeconds     string
	JwtSectet                string
	MaxFileSizeBytes         string
	Port                     string
	RedisHost                string
	RedisPassword            string
	RedisPort                string
	RootEmail                string
	UplaodsDirectoryName     string
}

type ResponseInfo struct {
	BadRequest            string
	FileNotAvailable      string
	Forbidden             string
	InternalServerError   string
	NotFound              string
	Ok                    string
	RequestEntityTooLarge string
	Unauthorized          string
}

type Roles struct {
	Manager string
	Root    string
}
