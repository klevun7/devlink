package email

import (
	"context"
	"devlink/internal/models"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Service struct {
	client *ses.Client
	sender string
}

func NewService(region, accessKey, secretKey, senderEmail string) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: ses.NewFromConfig(cfg),
		sender: senderEmail,
	}, nil
}

func (s *Service) Broadcast(subscribers []models.Subscriber, newJobs []models.Job) {
	if s.sender == "" {
		log.Println("Email sender not configured, skipping broadcast")
		return
	}

	for _, sub := range subscribers {
		var relevantJobs []models.Job
		for _, j := range newJobs {
			if j.Type == sub.Preference {
				relevantJobs = append(relevantJobs, j)
			}
		}

		if len(relevantJobs) > 0 {
			s.sendToUser(sub.Email, relevantJobs)
		}
	}
}

func (s *Service) sendToUser(recipient string, jobs []models.Job) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h2>Found %d new %s jobs for you!</h2><ul>", len(jobs), jobs[0].Type))
	for _, j := range jobs {
		sb.WriteString(fmt.Sprintf("<li><strong>%s</strong> at %s <br/> <a href='%s'>Apply Here</a></li>", j.Title, j.Company, j.URL))
	}
	sb.WriteString("</ul>")

	input := &ses.SendEmailInput{
		Destination: &types.Destination{ToAddresses: []string{recipient}},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{Data: aws.String(sb.String())},
			},
			Subject: &types.Content{Data: aws.String(fmt.Sprintf("DevLink: %d New Jobs", len(jobs)))},
		},
		Source: aws.String(s.sender),
	}

	_, err := s.client.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
	} else {
		log.Printf("Sent email to %s", recipient)
	}
}