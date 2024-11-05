package utilities

import "os"

func GetEnv(envName string, defaultValue ...string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		if len(defaultValue) == 0 {
			return ""
		}
		return defaultValue[0]
	}
	return envValue
}
