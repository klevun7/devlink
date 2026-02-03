import requests
from bs4 import BeautifulSoup
from datetime import datetime, timedelta
import re

def parse_github_date(date_str):

    today = datetime.now()
    clean_str = re.sub(r'[^a-zA-Z0-9 ]', '', date_str).strip().lower()

    try:
        # CASE 1: Relative "Months" (e.g., "1mo")
        if 'mo' in clean_str:
            match = re.search(r'\d+', clean_str)
            if match:
                num = int(match.group())
                return today - timedelta(days=num * 30)

        # CASE 2: Relative "Days" (e.g., "5d")
        elif 'd' in clean_str and 'm' not in clean_str:
            match = re.search(r'\d+', clean_str)
            if match:
                num = int(match.group())
                return today - timedelta(days=num)

        # CASE 3: Relative "Hours" (e.g., "12h")
        elif 'h' in clean_str:
            return today 

        # CASE 4: Standard "Jan 20"
        else:
            current_year = today.year
            clean_str = clean_str.title() 
            if len(clean_str) < 3: return None
            
            date_obj = datetime.strptime(f"{clean_str} {current_year}", "%b %d %Y")
            
            if date_obj > today + timedelta(days=5):
                 date_obj = date_obj.replace(year=current_year - 1)
            return date_obj

    except Exception:
        return None
    
    return None

def scrape_simplify_repo(repo_url):
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    
    jobs_found = []
    
    try:
        response = requests.get(repo_url, headers=headers)
        soup = BeautifulSoup(response.text, "html.parser")
        
        tables = soup.find_all("table")
        target_table = None
        
        for table in tables:
            headers = [th.get_text(strip=True).lower() for th in table.find_all("th")]
            if "company" in headers:
                target_table = table
                break
        
        if not target_table:
            print("[ERROR] No Job Table found.")
            return []

        rows = target_table.find_all("tr")[1:] 
        today = datetime.now()
        
        last_company = "Unknown"

        for row in rows:
            cols = row.find_all("td")
            if len(cols) < 4: continue

            raw_company = cols[0].get_text(strip=True)
            if "↳" in raw_company or not raw_company:
                company_name = last_company
            else:
                company_name = raw_company
                last_company = raw_company

            role_name = cols[1].get_text(strip=True)
            location = cols[2].get_text(strip=True)
            

            date_str = cols[-1].get_text(strip=True)
            job_date = parse_github_date(date_str)

            if job_date:
                days_old = (today - job_date).days
                
                if days_old > 5:
                    continue 

                link_elem = None
                if len(cols) >= 4: link_elem = cols[3].find("a")
                if not link_elem: link_elem = cols[1].find("a")

                if link_elem:
                    jobs_found.append({
                        "title": role_name,
                        "company": company_name,
                        "location": location,
                        "url": link_elem['href'],
                        "posted_date": date_str
                    })
                    print(f"   [FOUND] {company_name} ({date_str})")

    except Exception as e:
        print(f"[ERROR] Scraper failed: {e}")

    return jobs_found