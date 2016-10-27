package main

import (
		"fmt"
		"bytes"
		"net/http"
		"io"
		"io/ioutil"
		"encoding/json"
		"strconv"
)

func CheckSpeeds(w http.ResponseWriter, r *http.Request, header Header) {
	var speeds Speeds
	if body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10485)); err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &speeds); err != nil {
		w.Header().Set("Content-Type", "applicatin/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	if i, err := strconv.ParseFloat(speeds.Download, 32); err != nil {
		panic(err)
	}

	if (i < 10.0) {
		SendTweet(speeds, header)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(speeds); err != nil {
		panic(err)
	}
}

func SendTweet(speeds Speeds, header Header) {
	status := "Switched from BCS 50 Mbps $50/month to @TWC 15 Mbps $50/month to actually get Download: " + speeds.Download + " Mbps Upload: " + speeds.Upload + " Mbps. Not a happy camper."
	status = QEscape(status)
	header = UpdateHeaderValues(header, status)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/update.json?include_entities=true", bytes.NewBufferString("status=" + status))
	if err != nil {
		panic(err)
	}
	headerString := GetCompleteHeaderString(header)
	fmt.Printf("%s\n\n", headerString)
	req.Header.Add("Authorization", headerString)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len("status=" + status)))
	if res, err := client.Do(req); err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(io.LimitReader(res.Body, 10485678))
	fmt.Printf("%s\n\n", string(body))
	res.Body.Close()
}

func main() {
	header := GetHeaderValues()
	http.HandleFunc("/speeds", func(w http.ResponseWriter, r *http.Request) {
		CheckSpeeds(w, r, *header)
	})

	http.ListenAndServe(":8080", nil)
}
