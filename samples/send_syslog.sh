#!/bin/bash
# send_syslog.sh - Send a test UDP Syslog message to localhost
# Usage: ./send_syslog.sh "Your message here"

MESSAGE=${1:-"Test message from script"}
SERVER="localhost"
PORT=514

echo "Sending to $SERVER:$PORT -> $MESSAGE"
echo -n "$MESSAGE" | nc -u -w1 $SERVER $PORT
