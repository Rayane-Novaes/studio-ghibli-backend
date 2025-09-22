package models

import (
	"backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func ConnectDb(cfg config.Config) (DB, error) {

	pg := postgres.New(
		postgres.Config{
			DSN: "host=" + cfg.Host + " port=" + cfg.Port + " user=" + cfg.User_DB + " password=" + cfg.Password_DB + " dbname=" + cfg.Db_name,
		},
	)

	db, err := gorm.Open(pg, &gorm.Config{})
	if err != nil {
		return DB{}, err
	}

	err = db.AutoMigrate(&user{})
	if err != nil{
		return DB{}, err
	}

	return DB{
		db: db,
	}, err
}

