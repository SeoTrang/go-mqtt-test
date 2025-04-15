package main

import (
	"encoding/json"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Cấu hình MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://103.56.158.48:1883").
		SetClientID("go_simple_client").
		SetUsername("test").
		SetPassword("test")

	client := mqtt.NewClient(opts)

	// Kết nối MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("❌ MQTT connection failed:", token.Error())
		os.Exit(1)
	}

	fmt.Println("✅ MQTT client connected")

	// Tạo payload giống như JS
	payload := map[string]string{
		"from":    "Simple Go Client",
		"message": "Hello from MQTT client!",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	// Gửi message
	token := client.Publish("test/topic", 1, false, payloadBytes)
	token.Wait() // Chờ gửi xong

	if token.Error() != nil {
		fmt.Println("❌ Publish failed:", token.Error())
	} else {
		fmt.Println("📤 Message published to test/topic")
	}

	// Đóng kết nối sau khi gửi
	client.Disconnect(250)
}
