package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	lock   = sync.RWMutex{}
	videos = []*video{}
	votes  = []*vote{}
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
	myVideos := []*video{}
	err = json.Unmarshal(b, &myVideos)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	videos = myVideos
}

func loadVotes() {
	lock.Lock()
	defer lock.Unlock()
	b, err := ioutil.ReadFile("votes.json")
	if err != nil {
		if os.IsNotExist(err) {
			videos = initVideos
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVotes := []*vote{}
	err = json.Unmarshal(b, &myVotes)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	votes = myVotes
}

func getVideos() []*video {
	lock.RLock()
	defer lock.RUnlock()
	newVideos := make([]*video, 0, len(videos))
	for _, v := range videos {
		newVideos = append(newVideos, v)
	}
	return newVideos
}

func writeVideos(myVideos []*video) {
	lock.Lock()
	defer lock.Unlock()
	b, _ := json.MarshalIndent(videos, "", "  ")
	err := ioutil.WriteFile("videos.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
	videos = myVideos
}

func getVotes() []*vote {
	lock.RLock()
	defer lock.RUnlock()
	newVotes := make([]*vote, 0, len(votes))
	for _, v := range votes {
		newVotes = append(newVotes, v)
	}
	return newVotes
}

func writeVotes(myVotes []*vote) {
	lock.Lock()
	defer lock.Unlock()
	b, _ := json.MarshalIndent(votes, "", "  ")
	err := ioutil.WriteFile("votes.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
	votes = myVotes
}
