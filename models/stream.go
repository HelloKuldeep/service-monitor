package models

import "gopkg.in/mgo.v2/bson"

// Represents a stream
type Stream struct {
	ID			bson.ObjectId `bson:"_id" json:"id"`
	HitTime		int        `bson:"hittime" json:"hittime"`
	ResponseTime	int        `bson:"responsetime" json:"responsetime"`
}

// Represents a Incoming Data
type Input struct {
	StartTime	int        `bson:"starttime" json:"starttime"`
	EndTime		int        `bson:"endtime" json:"endtime"`
}

// Represents a Response
type Output struct {
	Time	int        `bson:"starttime" json:"starttime"`
	// Window		int        `bson:"window" json:"window"`
	MinResponseTime		int        `bson:"minresponsetime" json:"minresponsetime"`
	MaxResponseTime		int        `bson:"maxresponsetime" json:"maxresponsetime"`
	AverageThroughPut		int        `bson:"averagethroughput" json:"averagethroughput"`
}