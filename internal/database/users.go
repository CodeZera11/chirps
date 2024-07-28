package database


func (db *DB) CreateUser(email string) (User, error) {
	user := User{}
	dbStructure, err := db.loadDB()

	if err != nil {
		return user, err
	}

	id := len(dbStructure.Users) + 1

	user = User{
		ID: id,
		Email: email,
	}

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}
	
	return user, nil
}