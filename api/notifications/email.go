package notifications

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)


func SendJobAlert(title, company, url, location string) {
	sender := os.Getenv("EMAIL_FROM") 
	recipient := os.Getenv("EMAIL_TO")

	if sender == "" || recipient == "" {
		log.Println("[Email] Skipped: EMAIL_FROM or EMAIL_TO not set in .env")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("[Email] Error loading AWS config: %v", err)
		return
	}

	client := sesv2.NewFromConfig(cfg)


	subject := fmt.Sprintf("New Job: %s at %s", title, company)
	body := fmt.Sprintf("Role: %s\nCompany: %s\nLocation: %s\n\nApply here: %s", title, company, location, url)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(sender),
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body:    &types.Body{Text: &types.Content{Data: aws.String(body)}},
			},
		},
	}

	_, err = client.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("[Email] Failed to send: %v", err)
	} else {
		log.Printf("[Email] Sent alert for: %s", title)
	}
}