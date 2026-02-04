package models

import (
	"backend/config"
	"encoding/base64"
	"encoding/json"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

type Cursor struct {
	Column string
	Value  any
}

type Lister interface {
	GetId() uint
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

	err = db.AutoMigrate(&user{}, &Movie{})
	if err != nil {
		return DB{}, err
	}

	return DB{
		db: db,
	}, err
}

func Create(db DB, resouce any) error {
	err := db.db.Create(resouce).Error
	if err != nil {
		return err
	}

	return nil
}

// Limitar buscar (Limite máximo, limite especificado)
// Ordenar as informações
// Pular páginas
func List[T Lister, S ~[]T](db DB, list *S, cursor string) (string, error) {
	tx := db.db.Order("id").Limit(100)
	if cursor != "" {
		json_cursor, err := base64.URLEncoding.DecodeString(cursor)
		if err != nil {
			return "", err
		}

		var tmp Cursor
		err = json.Unmarshal(json_cursor, &tmp)
		if err != nil {
			return "", err
		}
		tx = tx.Where("id > ?", tmp.Value)
	}

	err := tx.Find(list).Error
	if err != nil {
		return "", err
	}

	if len(*list) > 0 {
		i := len(*list) - 1
		last := (*list)[i]
		id := last.GetId()

		byteJson, err := json.Marshal(Cursor{Value: id, Column: "id"})
		if err != nil {
			return "", err
		}

		cursor = base64.URLEncoding.EncodeToString(byteJson)
	} else {
		cursor = ""
	}

	return cursor, nil
}
