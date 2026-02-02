FROM python:3.10-slim

WORKDIR /app

<<<<<<< HEAD

RUN apt-get update && apt-get install -y --no-install-recommends gcc \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt


COPY . .
=======
# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# We don't specify CMD here because we might want to run it via Cron
# inside the container, or just trigger it manually.
# For now, we set it to run the script directly.
>>>>>>> fb3e5d249783e4b168e22c0761f6d9b8fa143afb
CMD ["python", "main.py"]