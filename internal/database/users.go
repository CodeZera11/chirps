package database

import (
	"errors"
	"time"

	"github.com/codezera11/chirps/internal/auth"
)

func (db *DB) CreateUser(email, password string) (User, error) {
	userResp:= User{}
	dbStructure, err := db.loadDB()

	if err != nil {
		return userResp, err
	}

	_, err = db.GetUserByEmail(email)

	if err == nil {
		return User{}, errors.New("user with this email already exists")
	}

	id := len(dbStructure.Users) + 1

	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		return User{}, err
	}

	user := User{
		ID: id,
		Email: email,
		Password: hashedPassword,
	}

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}
	
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) UpdateUser(email string, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]

	if !ok {
		return User{}, errors.New("User does not exist")
	}

	if user.ID != id {
		return User{}, errors.New("Unauthorizezed")
	}

	user.Email = email
	user.Password = password

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}

	return User{
		Email: user.Email,
		ID: user.ID,
	}, nil
}

func (db *DB) AddRefTokenToUser(id int, refToken string, expirationTime time.Time) error {

	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]

	if !ok {
		return ErrNotFound
	}

	user.RefreshToken = refToken
	user.ExpiresAt = expirationTime

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserByToken(token string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken == token {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) GetNewAccessToken(token, secret string) (string, error) {
	user, err := db.GetUserByToken(token)

	if err != nil {
		return "", err
	}

	jwtToken, err := auth.MakeJWT(user.ID, secret, time.Duration(60 * 60) * time.Second)

	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (db *DB) RevokeRefreshToken(token, secret string) error {

	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	for _, dbUser := range dbStructure.Users {
		if dbUser.RefreshToken == token {
			dbUser.RefreshToken = ""
			dbStructure.Users[dbUser.ID] = dbUser
			err := db.writeDB(dbStructure)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("error revoking refresh token")
}