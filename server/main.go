package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var fileName = flag.String("file", "data.txt", "file name")

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/code", getCode)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func getCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := os.Stat(*fileName); err != nil {
		os.WriteFile(*fileName, []byte("//write your code here"), 0644)
	}
	data, err := os.ReadFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	if _, err := os.Stat(*fileName); err == nil {
		data, err := os.ReadFile(*fileName)
		if err != nil {
			log.Fatal(err)
		}
		c.WriteMessage(websocket.TextMessage, data)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		os.WriteFile(*fileName, message, 0644)
	}
}
