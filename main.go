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
	clientIP := r.RemoteAddr
	ip, _, _ := net.SplitHostPort(clientIP)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Your IP is: %s", ip)
}

func getHelp(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	w.WriteHeader(http.StatusOK)
	port := GetOutboundPort()
	xip := fmt.Sprintf("%s", GetOutboundIP())
	fmt.Fprintf(w, "Help\n")
	fmt.Fprintf(w, "http://"+xip+":"+port+"/help - returns this help\n")
	fmt.Fprintf(w, "http://"+xip+":"+port+"/ping - returns the remote IP address\n")

}
func handleRequests(port string) {
	http.Handle("/ping", http.HandlerFunc(getPing))
	http.Handle("/help", http.HandlerFunc(getHelp))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	port := GetOutboundPort()
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

func GetOutboundPort() string {
	port := "8080"
	return port
}
