package main

import "os"

func GetEnv(env string, fallback string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return fallback
	}
	return value
}
