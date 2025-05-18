#!/bin/bash

url="http://localhost:17000"

# Initial POST with fixed body
curl -s -X POST "$url" -H "Content-Type: text/plain" --data-binary $'green\nfigure 0.01 0.01\nupdate'

x=0.01

while (($(echo "$x < 1" | bc -l))); do
  # Prepare body with updated coordinates
  body=$(printf "move %.2f %.2f\nupdate\n" "$x" "$x")

  # Send POST request with updated body
  curl -s -X POST "$url" -H "Content-Type: text/plain" --data-binary "$body"

  # Increment x by 0.01
  x=$(echo "$x + 0.01" | bc)

  # Wait 0.1 seconds
  sleep 0.1
done
