import os
import requests
import time
from scrapers.github import scrape_simplify_repo


API_URL = os.getenv("API_URL", "http://localhost:8080/jobs")


TARGET_REPOS = [
    "https://github.com/SimplifyJobs/New-Grad-Positions"
]

def post_job(job):
    """
    Sends the job to the Go API.
    Returns True if saved (new), False if duplicate or error.
    """
    try:
        response = requests.post(API_URL, json=job)
        
        if response.status_code == 201:
            print(f"[NEW] {job['company']} - {job['title']}")
            return True
        elif response.status_code == 409:
            return False 
        else:
            print(f"[ERR] API {response.status_code}: {response.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"[FAIL] Could not connect to API: {e}")
        return False

def run():
    print("--- Starting Scrape Cycle ---")
    
    for repo in TARGET_REPOS:
        jobs = scrape_simplify_repo(repo)
        print(f"Found {len(jobs)} potential jobs in {repo.split('/')[-1]}")
        
        new_count = 0
        for job in jobs:
            if post_job(job):
                new_count += 1
        
        print(f"Cycle complete. Added {new_count} new jobs from this repo.")

if __name__ == "__main__":
    run()
    # while True:
    #     run()
    #     time.sleep(3600) # Sleep 1 hour