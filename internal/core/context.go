package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type Context struct {
	S *Settings
	D *sql.DB
}

type Settings struct {
	DbHost      string
	DbPort      int
	DbName      string
	DbUser      string
	DbPassword  string
	MaxDbConnns int `yaml:"max_db_connections"`

	ListenPort int

	CertificateFile string
	PrivateKeyFile  string
	AuthKey         string
}

func NewContext() *Context {
	s := newSettings("settings.yaml")

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			s.DbUser, s.DbPassword, s.DbHost, s.DbPort, s.DbName))
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(s.MaxDbConnns)

	return &Context{
		S: s,
		D: db,
	}
}

func newSettings(settingsPath string) *Settings {
	s := &Settings{}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", settingsPath, err)
	}

	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", settingsPath, err)
	}

	s.DbHost = os.Getenv("SPIRE_DB_HOST")

	port, err := strconv.Atoi(os.Getenv("SPIRE_DB_PORT"))
	if err != nil {
		panic(err)
	}
	s.DbPort = port

	s.DbName = os.Getenv("SPIRE_DB_NAME")

	s.DbUser = os.Getenv("SPIRE_DB_USER")

	data, err = os.ReadFile(os.Getenv("SPIRE_DB_PASSWORD_FILE"))
	if err != nil {
		panic(err)
	}
	s.DbPassword = strings.TrimSpace(string(data))

	port, err = strconv.Atoi(os.Getenv("SPIRE_LOBBY_PORT"))
	if err != nil {
		panic(err)
	}
	s.ListenPort = port

	data, err = os.ReadFile(os.Getenv("SPIRE_AUTH_KEY_FILE"))
	if err != nil {
		panic(err)
	}
	s.AuthKey = strings.TrimSpace(string(data))

	s.CertificateFile = os.Getenv("SPIRE_LOBBY_CERTIFICATE_FILE")
	s.PrivateKeyFile = os.Getenv("SPIRE_LOBBY_PRIVATE_KEY_FILE")

	if !s.validate() {
		panic("Invalid settings")
	}

	return s
}

func (s *Settings) validate() bool {
	if (s.DbHost == "") || (s.DbName == "") || (s.DbUser == "") || (s.DbPassword == "") {
		return false
	}
	if s.AuthKey == "" {
		return false
	}

	return true
}
