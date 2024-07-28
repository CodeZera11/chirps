package database

import "os"

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID: id,
		Body: body,
	}

	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}


func (db *DB) GetOneChirp(id int) (Chirp, error){
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{}

	for _, dbChirp := range dbStructure.Chirps {
		if dbChirp.ID == id{
			chirp = Chirp{
				ID: dbChirp.ID,
				Body: dbChirp.Body,
			}
		}
	}

	if chirp.ID == 0 {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}