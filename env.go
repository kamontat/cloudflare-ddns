package main

import (
	"log"
	"os"
	"path/filepath"
)

type Env struct {
	ApiToken   string
	ZoneId     string
	ConfigFile string
}

func GetEnv() *Env {
	var wd, err = os.Getwd()
	if err != nil {
		log.Panicln("cannot get current working directory")
	}

	return &Env{
		ApiToken:   requireEnv("CF_API_TOKEN"),
		ZoneId:     requireEnv("CF_ZONE_ID"),
		ConfigFile: getEnv("CF_DDNS__CONFIG_FILE", filepath.Join(wd)),
	}
}

func requireEnv(name string) string {
	var value = os.Getenv(name)
	if value == "" {
		log.Panicf("Environment $%s is required", name)
	}

	return value
}

func getEnv(name string, def string) string {
	var value = os.Getenv(name)
	if value == "" {
		return def
	}

	return value
}
