package main

import (
	"encoding/json"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ Không thể load file .env:", err)
		os.Exit(1)
	}

	brokerIP := os.Getenv("BROKER_IP")
	brokerPort := os.Getenv("BROKER_PORT")
	brokerUsername := os.Getenv("BROKER_USERNAME")
	brokerPassword := os.Getenv("BROKER_PASSWORD")

	// Cấu hình MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", brokerIP, brokerPort)).
		SetClientID("go_mqtt_client_push").
		SetUsername(brokerUsername).
		SetPassword(brokerPassword)

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
