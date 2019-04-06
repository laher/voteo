package main

import "sort"

type videoList struct {
	ID        string   `json:"id" dynamodbav:"id"`
	Videos    []*video `json:"videos" dynamodbav:"videos"`
	Votes     []*vote  `json:"votes" dynamodbav:"votes"`
	CreatorID string   `json:"creatorId" dynamodbav:"creatorId"`
}

type video struct {
	ID        string `json:"id" dynamodbav:"id"`
	Title     string `json:"title" dynamodbav:"title"`
	Votes     int    `json:"votes"`
	CreatorID string
}

type vote struct {
	VideoID    string `json:"videoId" dynamodbav:"videoId"`
	PersonID   string `json:"personId,omitempty" dynamodbav:"personId"`
	PersonHash string `json:"personHash"`
	Up         bool   `json:"up" dynamodbav:"up"`
}

func sortByVotes(videos []*video, votes []*vote) {
	// reset votes ...
	for _, video := range videos {
		video.Votes = 0
	}
	for _, vote := range votes {
		for _, video := range videos {
			if video.ID == vote.VideoID {
				if vote.Up {
					video.Votes++
				} else {

					video.Votes--
				}
			}
		}
	}
	sort.Sort(sort.Reverse(videosByVotes(videos)))
}

type videosByVotes []*video

func (a videosByVotes) Len() int           { return len(a) }
func (a videosByVotes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a videosByVotes) Less(i, j int) bool { return a[i].Votes < a[j].Votes }
