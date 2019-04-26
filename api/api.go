package api

import (
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	cli "service-monitor/redis"
	"io/ioutil"
	"net/http"
	"strconv"
)

var conn *cli.PoolConn

func Start(port int) {
	fmt.Println("Start")//
	conn = cli.New()
	router := httprouter.New()
	router.POST("/newstream", Add)
	router.GET("/newstream", List)
	http.ListenAndServe(":"+strconv.Itoa(port), router)
}

func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// fmt.Println("Into Add Request")//

	b, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	receiver := cli.Stream{}

	if err = json.Unmarshal(b, &receiver); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// fmt.Println(receiver)//

	if err = conn.Add(receiver.HitTime, receiver.ResponseTime); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// fmt.Println("Added")//

	w.WriteHeader(201)
}

func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	receiver := cli.Stream{}

	if err = json.Unmarshal(b, &receiver); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// fmt.Println(receiver)//
	receivers, err := conn.Get(receiver.HitTime)
	if err != nil {
		w.Write([]byte("No content"))
		return
	}
	fmt.Println("Result1:", receivers)//
	out, _ := json.Marshal(receivers)
	fmt.Println("Result2:", out)//
	w.WriteHeader(200)
	w.Write(out)
}