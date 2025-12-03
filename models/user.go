package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;not null"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Hash      []byte
	ResetCode []byte
	ResetDate *time.Time `gorm:"default:null"`
}

func CreateUser(db DB, username string, password []byte, email string) error {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := user{
		Username: username,
		Hash:     hash,
		Email:    email,
	}

	err = db.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil

}

func ValidUser(db DB, username string, password []byte) error {
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

func UpdateUserPasswordReset(db DB, email string) (pwd string, username string, err error) {
	var pwd_raw [16]byte
	rand.Read(pwd_raw[:])
	pwd = base64.URLEncoding.EncodeToString(pwd_raw[:])
	user_got := user{}
	// Take: Pega a primeira coluna que der match
	err = db.db.Where(user{Email: email}).Take(&user_got).Error
	if err != nil {
		return
	}

	date := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	err = db.db.Where(user{ID: user_got.ID}).Updates(user{ResetCode: hash, ResetDate: &date}).Error
	if err != nil {
		return
	}

	username = user_got.Username

	return
}

func ResetPassword(db DB, email string, new_password string, hash_code string) error {
	user_temp := user{}
	err := db.db.Where(user{Email: email}).Take(&user_temp).Error

	if err != nil {
		return err
	}

	if user_temp.ResetDate == nil {
		return errors.New("no pending reset password")
	}

	if user_temp.ResetDate.AddDate(0, 0, 1).Before(time.Now()) {
		return errors.New("password reset expired")
	}

	err = bcrypt.CompareHashAndPassword(user_temp.ResetCode, []byte(hash_code))
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.db.Where(user{Email: email}).Select("Hash", "ResetDate", "ResetCode").Updates(user{Hash: hash, ResetDate: nil, ResetCode: nil}).Error
	if err != nil {
		return err
	}

	return nil
}
