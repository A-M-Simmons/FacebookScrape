package fb

import (
	"fmt"
)

// PrintTitle ...
func (s *Session) PrintTitle() {
	fmt.Println(s.browser.Title())
}

// Open ...
func (s *Session) Open(url string) error {
	return s.browser.Open(url)
}
