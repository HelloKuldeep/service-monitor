package dao

import (
	"math"
	_"fmt"
	"log"
	_"os"
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
	SLIDE      = 1
	MAXINT     = 2147483647
	MININT     = -2147483647
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
	var streams model.Stream
	slider := SLIDE
	outputs := []model.Output{}
	var err error
	// var chanNumber int = (input.EndTime - input.StartTime) / 60
	// // quit := make(chan int)
	// outputChan := make(chan model.Output, chanNumber)
	for i := input.StartTime; i < input.EndTime; i = i + (slider * 60) {
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
				break
			}
			if count > 0 {
				output := model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageResponseTime: total/count}
				outputs = append(outputs, output)
			}
		// go ForEachSlide(i, outputChan)
	}
	// for {
	// 	outputC, more := <- outputChan
	// 	if more {
	// 		outputs = append(outputs, outputC)
	// 	} else {
	// 		break
	// 	}
	// }
	// for outputC := range outputChan {
	// 	outputs = append(outputs, outputC)
	// }
	return outputs, err
}

// Insert a stream into database
func (m *StreamsDAO) Insert(stream model.Stream) error {
	err := db.C(COLLECTION).Insert(&stream)
	return err
}

func (m *StreamsDAO) Insert1Min(streams []model.Stream) error {
	var err error
	for stream := range streams {
		err = db.C(COLLECTION).Insert(&stream)
		if err != nil {
			return err
		}
	}
	return err
}

func (m *StreamsDAO) FindOutput1Min(input model.Input) ([]model.Output, error) {
	var err error
	timeDiff := input.EndTime - input.StartTime
	newStartTime := math.Round(float64(input.StartTime)/60)
	newEndTime := math.Round(float64(input.EndTime)/60)
	if timeDiff < 6000 {
		var outputs []model.Output
		err = db.C(COLLECTION).Find(bson.M{}).All(&outputs)
	} else {
		// var := rand.Intn(max - min) + min
	}
	
	var streams model.Stream
	slider := SLIDE
	outputs := []model.Output{}
	
	for i := input.StartTime; i < input.EndTime; i = i + (slider * 60) {
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
				break
			}
			if count > 0 {
				output := model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageResponseTime: total/count}
				outputs = append(outputs, output)
			}
		// go ForEachSlide(i, outputChan)
	}
	return outputs, err
}

// func writefile() {  
//     f, err := os.Create("test.txt")
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     l, err := f.WriteString("Hello World")
//     if err != nil {
//         fmt.Println(err)
//         f.Close()
//         return
//     }
//     fmt.Println(l, "bytes written successfully")
//     err = f.Close()
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
// }


// // Method for requesting handeling using goroutine and channels
// func ForEachSlide(i int, outputChan chan model.Output) {
// 	var streams model.Stream
// 	slider := SLIDE
// 	var err error
// 	// outputs := []model.Output{}
// 	iter := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$gte": i, "$lt": i + (slider * 60)}}).Iter()
// 	var minVal, maxVal, total, count int = MAXINT, MININT, 0, 0
// 	for iter.Next(&streams) {
// 		if maxVal < streams.ResponseTime {
// 			maxVal = streams.ResponseTime
// 		}
// 		if minVal > streams.ResponseTime {
// 			minVal = streams.ResponseTime
// 		}
// 		total += streams.ResponseTime
// 		count++
// 	}
// 	if err = iter.Close(); err != nil {
// 		return
// 	}
// 	var output model.Output
// 	if count > 0 {
// 		output = model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageResponseTime: total / count}
// 		// outputs = append(outputs, output)
// 	}
// 	outputChan <- output
// 	// close(outputChan)
// }