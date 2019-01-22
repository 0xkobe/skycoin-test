package main

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"sort"
	"strconv"
)

type Nodeinfo struct {
	Key				string		`json:"key"`
	Type			string		`json:"type"`
	SendBytes		int			`json:"send_bytes"`
	RecvBytes		int			`json:"recv_bytes"`
	LastAckTime		int			`json:"last_ack_time"`
	StartTime		int			`json:"start_time"`
}

func returnResponse(status bool, err error, data []byte, w http.ResponseWriter) {

}

func skywireNodes(w http.ResponseWriter, r *http.Request) {
	count, ok := r.URL.Query()["n"]
	if !ok || len(count[0]) < 1 {
		http.Error(w, "Url Param 'n' is missing", http.StatusBadRequest)
		return
	}

	n, err := strconv.Atoi(count[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	n++

	response, err := http.Get("http://discovery.skycoin.net:8001/conn/getAll")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body);
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data []Nodeinfo
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].StartTime > data[j].StartTime
	})

	//Marshal or convert user object back to json and write to response 
	dataJson, err := json.Marshal(data[1:n])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response 
	w.Write(dataJson)
}

func handleRequests() {
	http.HandleFunc("/api/nodesbytime", skywireNodes)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main() {
	handleRequests()
}
