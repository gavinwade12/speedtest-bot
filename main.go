package main

import (
		"fmt"
		"net/http"
		"net/url"
		"io"
		"io/ioutil"
		"encoding/json"
		"strconv"
)

type Speeds struct {
	Download string `json:'download'`
	Upload  string `json:'upload'`
}

func CheckSpeeds(w http.ResponseWriter, r *http.Request) {
	var speeds Speeds
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10485))
	if err != nil {
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

	i, err := strconv.ParseFloat(speeds.Download, 32)
	if err != nil {
		return
	}

	if (i < 10.0) {
		SendTweet(*speeds)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(speeds); err != nil {
		panic(err)
	}
}

func SendTweet(&speeds Speeds) {
	status := "Switched from BCS 50 Mbps $50/month to @TWC 15 Mbps $50/month to actually get Download: " + speeds.Downlad + " Mbps Upload: " + speeds.Upload + " Mbps. Not a happy camper."
	status = url.QueryEscape(status)
	client := &http.client
	req.Header.Add("Authorization", "oauth_consumer_key=\"" + consumerKey + "\"")
	req.Header.Add("Authorization", "oauth_nonce=\"" + nonce + "\"")
	req.Header.Add("Authorization", "oauth_signature=\"" + )
	req.Header.Add("Authorization", "")
	req.Header.Add("Authorization", "")
	req.Header.Add("Authorization", "")
	req.Header.Add("Authorization", "")
	req.Header.Add("Authorization", "")
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/update.json?status=" + status)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/speeds", CheckSpeeds)

	http.ListenAndServe(":8080", nil)
}