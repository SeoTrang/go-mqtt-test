
# MQTT Client and Server in Go

## Mục lục
1. [Giới thiệu](#giới-thiệu)
2. [Cài đặt](#cài-đặt)
3. [Chi tiết và giải thích](#chi-tiết-và-giải-thích)
   - [MQTT Client](#mqtt-client)
   - [MQTT Server](#mqtt-server)
   - [Câu hỏi thường gặp](#câu-hỏi-thường-gặp)

---

## Giới thiệu

Dự án này bao gồm hai phần chính:
1. **MQTT Client**: Được viết bằng Go, giúp kết nối và gửi thông điệp đến một MQTT broker.
2. **MQTT Server**: Được viết bằng Go, giúp nhận và xử lý thông điệp từ một MQTT broker.

Cả hai dự án này giao tiếp thông qua MQTT, một giao thức nhắn tin phổ biến cho các ứng dụng IoT.

---

## Cài đặt

1. Đảm bảo rằng đã cài đặt **Go** và có sẵn môi trường phát triển.
2. Clone repository này về:

   ```bash
   git clone https://github.com/SeoTrang/go-mqtt-test.git
   cd go-mqtt-test
   ```

3. Cài đặt thư viện MQTT cho Go:

   ```bash
   go get github.com/eclipse/paho.mqtt.golang
   ```

4. Chạy các file Go:

   - Để chạy **MQTT Client**:

     ```bash
     go run service-push/main.go
     ```

   - Để chạy **MQTT Server**:

     ```bash
     go run service-receive/main.go
     ```

---

## Chi tiết và giải thích

### MQTT Client

**File: `service-push/main.go`**

Đoạn code này giúp bạn kết nối tới một MQTT broker và gửi thông điệp đến topic `test/topic`.

```go
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
        AddBroker("tcp://{mqtt-ip}:1883").
        SetClientID("go_simple_client").
        SetUsername("test").
        SetPassword("test")

    client := mqtt.NewClient(opts)

    // Kết nối đến broker
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
```

**Giải thích**:
- **`mqtt.NewClientOptions()`**: Tạo các tùy chọn cho client, bao gồm địa chỉ broker, client ID, username và password.
- **`client.Connect()`**: Kết nối tới broker. Nếu kết nối thành công, client sẽ tiếp tục.
- **`client.Publish()`**: Gửi thông điệp JSON đến topic `test/topic`.
- **`client.Disconnect(250)`**: Đóng kết nối sau khi gửi xong.

#### Các câu hỏi giải đáp:
- **Tại sao dùng `time.Sleep(1 * time.Second)` hoặc `select {}` trong MQTT server?**
  - Trong Go, nếu không giữ chương trình chạy, client sẽ tự động thoát sau khi hoàn thành một lần gửi message. Để giữ chương trình chạy để nhận tin nhắn, chúng ta sử dụng `time.Sleep` (hoặc `select {}`) để tạo vòng lặp vô hạn, giúp chương trình lắng nghe tin nhắn mà không thoát.

---

### MQTT Server

**File: `service-receive/main.go`**

Đoạn mã này giúp bạn kết nối đến broker và lắng nghe các tin nhắn từ topic `test/topic`.

```go
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
        AddBroker("tcp://{mqtt-ip}:1883").
        SetClientID("go_mqtt_client").
        SetUsername("test").
        SetPassword("test").
        SetDefaultPublishHandler(messageHandler)

    // Tạo client
    client := mqtt.NewClient(opts)

    // Kết nối đến broker
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

    // Chờ nhận message (vòng lặp vô hạn)
    for {
        time.Sleep(1 * time.Second)
    }
}
```

**Giải thích**:
- **`SetDefaultPublishHandler()`**: Đặt handler cho việc nhận thông điệp. Mỗi khi có message mới từ topic, hàm `messageHandler` sẽ được gọi.
- **`client.Subscribe()`**: Đăng ký để lắng nghe các thông điệp từ topic `test/topic`.
- **Vòng lặp `for {}`**: Giữ chương trình không thoát, giúp chương trình tiếp tục lắng nghe tin nhắn.

#### Các câu hỏi giải đáp:
- **Tại sao phải dùng `time.Sleep(1 * time.Second)` hoặc `select {}`?**
  - Đây là cách để chương trình **không thoát** sau khi đăng ký lắng nghe. Nếu không có vòng lặp này, chương trình sẽ dừng ngay sau khi đăng ký mà không nhận được tin nhắn.

---

Cảm ơn bạn đã chỉ ra sai sót. Dưới đây là phiên bản chính xác và chi tiết hơn cho phần giải thích về lý do tại sao trong Node.js không cần chỉ định port, còn trong Go lại cần chỉ định và sự khác biệt giữa `mqtt://` và `tcp://`:

---

## Câu hỏi thường gặp

### 1. **Tại sao không cần chỉ định port trong URL của MQTT client trong Node.js nhưng lại cần trong Go?**
   - **Node.js**: Khi sử dụng thư viện `mqtt` cho Node.js, nếu không chỉ định port, thư viện này sẽ tự động sử dụng port mặc định `1883` cho giao thức MQTT không mã hóa (TCP). Vì vậy, khi bạn sử dụng `mqtt://{mqtt-ip}`, Node.js sẽ mặc định kết nối tới `{mqtt-ip}:1883`.
   - **Go**: Trong Go, thư viện `paho.mqtt.golang` yêu cầu bạn phải chỉ định rõ ràng **port** khi kết nối đến broker MQTT. Do đó, bạn phải chỉ định cổng rõ ràng, ví dụ `tcp://{mqtt-ip}:1883`. Go không tự động giả định port mặc định như Node.js.

### 2. **Tại sao trong Node.js dùng `mqtt://`, còn trong Go dùng `tcp://`?**
   - **Node.js**: Thư viện MQTT của Node.js mặc định sử dụng giao thức **TCP** cho kết nối không mã hóa khi bạn sử dụng `mqtt://`. Bạn không cần chỉ định rõ giao thức TCP vì thư viện đã hiểu mặc định là sử dụng giao thức TCP.
   - **Go**: Trong Go, bạn cần chỉ định rõ giao thức khi kết nối đến broker, ví dụ: `tcp://` cho kết nối không mã hóa hoặc `tls://` cho kết nối mã hóa (SSL/TLS). Vì vậy, khi bạn sử dụng `tcp://{mqtt-ip}:1883`, bạn đang nói rõ rằng bạn muốn sử dụng giao thức **TCP**.

---