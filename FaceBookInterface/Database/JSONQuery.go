package database

import (
	"fmt"

	"github.com/A-M-Simmons/FacebookScrape/datastructs"
)

// DBJSONQuery ...
type DBJSONQuery struct {
	DB DB
}

func compareUserID(author string, u datastructs.User) bool {
	return u.ID == author || u.UserName == author || u.RealName == author
}

// GetPhotosPostedBy ...
func (db *DBJSONQuery) GetPhotosPostedBy(author string) {
	photos := db.DB.Photos.getPhotos()
	for _, photo := range photos {
		if compareUserID(author, photo.UserID) {
			fmt.Printf("%s posted a photo on %s\n", author, photo.DateTime)
		}
	}
}

// GetPhotosCommentedOnBy ...
func (db *DBJSONQuery) GetPhotosCommentedOnBy(author string) {
	photos := db.DB.Photos.getPhotos()
	for _, photo := range photos {
		for _, c := range photo.Comments {
			if compareUserID(author, c.UserID) {
				fmt.Printf("%s commented on the %s to a photo posted on the %s\n", author, c.DateTime, photo.DateTime)
				break
			}
			for _, cc := range c.Comments {
				if compareUserID(author, cc.UserID) {
					fmt.Printf("%s replied to a comment on the %s to a photo posted on the %s\n", author, cc.DateTime, photo.DateTime)
					break
				}
			}
		}
	}
}
