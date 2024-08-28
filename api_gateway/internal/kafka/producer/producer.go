package producer

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func Producer(key string, req []byte) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("broker:29092"),
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
	// topic := "unknown13"
	// record := kgo.Record{
	// 	Topic: topic,
	// 	Value: []byte("hello"),
	// }

	// CreateTopic(ctx, client, topic)
	//with key
	record := kgo.Record{
		Key:   []byte(key),
		Topic: "hoteluser17",
		Value: req,
	}
	if err := client.ProduceSync(ctx, &record).FirstErr(); err != nil {
		log.Println(err)
	}
	return nil
}
