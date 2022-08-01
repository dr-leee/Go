package main

import (
	"fmt"
	"net/http"
	"time"
)

func time_out(t string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, t)
	}
}

func main() {

	t := time.Now()
	//fmt.Println(t.Format("2006-01-02 15:04"))

	http.HandleFunc("/time", time_out(t.Format("2006-01-02 15:04")))

	http.ListenAndServe(":8090", nil)

}
