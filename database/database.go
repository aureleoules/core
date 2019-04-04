package database

import (
	"github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

// DB : Database holder
var DB *mgo.Database

const (
	usersCollection              string = "users"
	emailVerificationsCollection string = "email_verifications"
	sitesCollection              string = "sites"
	contactCollection            string = "contact"
	aboutCollection              string = "about"
	articlesCollection           string = "articles"
	projectsCollection           string = "projects"
	galleriesCollection          string = "galleries"
	photosCollection             string = "photos"
	videoGroupsCollection        string = "videogroups"
	videosCollection             string = "videos"
	filesCollection              string = "files"
	albumsCollection             string = "albums"
	tracksCollection             string = "tracks"
)

//Connect : Connect to MongoDB
func Connect(server string, database string) {
	session, err := mgo.Dial(server)
	logrus.Infoln("Attempting to connect to " + database)
	if err != nil {
		logrus.Fatal(err)
	}
	DB = session.DB(database)
	logrus.Info("Successfully connected to " + database)
}
