package core

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Context struct {
	S *Settings

	client *mongo.Client
}

type Settings struct {
	MongoHost     string
	MongoPort     int
	MongoDatabase string
	MongoUser     string
	MongoPassword string

	ListenPort int

	CertificateFile string
	PrivateKeyFile  string
	AuthKey         string
}

func NewContext() *Context {
	s := newSettings()

	client, err := mongo.Connect(options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		s.MongoUser, s.MongoPassword, s.MongoHost, s.MongoPort, s.MongoDatabase)))
	if err != nil {
		panic(err)
	}

	return &Context{
		S:      s,
		client: client,
	}
}

func (c *Context) Close() {
	_ = c.client.Disconnect(context.Background())
}

func (c *Context) Collection(collection string) *mongo.Collection {
	return c.client.Database(c.S.MongoDatabase).Collection(collection)
}

func (c *Context) StartSession() (*mongo.Session, error) {
	return c.client.StartSession()
}

func newSettings() *Settings {
	s := &Settings{}

	s.MongoHost = os.Getenv("SPIRE_MONGO_HOST")

	port, err := strconv.Atoi(os.Getenv("SPIRE_MONGO_PORT"))
	if err != nil {
		panic(err)
	}
	s.MongoPort = port

	s.MongoDatabase = os.Getenv("SPIRE_MONGO_DATABASE")

	s.MongoUser = os.Getenv("SPIRE_MONGO_USER")

	data, err := os.ReadFile(os.Getenv("SPIRE_MONGO_PASSWORD_FILE"))
	if err != nil {
		panic(err)
	}
	s.MongoPassword = strings.TrimSpace(string(data))

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
	if (s.MongoHost == "") || (s.MongoDatabase == "") || (s.MongoUser == "") || (s.MongoPassword == "") {
		return false
	}
	if s.AuthKey == "" {
		return false
	}

	return true
}
