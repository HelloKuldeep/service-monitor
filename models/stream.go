package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Stream struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	HitTime        int        `bson:"hittime" json:"hittime"`
	ThroughPut  int        `bson:"throughput" json:"throughput"`
}
