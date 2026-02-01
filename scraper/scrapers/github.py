import requests
from bs4 import BeautifulSoup

def scrape_simplify_repo(repo_url):
    """
    Scrapes job tables from SimplifyJobs-style GitHub repositories.
    """
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    
    jobs_found = []
    
    try:
        print(f"Checking {repo_url}...")
        response = requests.get(repo_url, headers=headers)
        if response.status_code != 200:
            print(f"Failed to fetch {repo_url}: {response.status_code}")
            return []

        soup = BeautifulSoup(response.text, "html.parser")
        
    
        tables = soup.find_all("table")
        
        target_table = None
        for table in tables:
            headers_text = [th.get_text(strip=True).lower() for th in table.find_all("th")]
            if "company" in headers_text and "role" in headers_text:
                target_table = table
                break
        
        if not target_table:
            print(f"No job table found in {repo_url}")
            return []

        
        rows = target_table.find_all("tr")[1:] 
        
        for row in rows:
            cols = row.find_all("td")
            if len(cols) < 3:
                continue

            company = cols[0].get_text(strip=True)
            role = cols[1].get_text(strip=True)
            location = cols[2].get_text(strip=True)
        
            link_elem = None
        
            if len(cols) >= 4:
                link_elem = cols[3].find("a")
    
            if not link_elem:
                link_elem = cols[1].find("a")

            if link_elem and link_elem.get("href"):
                raw_url = link_elem["href"]
                
       
                if "github.com" in raw_url and "SimplifyJobs" in raw_url:
                    continue 

                jobs_found.append({
                    "title": role,
                    "company": company,
                    "location": location,
                    "url": raw_url
                })

    except Exception as e:
        print(f"Error scraping {repo_url}: {e}")

    return jobs_found