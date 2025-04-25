#!/bin/bash

# Спочатку створюємо зелений фон
curl -X POST -d "green" http://localhost:8080
curl -X POST -d "bgrect 0 0 400 400" http://localhost:8080
curl -X POST -d "update" http://localhost:8080

# Малюємо фігуру T-90
curl -X POST -d "blue" http://localhost:8080
curl -X POST -d "figure 200 200" http://localhost:8080
curl -X POST -d "update" http://localhost:8080

# Переміщуємо фігуру по діагоналі
for i in {1..10}; do
    x=$((200 + i * 10))
    y=$((200 + i * 10))
    curl -X POST -d "move $x $y" http://localhost:8080
    curl -X POST -d "update" http://localhost:8080
    sleep 1
done 