package database

import "os"

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID: id,
		Body: body,
		AuthorId: authorId,
	}

	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps(authorId int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	if authorId != 0 {
		var chirpsWithAuthorId []Chirp
		for _, dbChirp := range dbStructure.Chirps {
			if dbChirp.AuthorId == authorId {
				chirpsWithAuthorId = append(chirpsWithAuthorId, dbChirp)
			}
		}
		chirps = chirpsWithAuthorId
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
				AuthorId: dbChirp.AuthorId,
			}
		}
	}

	if chirp.ID == 0 {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpId, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[chirpId]

	if !ok {
		return Chirp{}, ErrNotFound
	}

	delete(dbStructure.Chirps, chirp.ID)

	return chirp, nil
}