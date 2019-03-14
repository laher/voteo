package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	db = &dtb{}
)

type dtb struct {
	lock   sync.RWMutex
	videos []*video
	votes  []*vote
}

func (db *dtb) loadVideos() {
	db.lock.Lock()
	defer db.lock.Unlock()
	b, err := ioutil.ReadFile("videos.json")
	if err != nil {
		if os.IsNotExist(err) {
			db.videos = initVideos
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVideos := []*video{}
	err = json.Unmarshal(b, &myVideos)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	db.videos = myVideos
}

func (db *dtb) loadVotes() {
	db.lock.Lock()
	defer db.lock.Unlock()
	b, err := ioutil.ReadFile("votes.json")
	if err != nil {
		if os.IsNotExist(err) {
			db.votes = []*vote{}
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVotes := []*vote{}
	err = json.Unmarshal(b, &myVotes)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	db.votes = myVotes
}

func (db *dtb) getVideos() []*video {
	db.lock.RLock()
	defer db.lock.RUnlock()
	newVideos := make([]*video, 0, len(db.videos))
	for _, v := range db.videos {
		vi := *v
		newVideos = append(newVideos, &vi)
	}
	return newVideos
}

func (db *dtb) writeVideos(myVideos []*video) {
	db.lock.Lock()
	defer db.lock.Unlock()
	b, _ := json.MarshalIndent(myVideos, "", "  ")
	err := ioutil.WriteFile("videos.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
	db.videos = myVideos
}

func (db *dtb) getVotes() []*vote {
	db.lock.RLock()
	defer db.lock.RUnlock()
	newVotes := make([]*vote, 0, len(db.votes))
	for _, v := range db.votes { // defensive copy
		vo := *v
		newVotes = append(newVotes, &vo)
	}
	return newVotes
}

func (db *dtb) writeVotes(myVotes []*vote) {
	db.lock.Lock()
	defer db.lock.Unlock()
	b, _ := json.MarshalIndent(myVotes, "", "  ")
	err := ioutil.WriteFile("votes.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
	db.votes = myVotes
}
