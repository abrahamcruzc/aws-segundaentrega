package aws

import (
    "context"
    "fmt"
    "log"
)

type SNSMock struct{}

func NewSNSMock() *SNSMock {
    return &SNSMock{}
}

func (s *SNSMock) SendEmail(ctx context.Context, email string, subject string, message string) error {
    log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    log.Println("📧 SNS MOCK - Email enviado:")
    log.Printf("   Para: %s", email)
    log.Printf("   Asunto: %s", subject)
    log.Printf("   Mensaje: %s", message)
    log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

    fmt.Printf("\n[SNS MOCK] Email a %s: %s\n", email, subject)
    return nil
}
