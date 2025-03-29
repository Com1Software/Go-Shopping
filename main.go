package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"time"
)

func getPing(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr // Get the client's IP address and port
	// If the IP address contains the port (e.g., "127.0.0.1:8080"), strip it
	ip, _, _ := net.SplitHostPort(clientIP)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Your IP is: %s", ip) // Return the IP address
}

func getHelp(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Help ")
}
func handleRequests(port string) {
	http.Handle("/ping", http.HandlerFunc(getPing))
	http.Handle("/help", http.HandlerFunc(getHelp))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	port := "8080"
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	xip := fmt.Sprintf("%s", GetOutboundIP())
	fmt.Println("Listening on " + xip + ":" + port)
	fmt.Println("http://" + xip + ":" + port + "/help")

	handleRequests(port)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
