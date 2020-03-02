package azure

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/globalsign/mgo"

	log "github.com/sirupsen/logrus"

	"github.com/MichaelPalmer1/simple-api-go/providers/base"
)

// MongoDBAPI : Implementation of SimpleAPI that interacts with MongoDB backends
type MongoDBAPI struct {
	*base.SimpleAPI
	Filtering MongoDBFiltering
	Client    *mgo.Database
}

// Init : Initialize connection to MongoDB
func (api *MongoDBAPI) Init(address, database, username, password string) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{address},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	session.SetSafe(&mgo.Safe{})
	api.Client = session.DB(database)
}

// Close : Close connection with MongoDB
func (api *MongoDBAPI) Close() {
	api.Client.Session.Close()
}
