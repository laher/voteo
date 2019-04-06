package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ddb struct {
	d *dynamodb.DynamoDB
}

func newDynamodb(cfg *aws.Config) *ddb {
	sess := session.Must(session.NewSession(cfg))
	db := dynamodb.New(sess)
	return &ddb{db}
}

func (db *ddb) dropTables() error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("videos"),
	}
	_, err := db.d.DeleteTable(input)
	if err != nil {
		return err
	}
	inputV := &dynamodb.DeleteTableInput{
		TableName: aws.String("votes"),
	}
	_, err = db.d.DeleteTable(inputV)
	if err != nil {
		return err
	}
	return nil
}

func (db *ddb) createTables() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("videos"),
	}
	_, err := db.d.CreateTable(input)
	if err != nil {
		return err
	}

	inputVotes := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("videoId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("personId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("videoId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("personId"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("votes"),
	}

	_, err = db.d.CreateTable(inputVotes)
	if err != nil {
		return err
	}
	return nil
}

func (db *ddb) getVideos() ([]*video, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String("videos"),
	}

	result, err := db.d.Scan(input)
	if err != nil {
		return nil, err
	}

	videos := make([]*video, 0, len(result.Items))
	for _, i := range result.Items {
		item := &video{}
		err = dynamodbattribute.UnmarshalMap(i, item)
		if err != nil {
			return videos, err
		}
		videos = append(videos, item)
	}
	return videos, nil
}

func (db *ddb) put(v *video) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("videos"),
	}
	_, err = db.d.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (db *ddb) getVotes() ([]*vote, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("votes"),
	}

	result, err := db.d.Scan(input)
	if err != nil {
		return nil, err
	}

	votes := make([]*vote, 0, len(result.Items))
	for _, i := range result.Items {
		item := &vote{}
		err = dynamodbattribute.UnmarshalMap(i, item)
		if err != nil {
			return votes, err
		}
		votes = append(votes, item)
	}
	return votes, nil
}

func (db *ddb) putVote(v *vote) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("votes"),
	}
	_, err = db.d.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}
