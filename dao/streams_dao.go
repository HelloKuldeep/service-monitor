package dao

import (
	"log"

	model "a-service-monitor/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StreamsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "streams"
)

// Establish a connection to database
func (m *StreamsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of streams
func (m *StreamsDAO) FindAll() ([]model.Stream, error) {
	var streams []model.Stream
	err := db.C(COLLECTION).Find(bson.M{}).All(&streams)
	return streams, err
}

// Insert a stream into database
func (m *StreamsDAO) Insert(stream model.Stream) error {
	err := db.C(COLLECTION).Insert(&stream)
	return err
}
