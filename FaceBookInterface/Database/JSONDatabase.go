package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/A-M-Simmons/FacebookScrape/datastructs"
)

// SaveJSON ...
func (photos *Photos) SaveJSON() error {
	file, err := json.MarshalIndent(photos.getPhotos(), "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(photos.DataJSONFileName, file, 0644) // Change Hardcoded FileName
}

// LoadJSON ...
func (photos *Photos) LoadJSON() error {
	plan, _ := ioutil.ReadFile(photos.DataJSONFileName) // Change Hardcoded FileName
	var data []datastructs.Photo
	err := json.Unmarshal([]byte(plan), &data)
	if err != nil {
		fmt.Println("oh no")
		fmt.Println(err)
		return err
	}
	photos.setPhotos(data)
	return err
}
