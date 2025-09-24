package settings

import (
	"github.com/joho/godotenv"
	"os"
)

type Setting struct {
	FilePath            string
	HomePath            string
	ScriptName          string
	FileName            string
	HashPass            string
	TrafficInterval     string
	CPUInterval         string
	MemoryInterval      string
	ConnectionsInterval string
}

func NewSetting() *Setting {
	_ = godotenv.Load() // игнорируем ошибку
	return &Setting{
		FilePath:            os.Getenv("FILE_PATH"),
		HomePath:            os.Getenv("HOME_PATH"),
		ScriptName:          os.Getenv("SCRIPT_NAME"),
		FileName:            os.Getenv("FILE_NAME"),
		HashPass:            os.Getenv("HASH_PASS"),
		TrafficInterval:     os.Getenv("TRAFFIC_INTERVAL"),
		CPUInterval:         os.Getenv("CPU_INTERVAL"),
		MemoryInterval:      os.Getenv("MEMORY_INTERVAL"),
		ConnectionsInterval: os.Getenv("CONNECTIONS_INTERVAL"),
	}
}

var Config = NewSetting()
