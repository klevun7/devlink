import os
import requests
import time
import datetime
from scrapers.github import scrape_simplify_repo


API_URL = os.getenv("API_URL", "http://localhost:8080") 
TARGET_REPOS = [
    "https://github.com/SimplifyJobs/Summer2025-Internships",
    "https://github.com/SimplifyJobs/New-Grad-Positions"
]

def post_job(job):
    try:
        response = requests.post(f"{API_URL}/jobs", json=job)
        if response.status_code == 201:
            print(f"[NEW] {job['company']} - {job['title']}")
            return job
        elif response.status_code == 409:
            return None 
    except Exception as e:
        print(f"[ERR] API Connection failed: {e}")
        return None

def trigger_daily_email(new_jobs):
    if not new_jobs:
        print("No new jobs today. Skipping email.")
        return

    print(f"Triggering email for {len(new_jobs)} jobs...")
    try:
        requests.post(f"{API_URL}/notifications/daily", json=new_jobs)
    except Exception as e:
        print(f"Failed to trigger email: {e}")

def run_cycle():
    print(f"--- Starting Daily Scrape: {datetime.datetime.now()} ---")
    collected_jobs = []

    for repo in TARGET_REPOS:
        found_jobs = scrape_simplify_repo(repo)
        
        for job in found_jobs:
            saved_job = post_job(job)
            if saved_job:
                collected_jobs.append(saved_job)

    trigger_daily_email(collected_jobs)
    print("--- Cycle Complete. Exiting Python Process. ---")

if __name__ == "__main__":
        run_cycle()
