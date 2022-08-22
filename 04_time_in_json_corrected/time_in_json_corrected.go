package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const layoutDT = "2006-01-02 15:04"

type Response struct {
	Time string `json:"time"`
}

func printOut(t string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, t)
	}
}

func jsonOut(ttype string, t string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp[ttype] = t
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.Write([]byte("Error happened in JSON marshal"))
		} else {
			w.Write(jsonResp)
		}

	}
}

func main() {

	t := time.Now()
	//fmt.Println(t.Format("2006-01-02 15:04"))

	http.HandleFunc("/time", printOut(t.Format("2006-01-02 15:04")))

	http.HandleFunc("/json", jsonOut("time", t.Format("2006-01-02 15:04")))

	http.ListenAndServe(":8090", nil)

}
