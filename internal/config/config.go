package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strconv"
)

type DBConfig struct {
	Host string `env:"DB_HOST"`
	Port uint16 `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
}

type AWSConfig struct {
	Secret string `env:"AWS_SECRET"`
	Key    string `env:"AWS_KEY"`
	Region string `env:"AWS_REGION"`
}

type CatsBucketConfig struct {
	CatsBucketName string `env:"CATS_BUCKET_NAME"`
	CatsBucketURL  string `env:"CATS_BUCKET_URL"`
}

type TgConfig struct {
	Token string `env:"TG_TOKEN"`
}

var dbConfig *DBConfig = nil
var awsConfig *AWSConfig = nil
var catsBucketConfig *CatsBucketConfig = nil
var tgConfig *TgConfig = nil

func fillByEnv(st interface{}) {
	stValue := reflect.ValueOf(st).Elem()
	fields := reflect.VisibleFields(reflect.TypeOf(st).Elem())

	for i, field := range fields {
		kind := field.Type.Kind()
		tag := field.Tag.Get("env")
		value := os.Getenv(tag)

		if len(value) == 0 {
			log.Printf("failed to parse env: %s", tag)
			continue
		}

		switch kind {
		case reflect.String:
			stValue.Field(i).SetString(value)
			break
		case reflect.Uint16:
			r, err := strconv.ParseUint(value, 10, 16)

			if err != nil {
				log.Println("failed to parse int from env")
				continue
			}

			stValue.Field(i).SetUint(r)
		}
	}
}

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("failed to load env")
	}
}

func GetDBConfig() *DBConfig {
	if dbConfig == nil {
		dbConfig = &DBConfig{}
		fillByEnv(dbConfig)
	}

	return dbConfig
}

func GetAWSConfig() *AWSConfig {
	if awsConfig == nil {
		awsConfig = &AWSConfig{}
		fillByEnv(awsConfig)
	}

	return awsConfig
}

func GetCatsBucketConfig() *CatsBucketConfig {
	if catsBucketConfig == nil {
		catsBucketConfig = &CatsBucketConfig{}
		fillByEnv(catsBucketConfig)
	}

	return catsBucketConfig
}

func GetTgConfig() *TgConfig {
	if tgConfig == nil {
		tgConfig = &TgConfig{}
		fillByEnv(tgConfig)
	}

	return tgConfig
}
