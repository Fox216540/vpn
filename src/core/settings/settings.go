package settings

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Setting struct {
	FilePath         string
	ClientConfigPath string
	HomePath         string
	ScriptName       string
}

func NewSetting() *Setting {
	err := godotenv.Load() // по умолчанию ищет файл .env в текущей папке
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return &Setting{
		FilePath:         os.Getenv("FILE_PATH"),
		ClientConfigPath: os.Getenv("CLIENT_CONFIG_PATH"),
		HomePath:         os.Getenv("HOME_PATH"),
		ScriptName:       os.Getenv("SCRIPT_NAME"),
	}
}

var Config = NewSetting()
