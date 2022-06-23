package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/jeansibelius/mouser/virtualMouse"
)

var addr = flag.String("addr", "192.168.1.5:8080", "http service address")
var vMouse = virtualMouse.NewVirtualMouse("Mouser")

var upgrader = websocket.Upgrader{} // use default options

func handleMouse(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// Get message (coordinates or button click)
		str := string(message)
		// If we have a comma separated array, it's coordinates
		arr := strings.Split(str, ",")
		if len(arr) > 1 {
			x, _ := strconv.Atoi(arr[0])
			y, _ := strconv.Atoi(arr[1])
			vMouse.Move(x, y)
		} else {
			// Else it's a mouse click
			vMouse.Click(str)
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

var htmlTemplate, _ = template.ParseFiles("index.html")

func home(w http.ResponseWriter, r *http.Request) {
	htmlTemplate.Execute(w, "ws://"+r.Host+"/mouse")
}

func main() {
	defer vMouse.Close()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/mouse", handleMouse)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
