package main

import "sort"

type user struct {
	ID         string   `json:"id" db:"id"`
	VideoLists []string `json:"videoLists"`
}

type videoList struct {
	ID        int      `json:"id" db:"id"`
	Title     string   `json:"title" db:"title"`
	Videos    []*video `json:"videos"`
	Votes     []*vote  `json:"votes"`
	CreatorID string   `json:"creatorId" db:"creator_id"`
}

type videoSource string

const (
	youtube videoSource = "youtube"
)

type video struct {
	ID          int         `json:"id" db:"id"`
	SourceID    int         `json:"sourceId" db:"source_id"`
	Source      videoSource `json:"source" db:"source"`
	VideoListID int         `json:"videoListId" db:"video_list_id"`
	Title       string      `json:"title" db:"title"`
	Votes       int         `json:"votes"`
	CreatorID   string      `json:"creatorId" db:"creator_id"`
}

type vote struct {
	ID          int    `json:"id" db:"id"`
	VideoID     int    `json:"videoId" db:"video_id"`
	VideoListID int    `json:"videoListId" db:"video_list_id"`
	PersonID    string `json:"personId,omitempty" db:"person_id"`
	PersonHash  string `json:"personHash"`
	Up          bool   `json:"up" db:"up"`
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
