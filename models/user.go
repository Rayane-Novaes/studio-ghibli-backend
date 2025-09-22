package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

type user struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique"`
	Hash []byte
} 

func CreateUser (db DB, username string, password []byte ) (error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := user{
		Username: username,
		Hash: hash,
	}

	err = db.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil

}

func ValidUser (db DB, username string, password []byte) (error) {
	users := []user{}

	err := db.db.Where(user{Username: username}).Find(&users).Error

	if err != nil {
		return err
	}

	if len(users) <= 0 {
		return errors.New("user ou password incorreto")
	}

	err = bcrypt.CompareHashAndPassword(users[0].Hash, password)
	if err != nil {
		return err
	}
	
	return nil
}