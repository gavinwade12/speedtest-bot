package main

import (
		"fmt"
		"net/http"
		"io"
		"io/ioutil"
		"encoding/json"
		"strconv"
)

type Speed struct {
	Download string `json:'download'`
	Upload  string `json:'upload'`
}

func Speeds(w http.ResponseWriter, r *http.Request) {
	var speeds Speed
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
		fmt.Printf("too low " + speeds.Download)
	}

	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(speeds); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/speeds", Speeds)

	http.ListenAndServe(":8080", nil)
}