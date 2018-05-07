package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type msg interface {
	isMsg()
}

type end struct{}

func (m *end) isMsg() {}

type data struct {
	payload *[]byte
}

func (m *data) isMsg() {}

const bufferSize = 300

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	channels := make(map[string](chan msg))
	cMutex := sync.Mutex{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var path string
		path = strings.Trim(r.URL.Path, "/")
		if path == "" {
			path = randStringBytes(16)
		}

		switch r.Method {
		case "PUT":
			cMutex.Lock()
			defer cMutex.Unlock()
			if _, ok := channels[path]; ok {
				w.WriteHeader(http.StatusConflict)
				return
			}
			channels[path] = make(chan msg, 300)
			w.Write([]byte(path))
		case "POST":
			cMutex.Lock()
			q, ok := channels[path]
			cMutex.Unlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			handlePost(q, w, r)
		case "DELETE":
			cMutex.Lock()
			q, ok := channels[path]
			cMutex.Unlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			q <- &end{}
		case "GET":
			cMutex.Lock()
			q, ok := channels[path]
			cMutex.Unlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			shouldDelete := handleGet(q, w, r)
			if shouldDelete {
				cMutex.Lock()
				delete(channels, path)
				cMutex.Unlock()
			}
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGet(queue chan msg, w http.ResponseWriter, r *http.Request) bool {
	select {
	case b := <-queue:
		switch d := b.(type) {
		case *data:
			w.Write(*(*d).payload)
		case *end:
			w.WriteHeader(http.StatusGone)
			return true
		default:
			fmt.Printf("Forgotten type %T", d)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusNoContent)
	}
	return false
}

func handlePost(queue chan msg, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	queue <- &data{&b}
	w.WriteHeader(http.StatusNoContent)
}

// From https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
