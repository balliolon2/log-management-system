import time

import requests

for i in range(10):  # ส่ง 10 ครั้ง
    requests.post(
        "http://localhost:8080/ingest",
        json={
            "tenant": "demo",
            "source": "api",
            "event_type": "login_failed",
            "user": "alice",
            "body": {"src_ip": "192.168.1.100"},
            "@timestamp": "2025-01-28T10:00:00Z",
        },
    )
    print(f"Sent log {i + 1}")
    time.sleep(1)
