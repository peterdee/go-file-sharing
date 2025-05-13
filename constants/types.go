package constants

type EnvNames struct {
	DatabaseConnectionString string
	DatabaseName             string
	JwtExpirationSeconds     string
	JwtSectet                string
	MaxFileSizeBytes         string
	Port                     string
	RedisHost                string
	RedisPassword            string
	RootEmail                string
	UplaodsDirectoryName     string
}

type ResponseInfo struct {
	BadRequest            string
	FileNotAvailable      string
	InternalServerError   string
	NotFound              string
	Ok                    string
	RequestEntityTooLarge string
}

type Roles struct {
	Manager string
	Root    string
}
