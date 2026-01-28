import requests
import json
import time

url = "https://localhost/ingest"

# Override verify=False because of self-signed cert
headers = {'Content-Type': 'application/json'}

logs = [
    {
        "tenant": "tenant_A",
        "source": "payment_service",
        "event_type": "transaction_success",
        "severity": 1,
        "body": {
             "amount": 500,
             "user_id": "u123",
             "currency": "USD"
        }
    },
    {
        "tenant": "tenant_A",
        "source": "payment_service",
        "event_type": "transaction_failed",
        "severity": 7,
        "body": {
             "amount": 12000,
             "user_id": "u999",
             "error": "insufficient_funds"
        }
    }
]

try:
    print(f"Sending {len(logs)} logs to {url}...")
    response = requests.post(url, data=json.dumps(logs), headers=headers, verify=False)
    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.text}")
except Exception as e:
    print(f"Error: {e}")
