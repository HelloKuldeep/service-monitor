package dao

import (
	"log"

	model "service-monitor/models"
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
	SLIDE = 1
	MAXINT = 2147483647
	MININT = -2147483647
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

// Find list of output
func (m *StreamsDAO) FindOutput(input model.Input) ([]model.Output, error) {
	//var streams model.Stream //Use

	// err := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$gte": input.StartTime, "$lte": input.EndTime}}).All(&streams)
	
	slider := SLIDE	//Use
	
	// window := input.EndTime - input.StartTime + 1
	// var output [window]model.Output
	// k := 0
	// for i := input.StartTime; i <= input.EndTime; i = i+slider {
	// count, err := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": input.StartTime}}).Count()
	// 	minVal := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": input.StartTime}}).Sort("responsetime").Limit(1)
	// 	output[k].AverageThroughPut = count
	// }
	// err := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": input.StartTime}}).Sort("responsetime").Select(bson.M{"responsetime": 1}).One(&val)
	// err := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": input.StartTime}}).Sort("-responsetime").Limit(1).All(&streams)
	
	outputs := []model.Output{}
	
	// for i := input.StartTime; i <= input.EndTime ;i = i+slider {
	// 	var minVal model.Stream
	// 	db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": i}}).Sort("responsetime").Select(bson.M{"responsetime": 1}).One(&minVal)
	// 	var maxVal model.Stream
	// 	db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$eq": i}}).Sort("-responsetime").Select(bson.M{"responsetime": 1}).One(&maxVal)


	// 	output := model.Output{Time: i, MinResponseTime: minVal.ResponseTime, MaxResponseTime: maxVal.ResponseTime}
	// 	outputs = append(outputs, output)
	// }
	var err error
	var chanNumber int = (input.EndTime-input.StartTime)/60
	outputC : make(chan []model.Output, )
	for i := input.StartTime; i < input.EndTime ;i = i+(slider*60) {
		// iter := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$gte": i, "$lt": i+(slider*60)}}).Iter()
		// var minVal, maxVal, total, count int = MAXINT, MININT, 0, 0
		// for iter.Next(&streams) {
		// 	if maxVal < streams.ResponseTime {
		// 		maxVal = streams.ResponseTime
		// 	}
		// 	if minVal > streams.ResponseTime {
		// 		minVal = streams.ResponseTime
		// 	}
		// 	total += streams.ResponseTime
		// 	count ++;
		// }
		// if err = iter.Close(); err != nil {
		// 	break
		// }
		// if count > 0 {
		// 	output := model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageThroughPut: total/count}
		// 	outputs = append(outputs, output)
		// }
		outputs = go ForEachSlide(i)
	}
	return outputs, err
}

func ForEachSlide(i int) ([]model.Output, error) {
	var streams model.Stream
	slider := SLIDE
	var err error
	outputs := []model.Output{}
	iter := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$gte": i, "$lt": i+(slider*60)}}).Iter()
	var minVal, maxVal, total, count int = MAXINT, MININT, 0, 0
	for iter.Next(&streams) {
		if maxVal < streams.ResponseTime {
			maxVal = streams.ResponseTime
		}
		if minVal > streams.ResponseTime {
			minVal = streams.ResponseTime
		}
		total += streams.ResponseTime
		count ++;
	}
	if err = iter.Close(); err != nil {
		return outputs, err
	}
	if count > 0 {
		output := model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageThroughPut: total/count}
		outputs = append(outputs, output)
	}
	return outputs, err
}

// Insert a stream into database
func (m *StreamsDAO) Insert(stream model.Stream) error {
	err := db.C(COLLECTION).Insert(&stream)
	return err
}
