package storage

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/codegangsta/cli"
)

// GoogleStore is a struct for handling storing on Google hardware
type GoogleStore struct {
	PubsubClient *pubsub.Client
	TopicID      string
}

// Name of this storage engine
func (s *GoogleStore) Name() string {
	return "GooglePubSub"
}

// IsConfigured tells us from the CLI flags passed if we're using this
// or not
func (s *GoogleStore) CanConfigure(c *cli.Context) bool {
	return c.String("googleProjectID") != ""
}

// AddCLIFlags extra flags for the watch command
func (s *GoogleStore) AddCLIFlags(fs []cli.Flag) []cli.Flag {
	return append(fs,
		cli.StringFlag{
			Name: "googleProjectID",
		},
		cli.StringFlag{
			Name:  "googleTopicID",
			Value: "bitping-event",
		},
	)
}

// Configure google storage
func (s *GoogleStore) Configure(c *cli.Context) error {
	projectID := c.String("googleProjectID")
	topicID := c.String("googleTopicID")
	// Configure storage
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	if exists, err := client.Topic(topicID).Exists(ctx); err != nil {
		return err
	} else if !exists {
		if _, err := client.CreateTopic(ctx, topicID); err != nil {
			return err
		}
	}

	s.TopicID = topicID
	s.PubsubClient = client

	return nil
}

// Push data to the storage engine
func (s *GoogleStore) Push(data interface{}) bool {
	if s.PubsubClient == nil {
		return false
	}
	ctx := context.Background()

	topic := s.PubsubClient.Topic(s.TopicID)
	exists, err := topic.Exists(ctx)

	if err != nil {
		log.Fatalf("Error checking for topic: %v", err)
		return false
	}
	if !exists {
		if _, err := s.PubsubClient.CreateTopic(ctx, s.TopicID); err != nil {
			log.Fatalf("Failed to create topic: %v", err)
			return false
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return false
	}

	_, err = topic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	return err == nil
}
