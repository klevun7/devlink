FROM python:3.10-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# We don't specify CMD here because we might want to run it via Cron
# inside the container, or just trigger it manually.
# For now, we set it to run the script directly.
CMD ["python", "main.py"]