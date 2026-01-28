import time
from datetime import datetime

import requests


# Test HTTP Ingestion
def test_http_ingestion():
    url = "http://localhost:8080/ingest"

    logs = [
        {
            "tenant": "demo",
            "source": "api",
            "event_type": "login_success",
            "user": "alice",
            "body": {"src_ip": "192.168.1.100"},
            "@timestamp": datetime.utcnow().isoformat() + "Z",
        },
        {
            "tenant": "demo",
            "source": "api",
            "event_type": "login_failed",
            "user": "bob",
            "body": {"src_ip": "192.168.1.101"},
            "@timestamp": datetime.utcnow().isoformat() + "Z",
        },
    ]

    for log in logs:
        response = requests.post(url, json=log)
        print(f"Status: {response.status_code} | Log: {log['event_type']}")
        time.sleep(0.5)


if __name__ == "__main__":
    test_http_ingestion()
