package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hashicorp/vault/api"

	"database/sql"

	_ "github.com/lib/pq"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	fmt.Println("Started processing event, secrets injected")

	// Check for required env vars env vars
	secretFile := os.Getenv("VAULT_SECRET_FILE_DB")
	if secretFile == "" {
		return errors.New("no VAULT_SECRET_FILE_DB environment variable, exiting")
	}

	dbURL := os.Getenv("DATABASE_ADDR")
	if dbURL == "" {
		return errors.New("no DATABASE_ADDR environment variable, exiting")
	}

	// Read the secret from the file before processing the event
	secretRaw, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Decode the JSON into a map[string]interface{}
	var secret api.Secret
	b := bytes.NewBuffer(secretRaw)
	dec := json.NewDecoder(b)
	dec.UseNumber()

	if err := dec.Decode(&secret); err != nil {
		return err
	}

	// Connect to the database and insert the registration
	connStr := fmt.Sprintf("postgres://%s:%s@%s", secret.Data["username"], secret.Data["password"], dbURL)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to the database")

	defer db.Close()

	// We don't need to do anything with the DB connection
	// The exercise is about securing the connection with Vault
	// move on to determining points at this stage

	for _, message := range sqsEvent.Records {
		// Unmarshal the JSON into a GameEvent struct
		var event GameEvent
		if err := json.Unmarshal([]byte(message.Body), &event); err != nil {
			return err
		}

		// if the participants have made it this far, they get points for every message
		lc, _ := lambdacontext.FromContext(ctx)
		arn := lc.InvokedFunctionArn
		var lbEvent LeaderboardEvent = LeaderboardEvent{
			FunctionARN:  arn,
			FunctionName: lambdacontext.FunctionName,
			AccountID:    strings.Split(arn, ":")[4],
			Points:       100, // TODO move to central record function
		}

		// Marshal the LeaderboardEvent struct into JSON
		lbEventJSON, err := json.Marshal(lbEvent)
		if err != nil {
			return err
		}

		// Publish the LeaderboardEvent to the SNS topic
		mySession := session.Must(session.NewSession())
		svc := sqs.New(mySession)
		body := string(lbEventJSON)
		_, err = svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(0),
			MessageBody:  &body,
			QueueUrl:     &event.LeaderboardQueue,
		})
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
