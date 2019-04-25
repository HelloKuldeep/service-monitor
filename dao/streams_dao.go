package dao

import (
	"math"
	"math/rand"
	"fmt"
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

// Find list of output
func (m *StreamsDAO) FindOutput(input model.Input) ([]model.Output, error) {
	var err error
	var streams model.Stream
	outputs := []model.Output{}
	slider := SLIDE
	// fmt.Println(input.StartTime, input.EndTime)
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
	}
	return outputs, err
}

func (m *StreamsDAO) FindOutput1Min(input model.Input) ([]model.Output, error) {
	var err error
	var slider int
	var streams model.Stream
	outputs := []model.Output{}
	newStartTime := int(math.Round(float64(input.StartTime)/60)*60)
	timeDiff := input.EndTime - input.StartTime
	newEndTime := int(math.Round(float64(input.EndTime)/60)*60)
	newTimeDiff := newEndTime - newStartTime
	totalError := (math.Abs(float64(newStartTime) - float64(input.StartTime)) + math.Abs(float64(newStartTime) - float64(input.StartTime)))/float64(timeDiff)
	totalError = math.Round(totalError*10000)/10000
	if timeDiff < 6000 {
		slider = SLIDE
	} else {
		sliderTemp := rand.Intn(100 - 15) + 15
		slider = (newTimeDiff/sliderTemp)/60
	}
	// fmt.Println(input.StartTime, newStartTime, input.EndTime, newEndTime, totalError)
	fmt.Println("Each Window Min: ", slider)
	fmt.Println("Error: ", totalError)
	// var chanNumber int = (input.EndTime - input.StartTime)/(slider*60)
    done := make(chan bool)
	outputChan := make(chan model.Output)
	for i := newStartTime; i < newEndTime; i = i + (slider * 60) {
		go func() {
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
				return
			}
			var output model.Output
			if count > 0 {
				output = model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageResponseTime: total/count}
				// outputs = append(outputs, output)
			} else {
				output = model.Output{Time: i, MinResponseTime: 0, MaxResponseTime: 0, AverageResponseTime: 0}
				// outputs = append(outputs, output)
			}
			// fmt.Println(output)
			outputChan <- output
			done <- true
		}()

		// go ForEachSlide(i, slider, outputChan)
		// for outputC := range outputChan {
		// 	fmt.Println("First", outputC)
		// 	outputs = append(outputs, outputC)
		// }
	}
	go func() {
		for outputC := range outputChan {
			// fmt.Println("First", outputC)
			outputs = append(outputs, outputC)
		}
	}()
	for i := newStartTime; i < newEndTime; i = i + (slider * 60) {
		<- done
	}
	return outputs, err
}

// Method for requesting handeling using goroutine and channels
func ForEachSlide(i int, slider int, outputChan chan model.Output) {
	var streams model.Stream
	var err error
	iter := db.C(COLLECTION).Find(bson.M{"hittime": bson.M{"$gte": i, "$lt": i + (slider * 60)}}).Iter()
	var minVal, maxVal, total, count int = MAXINT, MININT, 0, 0
	for iter.Next(&streams) {
		if maxVal < streams.ResponseTime {
			maxVal = streams.ResponseTime
		}
		if minVal > streams.ResponseTime {
			minVal = streams.ResponseTime
		}
		total += streams.ResponseTime
		count++
	}
	if err = iter.Close(); err != nil {
		return
	}
	var output model.Output
	if count > 0 {
		output = model.Output{Time: i, MinResponseTime: minVal, MaxResponseTime: maxVal, AverageResponseTime: total / count}
		// outputs = append(outputs, output)
	} else {
		output = model.Output{Time: i, MinResponseTime: 0, MaxResponseTime: 0, AverageResponseTime: 0}
		// outputs = append(outputs, output)
	}
	fmt.Println("Inside ", output)
	outputChan <- output
	// close(outputChan)
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