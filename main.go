package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	ws *websocket.Conn
}

type Msg struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

var wsClient Client

var subscriber mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	sendMessage(msg.Topic(), msg.Payload())
}

func sendMessage(topic string, message []byte) {
	msg := &Msg{
		Action:  topic,
		Message: string(message),
	}

	json, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error json marshal", err)
	}

	wsClient.ws.WriteMessage(websocket.TextMessage, []byte(json))
}

func main() {
	var addr string
	var broker string
	var username string
	var password string
	flag.StringVar(&addr, "addr", "127.0.0.1:1234", "IP Adrress to listen on")
	flag.StringVar(&broker, "broker", "127.0.0.1:1883", "MQTT Broker address")
	flag.StringVar(&username, "username", "", "Broker username")
	flag.StringVar(&password, "password", "", "Broker password")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", broker))
	opts.SetClientID("go-client")
	opts.SetUsername(username)
	opts.SetPassword(password)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer c.Disconnect(250)

	log.Printf("Connected to broker at %s\n", broker)

	http.Handle("/", http.FileServer(http.Dir("./html")))

	http.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		topic := "esp32/led/status"
		payload := "off"
		token := c.Publish(topic, 0, true, payload)

		if token.Error() != nil {
			log.Println(token.Error())
			return
		}

		sendMessage(topic, []byte(payload))
	})

	http.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		topic := "esp32/led/status"
		payload := "on"
		token := c.Publish(topic, 0, true, payload)

		if token.Error() != nil {
			log.Println(token.Error())
			return
		}

		sendMessage(topic, []byte(payload))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		wsClient = Client{ws}

		log.Println("Client Connected")

		if token := c.Subscribe("esp32/status", 0, subscriber); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := c.Subscribe("esp32/led/status", 0, subscriber); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	})

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
