import requests
import json
import sys
import time
import os

# Configuration
API_URL = "https://localhost/ingest"
TENANT_ID = "tenant_default"
SOURCE = "file_upload"

def ingest_file(filepath):
    if not os.path.exists(filepath):
        print(f"❌ File not found: {filepath}")
        return

    print(f"🚀 Starting ingestion for: {filepath}")
    
    batch = []
    batch_size = 10
    total_sent = 0

    with open(filepath, 'r', encoding='utf-8') as f:
        for line in f:
            line = line.strip()
            if not line:
                continue

            # Construct Log Entry
            entry = {
                "tenant": TENANT_ID,
                "source": SOURCE,
                "event_type": "raw_log",
                "severity": 1,
                "body": {
                    "message": line
                }
            }
            batch.append(entry)

            # Send batch if full
            if len(batch) >= batch_size:
                send_batch(batch)
                total_sent += len(batch)
                batch = []
                time.sleep(0.1) # Prevent flooding

    # Send remaining
    if batch:
        send_batch(batch)
        total_sent += len(batch)

    print(f"✅ Ingestion complete! Total logs sent: {total_sent}")

def send_batch(logs):
    headers = {'Content-Type': 'application/json'}
    try:
        # verify=False because of self-signed cert
        response = requests.post(API_URL, data=json.dumps(logs), headers=headers, verify=False)
        if response.status_code != 200:
            print(f"⚠️ Error sending batch: {response.status_code} - {response.text}")
        else:
            print(f"host -> sent {len(logs)} logs")
    except Exception as e:
        print(f"❌ Connection error: {e}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python file_ingestor.py <path_to_log_file>")
        print("Example: python file_ingestor.py ./sample.log")
    else:
        ingest_file(sys.argv[1])
