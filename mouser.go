package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/jeansibelius/mouser/virtualMouse"
)

var addr = flag.String("addr", ":8080", "http service address")
var vMouse = virtualMouse.NewVirtualMouse("Mouser")

func handleMouseHttps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read:", err)
	}

	// Get message (coordinates or button click)
	str := string(body)
	// If we have a comma separated array, it's coordinates
	arr := strings.Split(str, ",")
	if len(arr) > 1 {
		x, _ := strconv.Atoi(arr[0])
		y, _ := strconv.Atoi(arr[1])
		vMouse.Move(x, y)
		fmt.Println("Move:", str)
	} else {
		// Else it's a mouse click
		vMouse.Click(str)
		fmt.Println("Click:", str)
	}
	_, err = w.Write(body)
	if err != nil {
		log.Println("write:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if pusher, ok := w.(http.Pusher); ok {
		// Push is supported.
		if err := pusher.Push("/static/app.js", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/static/style.css", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	htmlTemplate, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing html file: %s\n", err)
	}
	htmlTemplate.Execute(w, "https://"+r.Host+"/mouse")
}

func main() {
	defer vMouse.Close()
	flag.Parse()
	log.SetFlags(0)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/mouse", handleMouseHttps)
	http.HandleFunc("/", home)
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Starting to listen to mouse events at %s%s\n", localAddr.IP.String(), *addr)
	log.Fatal(http.ListenAndServeTLS(*addr, "./cert/server.crt", "./cert/server.key", nil))
}
