# DevLink - Automated Job Aggregator & Alert System

DevLink is a microservices-based data pipeline that aggregates software engineering job postings from various sources, filters them for relevance, and broadcasts daily email summaries to subscribers. It is designed to be self-hosting, idempotent, and scalable.

## System Architecture

The system follows a Producer-Consumer pattern utilizing Docker for containerization and orchestration.

* **Producer (Scraper):** A Python-based service that runs on a scheduled interval (cron). It extracts data from GitHub repositories, cleans/normalizes the data, and transmits it to the internal API.
* **Consumer/Gatekeeper (API):** A Go (Golang) REST API that validates incoming data, enforces schema constraints, and manages persistence. It also handles the business logic for dispatching notifications.
* **Storage (Database):** A PostgreSQL instance used for persistent storage of job listings and subscriber data. It enforces unique constraints to prevent duplicate entries.
* **Proxy (Nginx):** A reverse proxy that manages ingress traffic, handling port forwarding and security for the public-facing API.
* **Notification Service:** Integrated with AWS SES (Simple Email Service) to deliver HTML-formatted email alerts using IAM Instance Profiles for secure authentication.

## Tech Stack

* **Backend API:** Go (Golang), Gorm (ORM), Standard Library net/http
* **ETL/Scraper:** Python, BeautifulSoup4, Requests
* **Database:** PostgreSQL
* **Infrastructure:** Docker, Docker Compose, AWS EC2
* **Networking:** Nginx Reverse Proxy
* **Cloud Services:** AWS SES, AWS IAM

## Key Features

1.  **Hybrid Microservices Architecture:** Leverages Python's ecosystem for efficient text processing and scraping while utilizing Go's concurrency and strict typing for the high-performance API.
2.  **Idempotent Data Ingestion:** The database schema employs unique constraints on job URLs. The API gracefully handles duplicate POST requests (returning 409 Conflict), ensuring the scraper can run multiple times without corrupting the dataset.
3.  **Resilient Networking:** Internal services communicate via a private Docker network (DNS service discovery), while public traffic is securely routed through Nginx.
4.  **Asynchronous Notifications:** Email dispatching is handled via Go goroutines, decoupling the user response time from the SMTP latency.
5.  **Secure Configuration:** Uses environment variables for all secrets and AWS IAM Roles for EC2, eliminating the need for long-lived credentials on the server.

## Getting Started

### Prerequisites

* Docker & Docker Compose
* Git

### Local Installation

1.  **Clone the repository**
    ```bash
    git clone [https://github.com/yourusername/devlink.git](https://github.com/yourusername/devlink.git)
    cd devlink
    ```

2.  **Configure Environment Variables**
    Create a .env file in the root directory:
    ```ini
    # Database Configuration
    DB_HOST=db
    DB_USER=devlink_user
    DB_PASS=your_secure_password
    DB_NAME=devlink
    DB_PORT=5432

    # API Configuration
    API_PORT=8080
    APP_ENV=development

    # AWS SES Configuration (Sandbox Mode)
    AWS_REGION=us-east-1
    EMAIL_FROM=verified_sender@example.com
    EMAIL_TO=verified_receiver@example.com
    ```

3.  **Build and Run**
    ```bash
    docker compose up --build -d
    ```

4.  **Verify Status**
    The API will be available at http://localhost:8080/jobs.

## API Documentation

### GET /jobs
Retrieves a list of the most recent job postings.

**Response:**
```json
[
  {
    "title": "Software Engineer New Grad",
    "company": "Tech Corp",
    "url": "[https://apply.example.com](https://apply.example.com)",
    "location": "San Francisco, CA",
    "posted_at": "2023-10-25T10:00:00Z"
  }
]
```
### POST /jobs
Internal endpoint used by the scraper to ingest new data.

**Payload:**
```json
{
  "title": "Backend Developer",
  "company": "Startup Inc",
  "url": "[https://apply.startup.com/job/123](https://apply.startup.com/job/123)",
  "location": "Remote"
}
```
### POST /notifications/daily
Triggers the daily email summary. Accepts a list of new jobs to include in the report.
