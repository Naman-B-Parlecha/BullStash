package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func SendToDiscord(webhook string, message string) error {
	payload := map[string]string{"content": message}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook error: %s", resp.Status)
	}
	return nil
}

func CallWebHook() {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	message := fmt.Sprintf("Task executed at %s\nOutput:\n```\n%s\n```",
		time.Now().Format(time.RFC1123),
		string("salam"))

	if err := SendToDiscord(webhookURL, message); err != nil {
		fmt.Printf("Failed to send to Discord: %v\n", err)
	}
}
