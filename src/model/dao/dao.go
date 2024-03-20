package dao

import (
	"dwelt/src/config"
	"dwelt/src/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DbCfg.Host,
		config.DbCfg.User,
		config.DbCfg.Password,
		config.DbCfg.DbName,
		config.DbCfg.Port,
	)

	Db = utils.Must(
		gorm.Open(postgres.Open(dsn), &gorm.Config{
			TranslateError: true,
			// log every SQL command
			Logger: logger.Default.LogMode(logger.Info),
		}),
	)
}
