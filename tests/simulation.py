import json
import random
import socket
import time
from datetime import datetime

import requests

# Config
# ใช้ 127.0.0.1 แทน localhost
API_URL = "http://127.0.0.1:8080/ingest"
UDP_IP = "127.0.0.1"
UDP_PORT = 514

# *** IMPORTANT ***
# เปลี่ยน Tenant ให้ตรงกับ User ที่ Login ใน Dashboard
# จากภาพ Screenshot ดูเหมือนUser จะใช้ Tenant "demo"
# ถ้า Login เป็น Admin ให้ตั้งค่านี้เป็นอะไรก็ได้ (เพราะ Admin เห็นหมด)
TARGET_TENANT = "demo" 

# 1. จำลอง Firewall / Syslog (UDP)
def send_syslog():
    # Syslog ปกติ Backend จะ hardcode เป็น "syslog_default" หรือตั้งค่าแยก
    # แต่ถ้า User เป็น "demo" จะไม่เห็น "syslog_default"
    # ดังนั้นเราจะยิง UDP ปกติ (เผื่อ Admin ดู) 
    # และยิง HTTP เพื่อจำลองให้ User "demo" เห็นด้วย
    
    # 1.1 Syslog UDP Standard
    log_msg = "<134>{} fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.8 spt=5353 dpt=53 proto=udp msg=DNS blocked policy=Block-DNS".format(
        datetime.now().strftime("%b %d %H:%M:%S")
    )
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.sendto(log_msg.encode(), (UDP_IP, UDP_PORT))
        print(f"[UDP] Sent to {UDP_IP}:{UDP_PORT} (Tenant: default)")
    except Exception as e:
        print(f"[UDP] Error: {e}")

    # 1.2 Simulate Syslog via HTTP for "demo" user visibility
    data = {
        "tenant": TARGET_TENANT,
        "source": "syslog-sim", # ตั้งชื่อให้รู้ว่าเป็น Simulation
        "event_type": "DNS_Block",
        "severity": 5,
        "message": log_msg,
        "raw": log_msg,
        "@timestamp": datetime.utcnow().isoformat() + "Z"
    }
    send_http(data)


def send_syslog_tcp():
    log_msg = "<134>{} fw01-tcp action=allow src=192.168.1.99 msg=TCP-Connection-Established\n".format(
        datetime.now().isoformat()
    )
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((UDP_IP, UDP_PORT))
        sock.sendall(log_msg.encode())
        sock.close()
        print(f"[TCP] Sent to {UDP_IP}:{UDP_PORT} (Tenant: default)")
    except Exception as e:
        print(f"[TCP] Error: {e}")
        
    # Simulate TCP via HTTP for visibility
    data = {
        "tenant": TARGET_TENANT,
        "source": "syslog-tcp-sim",
        "event_type": "TCP_Connection",
        "severity": 1,
        "message": log_msg.strip(),
        "raw": log_msg.strip(),
        "@timestamp": datetime.utcnow().isoformat() + "Z"
    }
    send_http(data)


# 2. จำลอง CrowdStrike (HTTP JSON)
def send_crowdstrike():
    data = {
        "tenant": TARGET_TENANT, # แก้เป็น TARGET_TENANT
        "source": "crowdstrike",
        "event_type": "malware_detected",
        "host": "WIN10-01",
        "process": "powershell.exe",
        "severity": 8,
        "sha256": "abc123456789...",
        "action": "quarantine",
        "@timestamp": datetime.utcnow().isoformat() + "Z",
    }
    send_http(data)


# 3. จำลอง AWS CloudTrail (HTTP JSON)
def send_aws():
    data = {
        "tenant": TARGET_TENANT, # แก้เป็น TARGET_TENANT
        "source": "aws",
        "cloud": {
            "service": "iam",
            "account_id": "123456789012",
            "region": "ap-southeast-1",
        },
        "event_type": "CreateUser",
        "user": "admin",
        "@timestamp": datetime.utcnow().isoformat() + "Z",
        "raw": {
            "eventName": "CreateUser",
            "requestParameters": {"userName": "temp-user"},
        },
    }
    send_http(data)


# 4. จำลอง Microsoft 365 (HTTP JSON)
def send_m365():
    data = {
        "tenant": TARGET_TENANT, # แก้เป็น TARGET_TENANT
        "source": "m365",
        "event_type": "UserLoggedIn",
        "user": "bob@demo.local",
        "ip": "198.51.100.23",
        "status": "Success",
        "workload": "Exchange",
        "@timestamp": datetime.utcnow().isoformat() + "Z",
    }
    send_http(data)


# 5. จำลอง AD / Windows Security (HTTP JSON)
def send_ad():
    data = {
        "tenant": TARGET_TENANT, # แก้เป็น TARGET_TENANT
        "source": "ad",
        "event_id": 4625,
        "event_type": "LogonFailed",
        "user": "demo\\eve",
        "host": "DC01",
        "ip": "203.0.113.77",
        "logon_type": 3,
        "@timestamp": datetime.utcnow().isoformat() + "Z",
    }
    send_http(data)


def send_http(data):
    try:
        res = requests.post(API_URL, json=data)
        if res.status_code == 200:
            print(f"[HTTP] Sent {data.get('source')} to Tenant: {data.get('tenant')} => Success")
        else:
            print(f"[HTTP] Failed: {res.status_code} - {res.text}")
    except Exception as e:
        print(f"[HTTP] Error connecting to {API_URL}: {e}")


if __name__ == "__main__":
    print(f"🚀 Starting Log Simulation v2...")
    print(f"🎯 Target Tenant: {TARGET_TENANT} (เพื่อให้ตรงกับ Dashboard)")
    print(f"📡 API: {API_URL}")
    print(f"📡 Syslog: {UDP_IP}:{UDP_PORT}")
    
    while True:
        send_syslog()
        send_syslog_tcp()
        time.sleep(0.5)
        send_crowdstrike()
        time.sleep(0.5)
        send_aws()
        time.sleep(0.5)
        send_m365()
        time.sleep(0.5)
        send_ad()
        print("--- Batch Sent ---")
        time.sleep(3)
