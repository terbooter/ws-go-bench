package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func main() {
	flag.Parse()
	fmt.Println(*addr)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func echo(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("recv: %s", message)
		start := time.Now()
		for i := 0; i < 10000; i++ {
			err = c.WriteMessage(mt, []byte("message : "+strconv.Itoa(i)))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

		totalInterval := time.Now().Sub(start)
		fmt.Println("Completed in :" + totalInterval.String())
		err = c.WriteMessage(mt, []byte("Completed in :"+totalInterval.String()))
		if err != nil {
			log.Println("write:", err)
			break
		}

	}
}

//go:embed index.html
var content []byte

var homeTemplate = template.Must(template.New("").Parse(string(content)))
