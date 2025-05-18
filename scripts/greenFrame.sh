#!/bin/bash

curl -X POST http://localhost:17000 \
  -H "Content-Type: text/plain" \
  --data-binary $'reset\nwhite\nbgrect 0.25 0.25 0.75 0.75\nfigure 0.5 0.5\ngreen\nfigure 0.6 0.6\nupdate'
