package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestDynamoLocal(t *testing.T) {

	tr := true
	cfg := &aws.Config{
		Endpoint:                      aws.String("http://localhost:8000"),
		Region:                        aws.String("us-east-1"),
		CredentialsChainVerboseErrors: &tr,
	}
	d := newDynamodb(cfg)
	if err := d.dropTables(); err != nil {
		// TODO switch on error type
		t.Logf("Error dropping tables: %v", err)
	}
	if err := d.createTables(); err != nil {
		// TODO switch on error type
		t.Logf("Error creating tables: %v", err)
	}

	if err := d.putVideoList(&videoList{
		ID:     "a123",
		Videos: []*video{{ID: "123", Title: "A time to remember"}},
	}); err != nil {
		t.Error("Error putting video:", err)
	}
	videoLists, err := d.getVideoLists()
	if err != nil {
		t.Error("Error getting videos:", err)
	}
	if len(videoLists) != 1 {
		t.Errorf("should be one video - got %+v", len(videoLists))
	} else {
		t.Logf("Success: videoLists[0]: %+v", videoLists[0])
	}

	err = d.putVote(&vote{VideoID: "123", PersonID: "Me", Up: true})
	if err != nil {
		t.Error("Error putting vote:", err)
	}
	votes, err := d.getVotes()
	if err != nil {
		t.Error("Error getting votes:", err)
	}
	if len(votes) != 1 {
		t.Errorf("should be one video - got %+v", len(votes))
	} else {
		t.Logf("Success: votes[0]: %+v", votes[0])
	}
}
