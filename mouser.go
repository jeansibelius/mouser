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
		str := string(message)
		arr := strings.Split(str, ",")
		if len(arr) > 1 {
			x, _ := strconv.Atoi(arr[0])
			y, _ := strconv.Atoi(arr[1])
			vMouse.Move(x, y)
		} else {
			vMouse.Click(str)
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/mouse")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
		const getCoordinates = function(x, y) {
				return String(x)+","+String(y);
		};
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = async function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("up").onclick = function(evt) {
        if (!ws) {
            return false;
        }
				const coord = getCoordinates(0, input.value*-1);
        print("UP: " + coord);
        ws.send(coord);
        return false;
    };
    document.getElementById("down").onclick = function(evt) {
        if (!ws) {
            return false;
        }
				const coord = getCoordinates(0, input.value);
        print("DOWN: " + coord);
        ws.send(coord);
        return false;
    };
    document.getElementById("left").onclick = function(evt) {
        if (!ws) {
            return false;
        }
				const coord = getCoordinates(input.value*-1, 0);
        print("LEFT: " + coord);
        ws.send(coord);
        return false;
    };
    document.getElementById("right").onclick = function(evt) {
        if (!ws) {
            return false;
        }
				const coord = getCoordinates(input.value, 0);
        print("RIGHT: " + coord);
        ws.send(coord);
        return false;
    };
    document.getElementById("leftClick").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("CLICK: LEFT");
        ws.send("left");
        return false;
    };
    document.getElementById("rightClick").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("CLICK: RIGHT");
        ws.send("right");
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
</form>
<br/>
<form>
<button id="leftClick">Left Click</button>
<button id="rightClick">Right Click</button>
</form>
<br/>
<form>
<button id="up">Up</button>
<button id="down">Down</button>
<button id="left">Left</button>
<button id="right">Right</button>
<br/>
<input id="input" type="number" value="100">
</form>
<br/>
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</body>
</html>
`))

func main() {
	defer vMouse.Close()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/mouse", handleMouse)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
