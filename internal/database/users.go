package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (RespUser, error) {
	userResp:= RespUser{}
	dbStructure, err := db.loadDB()

	if err != nil {
		return userResp, err
	}

	_, err = db.GetOneUser(email)

	if err == nil {
		return RespUser{}, errors.New("user with this email already exists")
	}

	id := len(dbStructure.Users) + 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 1)

	if err != nil {
		return RespUser{}, err
	}

	userInput := User{
		ID: id,
		Email: email,
		Password: string(hashedPassword),
	}

	dbStructure.Users[id] = userInput

	err = db.writeDB(dbStructure)

	if err != nil {
		return RespUser{}, err
	}

	userResp = RespUser{
		ID: id,
		Email: email,
	}
	
	return userResp, nil
}

func (db *DB) GetOneUser(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users

	for _, user := range users {
		if user.Email == email {
			// return RespUser{
			// 	Email: user.Email,
			// 	ID: user.ID,
			// }, nil
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) LoginUser(email, password string) (RespUser, error) {
	// dbStructure, err := db.loadDB()

	// if err != nil {
	// 	return RespUser{}, err
	// }

	user, err := db.GetOneUser(email)

	if err != nil {
		return RespUser{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return RespUser{}, err
	}

	respUser := RespUser{
		Email: user.Email,
		ID: user.ID,
	}

	return respUser, nil
}