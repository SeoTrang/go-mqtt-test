package main

import (
	"fmt"
	"os"
	"time"

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
	// Định nghĩa hàm callback khi nhận tin nhắn
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📥 Received on [%s]: %s\n", msg.Topic(), msg.Payload())
	}

	// Cấu hình MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", brokerIP, brokerPort)).
		SetClientID("go_mqtt_client").
		SetUsername(brokerUsername).
		SetPassword(brokerPassword).
		SetDefaultPublishHandler(messageHandler)

	// Tạo client
	client := mqtt.NewClient(opts)

	// Kết nối đến broker
	// token := client.Connect()
	// token.Wait()
	// if token.Error() != nil {
	// 	fmt.Println("❌ MQTT connection failed:", token.Error())
	// 	os.Exit(1)
	// }

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("❌ MQTT connection failed:", token.Error())
		os.Exit(1)
	}

	fmt.Println("✅ MQTT Server connected and ready")

	// Đăng ký topic
	if token := client.Subscribe("test/topic", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println("❌ Failed to subscribe:", token.Error())
		os.Exit(1)
	}
	fmt.Println("📡 Subscribed to topic: test/topic")

	// Chờ nhận message (ở đây để chạy vô hạn – hoặc bạn có thể dùng select{})
	for {
		time.Sleep(1 * time.Second)
	}
}
