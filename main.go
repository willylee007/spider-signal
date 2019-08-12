package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

var addr = flag.String("addr", ":8080", "websocket service")

func main() {

	flag.Parse()
	log.Println("let's go!")
	hub := newHub()
	go hub.run()
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})
	r.GET("/test", func(c *gin.Context) {
		c.Writer.WriteString("hello world")
	})
	err := r.Run(*addr)
	// err := r.RunTLS(*addr, "1.crt", "1.key")
	if err != nil {
		log.Fatal("ListenAndServe err: ", err)
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "Hello, TLS!\n")
	// })
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// err := http.ListenAndServe(*addr, nil)
	// err := http.ListenAndServeTLS(*addr, "1.crt", "1.key", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}
