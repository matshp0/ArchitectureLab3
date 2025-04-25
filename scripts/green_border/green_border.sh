#!/bin/bash

# Команди для створення зеленої рамки
curl -X POST -d "green" http://localhost:8080
curl -X POST -d "bgrect 0 0 400 400" http://localhost:8080
curl -X POST -d "update" http://localhost:8080 