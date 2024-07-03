package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type Response struct {
	ClientIP  string `json:"client_ip"`
	Location  string `json:"location"`
	Greeting  string `json:"greeting"`
}

type IPAPIResponse struct {
	City string `json:"city"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/hello", helloHandler)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	visitorName := r.URL.Query().Get("visitor_name")
	if visitorName == "" {
		visitorName = "Guest"
	}

	clientIP := getClientIP(r)
	location, err := getLocation(clientIP)
	if err != nil {
		http.Error(w, "Error getting location", http.StatusInternalServerError)
		return
	}

	response := Response{
		ClientIP: clientIP,
		Location: location,
		Greeting: fmt.Sprintf("Hello, %s! The temperature is 11 degrees Celsius in %s", visitorName, location),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func getLocation(ip string) (string, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipAPIResponse IPAPIResponse
	err = json.Unmarshal(body, &ipAPIResponse)
	if err != nil {
		return "", err
	}

	return ipAPIResponse.City, nil
}