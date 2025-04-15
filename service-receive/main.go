package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Định nghĩa hàm callback khi nhận tin nhắn
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📥 Received on [%s]: %s\n", msg.Topic(), msg.Payload())
	}

	// Cấu hình MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://103.56.158.48:1883").
		SetClientID("go_mqtt_client").
		SetUsername("test").
		SetPassword("test").
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
