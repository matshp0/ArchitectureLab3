package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Спочатку створюємо зелений фон
	commands := []string{
		"green",
		"bgrect 0 0 400 400",
		"update",
	}

	// Відправляємо команди для створення фону
	for _, cmd := range commands {
		resp, err := http.Post("http://localhost:8080", "text/plain", strings.NewReader(cmd))
		if err != nil {
			fmt.Printf("Error sending command %s: %v\n", cmd, err)
			return
		}
		resp.Body.Close()
	}

	// Малюємо фігуру T-90
	tCommands := []string{
		"blue",
		"figure 200 200",
		"update",
	}

	// Відправляємо команди для малювання фігури
	for _, cmd := range tCommands {
		resp, err := http.Post("http://localhost:8080", "text/plain", strings.NewReader(cmd))
		if err != nil {
			fmt.Printf("Error sending command %s: %v\n", cmd, err)
			return
		}
		resp.Body.Close()
	}

	// Переміщуємо фігуру по діагоналі
	x, y := 200, 200
	for i := 0; i < 10; i++ {
		x += 10
		y += 10
		moveCmd := fmt.Sprintf("move %d %d", x, y)
		resp, err := http.Post("http://localhost:8080", "text/plain", strings.NewReader(moveCmd))
		if err != nil {
			fmt.Printf("Error sending command %s: %v\n", moveCmd, err)
			return
		}
		resp.Body.Close()

		// Оновлюємо екран
		resp, err = http.Post("http://localhost:8080", "text/plain", strings.NewReader("update"))
		if err != nil {
			fmt.Printf("Error sending update command: %v\n", err)
			return
		}
		resp.Body.Close()

		time.Sleep(time.Second)
	}
} 