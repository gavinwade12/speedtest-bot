#!/bin/bash
#This script will test internet speeds and send them to a local api that will
#send a tweet if the speeds are too low

SPEEDS=(`speedtest-cli | grep -e Download -e Upload | grep -Eo '[0-9]?[0-9]\.[0-9][0-9]'`)
JSON="{\"download\": \"${SPEEDS[0]}\", \"upload\": \"${SPEEDS[1]}\"}"
curl -H "Content-Type: application/json" -X POST -d "$JSON" "http://localhost:8080/speeds"
exit 0
