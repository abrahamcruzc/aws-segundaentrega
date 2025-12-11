package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Notifier interface para envío de notificaciones
type Notifier interface {
	Publish(ctx context.Context, subject, message string) error
}

// SNSClient implementación real de SNS
type SNSClient struct {
	client   *sns.Client
	topicARN string
}

func NewSNSClient(client *sns.Client, topicARN string) *SNSClient {
	return &SNSClient{
		client:   client,
		topicARN: topicARN,
	}
}

func (s *SNSClient) Publish(ctx context.Context, subject, message string) error {
	_, err := s.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Subject:  aws.String(subject),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("error al publicar en SNS: %w", err)
	}
	log.Printf("SNS: Mensaje - Asunto: %s", subject)
	return nil
}

// SNSMock implementación mock de SNS
type SNSMock struct{}

func NewSNSMock() *SNSMock {
	return &SNSMock{}
}

func (s *SNSMock) Publish(ctx context.Context, subject, message string) error {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Mensaje:")
	log.Printf("   Asunto: %s", subject)
	log.Printf("   Mensaje: %s", message)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	return nil
}
