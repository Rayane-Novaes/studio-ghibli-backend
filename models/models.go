package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct{
	db *gorm.DB
}

func ConnectDb() (DB, error) {

	pg := postgres.New(
		postgres.Config{
			DSN: "host=localhost port=5555 user=user password=password dbname=ghibliApi",
		},
	)

	db, err := gorm.Open(pg, &gorm.Config{})
	if err != nil{
		return DB{}, err
	}

	return DB{
		db : db,
	}, err
}
