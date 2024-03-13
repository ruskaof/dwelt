package config

import (
	"dwelt/src/utils"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	DbName   string `envconfig:"DB_NAME"`
	User     string
	Password string
}

type DweltConfig struct {
	JwtKey            string `envconfig:"JWT_KEY"`
	WorkflowRunNumber int    `envconfig:"WORKFLOW_RUN_NUMBER"`
}

var (
	DbCfg    DatabaseConfig
	DweltCfg DweltConfig
)

func InitCfg() {
	_ = godotenv.Load("local.env")

	utils.MustNoErr(envconfig.Process("db", &DbCfg))
	utils.MustNoErr(envconfig.Process("dwelt", &DweltCfg))

	slog.Debug("Loaded configs: " +
		string(utils.Must(json.Marshal(DbCfg))) +
		" " +
		string(utils.Must(json.Marshal(DweltCfg))))
}
