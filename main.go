package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/", fs)
	http.HandleFunc("/videos", crudHandler)
	http.HandleFunc("/yt/data", ytDataHandler)
	loadVideos()
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

type video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Votes int    `json:"votes"`
}

var (
	videos     = []video{}
	initVideos = []video{
		{ID: "JntjzuI5rGM", Title: "Dave Grohl Tells ...", Votes: 0},
		{ID: "X7hFERntlog", Title: "Fearless Organization", Votes: 0},
		{ID: "d_HHnEROy_w", Title: "Stop managing, start ...", Votes: -1},
		{ID: "BCkCvay4-DQ", Title: "Push Kick", Votes: 1},
	}
	lock = sync.RWMutex{}
)

func loadVideos() {
	lock.Lock()
	defer lock.Unlock()
	b, err := ioutil.ReadFile("videos.json")
	if err != nil {
		if os.IsNotExist(err) {
			videos = initVideos
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVideos := []video{}
	err = json.Unmarshal(b, &myVideos)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	videos = myVideos
}

func writeVideos() {
	lock.Lock()
	defer lock.Unlock()
	b, _ := json.MarshalIndent(videos, "", "  ")
	err := ioutil.WriteFile("videos.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
}

func crudHandler(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Query().Get("id")
	switch r.Method {
	case http.MethodGet:
		// just get all
	case http.MethodPut:
		myVideos := []video{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&myVideos)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldnt decode: %v", err)
			return
		}
		lock.Lock()
		videos = myVideos
		lock.Unlock()
		writeVideos()
		log.Printf("Updated: %+v", myVideos)
	}
	lock.RLock()
	defer lock.RUnlock()
	b, _ := json.Marshal(videos)
	w.Write(b)
	log.Printf("Returned: %+v", videos)
}

func ytDataHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// no video id
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("No video ID")
		return
	}
	url := "https://www.youtube.com/oembed?url=http%3A//youtube.com/watch%3Fv%3D" + id
	resp, err := http.Get(url)
	if err != nil {
		// could not connect
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	l := resp.Header.Get("Content-Length")
	if l != "" {
		w.Header().Set("Content-Length", l)
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("Retrieved video metadata for: %+v", id)
	} else {
		log.Printf("Error retrieving video metadata for %+v: %v", id, resp.Status)
	}
}
