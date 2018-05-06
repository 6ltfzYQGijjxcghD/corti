package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := make(chan []byte, 300)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlePost(c, w, r)
		} else if r.Method == "GET" {
			handleGet(c, w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGet(queue chan []byte, w http.ResponseWriter, r *http.Request) {
	select {
	case b := <-queue:
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func handlePost(queue chan []byte, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	queue <- b
	w.WriteHeader(http.StatusNoContent)
}
