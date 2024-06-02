package main

import "os"

type enviroment struct {
}

func NewEnviroment() *enviroment {
	return &enviroment{}
}

func (e *enviroment) getEnv(env string, fallback string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return fallback
	}
	return value
}
