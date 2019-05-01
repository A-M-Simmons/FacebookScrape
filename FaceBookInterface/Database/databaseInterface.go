package database

import (
	"fmt"

	"github.com/A-M-Simmons/FacebookScrape/datastructs"
)

type dbType struct {
	value string
}

// Photos ...
type Photos struct {
	DataJSONFileName string
	Photos           []datastructs.Photo
}

// DB ...
type DB struct {
	dbType  dbType
	Photos  Photos
	UserIDs []datastructs.User
	Posts   []datastructs.Post
}

// NewDacebookDB ...
func NewFacebookDB() DB {
	var DB DB
	DB.photosConstructor()
	return DB
}

// SetPhotos ...
func (db *DB) SetPhotos(p []datastructs.Photo) error {
	return db.Photos.setPhotos(p)
}

// GetPhotos ...
func (db *DB) GetPhotos() []datastructs.Photo {
	return db.Photos.getPhotos()
}

// photosConstructor ...
func (db *DB) photosConstructor() {
	db.Photos.DataJSONFileName = "FBPhotos.json"
}

// setPhotos ...
func (photos *Photos) setPhotos(p []datastructs.Photo) error {
	photos.Photos = p
	return nil // TODO: ERROR HANDLE OR SOMETHING
}

// getPhotos ...
func (photos *Photos) getPhotos() []datastructs.Photo {
	return photos.Photos
}

func isAllowedType(t string) bool {
	return ((t == "JSON") || (t == "SQL"))
}

// SetType ...
func (dbType *dbType) SetType(t string) error {
	if isAllowedType(t) {
		dbType.value = t
		return (nil)
	}
	return (fmt.Errorf("DB: Can't set database type too %s", t))
}

// CheckType ...
func (dbType *dbType) CheckType(t string) (bool, error) {
	// Check input
	if isAllowedType(t) == false {
		return false, fmt.Errorf("DB: Can't be type %s", t)
	}

	// Compare type
	if dbType.value == t {
		return true, nil
	} else {
		return false, nil
	}
}

// SetType ...
func (db *DB) SetType(t string) error {
	return (db.dbType.SetType(t))
}

// CheckType ...
func (db *DB) CheckType(t string) (bool, error) {
	return db.dbType.CheckType(t)
}

// SavePhotos ...
func (db *DB) SavePhotos() error {

	// Check for JSON DB Type
	flag, err := db.CheckType("JSON")
	if err != nil { // Check if it errored and return error
		return err
	} else if flag == true {
		return db.Photos.SaveJSON()
	}

	// Check for SQL DB Type
	flag, err = db.CheckType("SQL")
	if err != nil { // Check if it errored and return error
		return err
	} else if flag == true {
		return fmt.Errorf("SQL Database not yet supported")
		//return sqldatabase.SavePhotosSQL(db.Photos)
	}

	// Catchall error at end
	return fmt.Errorf("Error in SavePhotos(). This error shouldn't be reached unless you added a new databse type and didn't add it into this function")
}

// LoadPhotosJSON ...
func (db *DB) LoadPhotosJSON() error {
	err := db.Photos.LoadJSON()
	if err != nil {
		return err
	}
	return nil
}
