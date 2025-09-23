package env

import "os"



func Getenv(key,defaultValue string) string{
	if value,exists := os.LookupEnv(key); exists{
		return value
	}

	return defaultValue
}