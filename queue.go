package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	MQUrl         = "https://message-queue.api.cloud.yandex.net"
	SigningRegion = "ru-central1"
)

func connectQueue(ctx context.Context) (client *sqs.Client, queueURL string, err error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           MQUrl,
				SigningRegion: SigningRegion,
			}, nil
		},
	)

	mqCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		err = fmt.Errorf("load config: %w", err)
		return
	}

	client = sqs.NewFromConfig(mqCfg)

	queue, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &cfg.QueueName,
	})
	if err != nil {
		err = fmt.Errorf("get queue url: %w", err)
		return
	}
	queueURL = *queue.QueueUrl
	return
}

// sendMesssage to mq
func sendMessage(ctx context.Context, msg message) (messageID string, err error) {
	client, queueURL, err := connectQueue(ctx)
	if err != nil {
		err = fmt.Errorf("connect to mq: %w", err)
		return
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("send message to queue ")
		return
	}

	msgBody := string(msgData)
	send, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &msgBody,
	})
	if err != nil {
		err = fmt.Errorf("send message to queue url: %s - %w", queueURL, err)
		return
	}
	messageID = *send.MessageId
	return
}
