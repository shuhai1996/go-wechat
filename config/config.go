package config

import (
	"github.com/verystar/ini"
	"log"
	"strconv"
)

var cfg *ini.Ini
var env string

func Setup(environment string) {
	var err error
	cfg, err = ini.Load("../conf.ini")
	if err != nil {
		log.Fatalln(err)
	}
	env = environment
}

func Get(key string) string {
	return cfg.Read(env, key)
}

func GetInt(key string) int {
	i, _ := strconv.Atoi(cfg.Read(env, key))
	return i
}

func Getenv() string {
	return env
}
