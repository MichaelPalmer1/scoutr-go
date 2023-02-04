package mongo

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/globalsign/mgo"

	log "github.com/sirupsen/logrus"

	"github.com/MichaelPalmer1/scoutr-go/providers/base"
)

// MongoAPI : Implementation of Scoutr that interacts with MongoDB backends
type MongoAPI struct {
	*base.Scoutr
	Filtering MongoDBFiltering
	Client    *mgo.Database
}

// Init : Initialize connection to MongoDB
func (api *MongoAPI) Init(address, database, username, password string) {
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
func (api *MongoAPI) Close() {
	api.Client.Session.Close()
}
