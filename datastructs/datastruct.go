package datastructs

import (
	"fmt"
	"html"
	"strconv"
)

// USER STRUCTURE AND FUNCTIONS
type FBProfile struct {
	UserID User
	Photos []Photo
	Posts  []Post
}

// User ...
type User struct {
	ID       string
	UserName string
	RealName string
}

// SetID ...
func (u *User) SetID(ID string) {
	u.ID = ID
}

// SetUsername ...
func (u *User) SetUsername(username string) {
	u.UserName = username
}

// SetRealName ...
func (u *User) SetRealName(realName string) {
	u.RealName = realName
}

// PHOTO STRUCTURE AND FUNCTIONS

// Photo ...
type Photo struct {
	URL         string
	DateTime    string
	Description string
	Filepath    string
	UserID      User
	Comments    []Comment
	ReactionSet []Reaction
}

// PrintCommentsPretty ...
func (p *Photo) PrintCommentsPretty() {
	for _, ct := range p.Comments {
		for _, com := range ct.Comments {
			fmt.Println(fmt.Sprintf("%s %s: %s", com.DateTime, com.UserID.RealName, com.Message.Value))
		}
		fmt.Println("")
	}
}

func (p *Photo) addUploader(name string) {
	var u User
	u.RealName = name
	p.UserID = u
}

// Post ...
type Post struct {
	UserID      User
	DateTime    string
	Message     MessageString
	Comments    []Comment
	ReactionSet []Reaction
}

// Reaction ...
type Reaction struct {
	ReactionType string
	UserID       User
}

// COMMENTS STRUCTURE AND FUNCTIONS

// Comment ...
type Comment struct {
	UserID      User
	Message     MessageString
	ReactionSet []Reaction
	DateTime    string
	Comments    []Comment
}

// AddComment ...
func (c *Comment) SetComment(str string) error {

	str = html.UnescapeString(str)
	str, err := strconv.Unquote("`" + str + "`")
	fmt.Println(err)
	c.Message.Value = str
	return nil
}

// AddTimestamp ...
func (c *Comment) SetTimestamp(time string) {
	c.DateTime = time
}

// Append ...
func (ct *Comment) Append(c Comment) {
	ct.Comments = append(ct.Comments, c)
}

// MessageString ...
type MessageString struct {
	Value      string
	PersonRefs []User
}
