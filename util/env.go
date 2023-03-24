package util

import (
	"fmt"
	"os"
)

func GetEnv(name, defaultValue, help string) string {
	fmt.Printf("LOAD ENVIRONMENT: %v, DEFAULT VALUE: %v, DESC: %v \n", name, defaultValue, help)
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return defaultValue
}
