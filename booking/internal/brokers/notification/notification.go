package notification

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func Producer(req []byte) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Println("client not connected to kafka", err)
	}
	topic := "notification"
	record := kgo.Record{
		Topic: topic,
		Value: req,
	}

	// CreateTopic(ctx, client, topic)
	//with key
	// record := kgo.Record{
	// 	Key:   []byte(key),
	// 	Topic: topic,
	// 	Value: req,
	// }
	if err := client.ProduceSync(ctx, &record).FirstErr(); err != nil {
		log.Println(err)
	}
	return nil
}
