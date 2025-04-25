package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	// Команди для створення зеленої рамки
	commands := []string{
		"green",
		"bgrect 0 0 400 400",
		"update",
	}

	// Відправляємо кожну команду на сервер
	for _, cmd := range commands {
		resp, err := http.Post("http://localhost:8080", "text/plain", strings.NewReader(cmd))
		if err != nil {
			fmt.Printf("Error sending command %s: %v\n", cmd, err)
			return
		}
		resp.Body.Close()
	}
} 