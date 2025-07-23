# devlink
A real-time job alert platform that delivers New-Grad and Internship postings to developers using Go and AWS SES.

## Decription
Devlink implements an automated system that scrapes job boards daily, identifies new job postings, and sends personalized email alerts to subscribed users. 
Built with Go for the backend API, Python with BeautifulSoup for scraping, and powered by AWS SES for email delivery, it ensures New-Grads and prospective interns
receive timely updates directly to their inbox.

## System Architecture
![System Diagram](/system-diagram.png)

## Technologies
`Go`: For the robust and scalable REST API server.

`Python`: For the web scraping and data aggregation logic (using BeautifulSoup).

`SQLite`: A lightweight, file-based database used for storing job data and user preferences.

`AWS SES (Simple Email Service)`: A cost-effective and scalable email sending service.

`Docker`: For containerizing the Go API server, ensuring consistent deployment.

`GitHub Actions`: For automating the daily cron job that triggers the scraping process.
