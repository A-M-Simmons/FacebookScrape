package session

import (
	fb "github.com/A-M-Simmons/FacebookScrape/FaceBookInterface"
)

// Session ...
type Session struct {
	FBSession fb.Session
}

// NewSession makes a new Facebook Session to store the browser and login details
func NewSession() *Session {
	var s Session
	s.FBSession = fb.NewSession()
	return &s
}
