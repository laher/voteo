package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	videoLists = "videoLists"
	users      = "users"
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
	for _, t := range []string{"videos", "votes", videoLists, users} {
		input := &dynamodb.DeleteTableInput{
			TableName: aws.String(t),
		}
		_, err := db.d.DeleteTable(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *ddb) createTables() error {

	inputs := []*dynamodb.CreateTableInput{
		&dynamodb.CreateTableInput{
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
		},
		&dynamodb.CreateTableInput{
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
		},
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(videoLists),
		},
		&dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(users),
		},
	}

	for _, input := range inputs {
		_, err := db.d.CreateTable(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *ddb) putUser(u *user) error {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(users),
	}
	_, err = db.d.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}
func (db *ddb) addVideoToList(v *videoList, video *video) error {
	av, err := dynamodbattribute.MarshalMap(video)
	if err != nil {
		return err
	}
	name := expression.Name("videos")
	value := expression.Value(av)
	update := expression.UpdateBuilder{}
	update = update.Add(name, value)
	builder := expression.NewBuilder().WithUpdate(update)
	expression, err := builder.Build()
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		Key:                       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(v.ID)}},
		TableName:                 aws.String(videoLists),
		UpdateExpression:          expression.Update(),
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err = db.d.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (db *ddb) putVideoList(v *videoList) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(videoLists),
	}
	_, err = db.d.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (db *ddb) getVideoList(id string) (*videoList, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(videoLists),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
	}

	result, err := db.d.GetItem(input)
	if err != nil {
		return nil, err
	}
	if len(result.Item) == 0 {
		return nil, errors.New("not found")
	}
	log.Printf("result: %+v", result.Item)
	item := &videoList{}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}
	log.Printf("video list: %+v", item)
	return item, nil
}

func (db *ddb) getUser(id string) (*user, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(users),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
	}

	result, err := db.d.GetItem(input)
	if err != nil {
		return nil, err
	}
	item := &user{}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (db *ddb) getVideoLists() ([]*videoList, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(videoLists),
	}

	result, err := db.d.Scan(input)
	if err != nil {
		return nil, err
	}

	vls := make([]*videoList, 0, len(result.Items))
	for _, i := range result.Items {
		item := &videoList{}
		err = dynamodbattribute.UnmarshalMap(i, item)
		if err != nil {
			return vls, err
		}
		vls = append(vls, item)
	}
	log.Printf("found %d video lists", len(vls))
	return vls, nil
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
	log.Printf("found %d videos", len(videos))
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
