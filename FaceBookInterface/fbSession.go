package fb

import (
	database "github.com/A-M-Simmons/FacebookScrape/FaceBookInterface/Database"
	"github.com/headzoo/surf/browser"
)

// Session stores the Logindetails and Browser
type Session struct {
	browser *browser.Browser
	DB      database.DB
	login   LoginDetails
}

// NewQuery ...
func (session *Session) NewQuery() database.DBJSONQuery {
	var dbQuery database.DBJSONQuery
	dbQuery.DB = session.DB
	return dbQuery
}

// Login logs into facebook.
func (session *Session) Login() error {
	bow, err := session.login.login()
	if err != nil {
		return err
	}
	session.browser = bow
	return nil
}

// SetUsername is the Username setter function for the LoginDetails structure
func (session *Session) SetUsername(str string) error {
	session.login.setUsername(str)
	return nil
}

// SetPassword is the Password setter function for the LoginDetails structure
func (session *Session) SetPassword(str string) error {
	session.login.setPassword(str)
	return nil
}

// NewSession makes a new Facebook Session to store the browser and login details
func NewSession() Session {
	var s Session
	s.DB = database.NewFacebookDB()
	return s
}
