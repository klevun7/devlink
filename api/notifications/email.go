package notifications

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type EmailJob struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	URL      string `json:"url"`
	Location string `json:"location"`
}

func SendDailySummary(jobs []EmailJob) {
	sender := os.Getenv("EMAIL_FROM")
	recipient := os.Getenv("EMAIL_TO")

	if sender == "" || recipient == "" || len(jobs) == 0 {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("[Email] Error loading AWS config: %v", err)
		return
	}

	client := sesv2.NewFromConfig(cfg)
	subject := fmt.Sprintf("Daily Job Summary: %d New Roles", len(jobs))

	// --- BUILD HTML BODY ---
	var sb strings.Builder
	
	// CSS Styling
	sb.WriteString("<html><body style='font-family: Arial, sans-serif;'>")
	sb.WriteString("<h2>Good Morning!</h2>")
	sb.WriteString(fmt.Sprintf("<p>We found <b>%d new jobs</b> matching your criteria in the last 24 hours.</p>", len(jobs)))
	sb.WriteString("<hr style='border: 0; border-top: 1px solid #eee;'>")
	
	// Job List
	sb.WriteString("<ul style='padding-left: 0; list-style-type: none;'>")
	for _, job := range jobs {
		sb.WriteString("<li style='margin-bottom: 20px; padding: 10px; background-color: #f9f9f9; border-radius: 5px;'>")
		
		// Title & Company
		sb.WriteString(fmt.Sprintf("<div style='font-size: 16px;'><b>%s</b> @ %s</div>", job.Title, job.Company))
		
		// Location
		if job.Location != "" {
			sb.WriteString(fmt.Sprintf("<div style='color: #666; font-size: 13px;'>📍 %s</div>", job.Location))
		}
		
		// The Hyperlink (Button Style)
		sb.WriteString(fmt.Sprintf("<div style='margin-top: 5px;'><a href='%s' style='color: #0066cc; text-decoration: none; font-weight: bold;'>👉 Apply Now</a></div>", job.URL))
		
		sb.WriteString("</li>")
	}
	sb.WriteString("</ul>")
	
	sb.WriteString("<hr style='border: 0; border-top: 1px solid #eee;'>")
	sb.WriteString("<p style='font-size: 12px; color: #888;'>Sent by DevLink Bot</p>")
	sb.WriteString("</body></html>")


	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(sender),
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Html: &types.Content{Data: aws.String(sb.String())},
					Text: &types.Content{Data: aws.String("Please view this email in an HTML-compatible client.")},
				},
			},
		},
	}

	_, err = client.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("[Email] Failed to send summary: %v", err)
	} else {
		log.Println("[Email] Summary sent successfully.")
	}
}