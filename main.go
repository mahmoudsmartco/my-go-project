package main

import (
	"app2_http_api_database/cache"
	"app2_http_api_database/config"
	"app2_http_api_database/middleware"
	"app2_http_api_database/routes"
	"app2_http_api_database/service/rabbitmq"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// existing init...
	// Init RabbitMQ publisher
	rabbitmq.InitDefaultPublisher("amqp://guest:guest@localhost:5672/", "students.exchange")
	defer rabbitmq.CloseDefaultPublisher()

	// Initialize DB & Redis
	config.InitDB()
	cache.InitRedis()
	defer config.CloseDB()

	// Setup routes
	mux := routes.SetupRoutes()

	// Logger middleware
	loggedMux := middleware.Logger(mux)

	// --------------------------
	// 1️⃣ HTTP Server
	// --------------------------
	go func() {
		fmt.Println("⚡ HTTP Server running on http://localhost:8080")
		err := http.ListenAndServe(":8080", loggedMux)
		if err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// --------------------------
	// 2️⃣ HTTPS Server
	// --------------------------
	fmt.Println("🚀 HTTPS Server running on https://localhost:8443")
	err := http.ListenAndServeTLS(
		":8443",
		"certs/server.crt",
		"certs/server.key",
		loggedMux,
	)
	if err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
