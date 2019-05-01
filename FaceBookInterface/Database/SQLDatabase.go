package database

import (
	"database/sql"
	"fmt"
	"log"
)

type databaseCredentials struct {
	endpoint string
	port     string
	user     string
	pass     string
	dbname   string
}

func (dbC *databaseCredentials) setEndpoint(str string) {
	dbC.endpoint = str
}
func (dbC *databaseCredentials) setPort(str string) {
	dbC.port = str
}
func (dbC *databaseCredentials) setUser(str string) {
	dbC.user = str
}
func (dbC *databaseCredentials) setPass(str string) {
	dbC.pass = str
}
func (dbC *databaseCredentials) setDBName(str string) {
	dbC.dbname = str
}

// Credentials saves the Database credentials
func (db *DBSQL) Credentials(user, pass, ep, port, dbn string) {
	db.creds.setUser(user)
	db.creds.setPass(pass)
	db.creds.setEndpoint(ep)
	db.creds.setPort(port)
	db.creds.setDBName(dbn)
}

// DB Main Database structure to handle connections
type DBSQL struct {
	conn     *sql.DB
	creds    databaseCredentials
	tables   map[string]string
	connOpen bool
}

// NewDatabase Initalise the Database structure and returns it
func NewDatabase() *DBSQL {
	var db DBSQL
	db.connOpen = false
	db.tables = make(map[string]string)
	db.tables["Logins"] = "Logins"
	db.tables["Photos"] = "Photos"
	db.tables["Comments"] = "Comments"

	return &db
}

func (db *DBSQL) Init() {
	db.createLoginTable()
	db.createPhotosTable()
	db.createCommentsTable()
}

// CreateLoginDatabase initialises the Login table
func (db *DBSQL) createLoginTable() (sql.Result, error) {
	command := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Account varchar(64) PRIMARY KEY, Username varchar(64), Pass varchar(64));", db.tables["Logins"])
	resp, err := db.conn.Exec(command)
	return resp, err
}

// createPhotosDatabase initialises the Photos table
func (db *DBSQL) createPhotosTable() (sql.Result, error) {
	command := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (URL varchar(1000), Uploader varchar(64), Uploaded varchar(20), Description varchar(8000), PRIMARY KEY (URL));", db.tables["Photos"])
	resp, err := db.conn.Exec(command)
	return resp, err
}

// createCommentsTable initialises the CommentThreads table
func (db *DBSQL) createCommentsTable() (sql.Result, error) {
	command := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (PhotoURL varchar(1000), CommentThreadPos INT, CommentPos INT, Uploader varchar(64), Uploaded varchar(20), Message varchar(8000), PRIMARY KEY (PhotoURL, CommentThreadPos, CommentPos));", db.tables["Comments"])
	resp, err := db.conn.Exec(command)
	return resp, err
}

// AddLogin ...
func (db *DBSQL) AddLogin(name, user, pass string) (sql.Result, error) {
	command := fmt.Sprintf("INSERT IGNORE INTO %s (Account, Username, Pass) VALUES ('%s', '%s', '%s');", db.tables["Logins"], name, user, pass)
	resp, err := db.conn.Exec(command)
	return resp, err
}

// GetLogin ...
func (db *DBSQL) GetLogin(name string) (string, string, error) {
	var user string
	var pass string
	command := fmt.Sprintf("SELECT Username, Pass FROM %s WHERE Account = '%s';", db.tables["Logins"], name)
	err := db.conn.QueryRow(command).Scan(&user, &pass)

	return user, pass, err
}

// Open opens a connection to database
func (db *DBSQL) Open() {
	loginStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.creds.user, db.creds.pass, db.creds.endpoint, db.creds.port, db.creds.dbname)
	dbSQL, err := sql.Open("mysql", loginStr)
	if err != nil {
		log.Fatal(err)
	}
	db.conn = dbSQL
}

//
func (db *DBSQL) AddCommentThreads() int64 {
	command := fmt.Sprintf("INSERT INTO CommentThreads VALUES (NULL);")
	db.conn.Exec(command) //TODO: Error handle
	var ID int64
	err := db.conn.QueryRow("SELECT LAST_INSERT_ID();").Scan(&ID)
	if err != nil {
		panic(err)
	}
	return ID
}

func (db *DBSQL) AddCommentThread() int64 {
	command := fmt.Sprintf("INSERT INTO CommentThreads VALUES (NULL);")
	db.conn.Exec(command) //TODO: Error handle
	var ID int64
	err := db.conn.QueryRow("SELECT LAST_INSERT_ID();").Scan(&ID)
	if err != nil {
		panic(err)
	}
	return ID
}
