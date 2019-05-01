package fb

import (
	"fmt"

	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

type user struct {
	value string
}

type pass struct {
	value string
}

// LoginDetails holds the Username and Password for a Facebook login session.
// Access is handled by the Getter Functions 'Username' and 'Password'
// and the Setter functions 'SetUsername' and 'SetPassword'
type LoginDetails struct {
	username user
	password pass
}

// SetUsername is the Username setter function for the LoginDetails structure
func (l *LoginDetails) setUsername(s string) error {
	l.username.value = s
	return nil
}

// SetPassword is the Password setter function for the LoginDetails structure
func (l *LoginDetails) setPassword(s string) error {
	l.password.value = s
	return nil
}

// Username is the Username setter function for the LoginDetails structure
func (l LoginDetails) getUsername() (string, error) {
	if l.username.value == "" {
		return "", fmt.Errorf("No username stored in session")
	}
	return l.username.value, nil
}

// Password is the Password setter function for the LoginDetails structure
func (l LoginDetails) getPassword() (string, error) {
	if l.password.value == "" {
		return "", fmt.Errorf("No password stored in session")
	}
	return l.password.value, nil
}

// Login logs into the mobile optimised version of facebook using the supplied Username and Password.
// Sets the browser to https://m.facebook.com/home.php if completed successfully
func (l LoginDetails) login() (*browser.Browser, error) {
	bow := surf.NewBrowser()
	err := bow.Open("https://m.facebook.com/home.php")
	if err != nil {
		panic(err)
		return nil, err
	}

	fm, err := bow.Form("form#login_form")
	if err != nil {
		panic(err)
		return nil, err
	}

	username, err := l.getUsername()
	if err != nil {
		return nil, err
	}
	err = fm.Input("email", username)
	if err != nil {
		panic(err)
		return nil, err
	}

	password, err := l.getPassword()
	if err != nil {
		return nil, err
	}
	err = fm.Input("pass", password)
	if err != nil {
		panic(err)
		return nil, err
	}

	err = fm.Submit()
	if err != nil {
		panic(err)
		return nil, err
	}

	err = bow.Open("https://m.facebook.com/home.php")
	if err != nil {
		panic(err)
		return nil, err
	}

	return bow, err
}
