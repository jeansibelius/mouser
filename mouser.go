package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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
	} else {
		// Else it's a mouse click
		vMouse.Click(str)
	}
	_, err = w.Write(body)
	if err != nil {
		log.Println("write:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	htmlTemplate, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing html file: %s\n", err)
	}
	//htmlTemplate.Execute(w, "ws://"+r.Host+"/mouse")
	htmlTemplate.Execute(w, "https://"+r.Host+"/mouse")
}

func main() {
	defer vMouse.Close()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/mouse", handleMouseHttps)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServeTLS(*addr, "./cert/server.crt", "./cert/server.key", nil))
}
