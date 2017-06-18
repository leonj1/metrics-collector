package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"metrics-collector/models"
	"metrics-collector/routes"
	"net/http"
)

func main() {
	var userName = flag.String("user", "", "Database username")
	var password = flag.String("pass", "", "Database password")
	var databaseHost = flag.String("dbHost", "", "Database host")
	var databaseName = flag.String("db", "", "Database name")
	var databasePort = flag.Int("dbPort", 3306, "Database port")
	var serverPort = flag.String("port", "", "Web Server port")
	flag.Parse()

	// open connection to db
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", *userName, *password, *databaseHost, *databasePort, *databaseName)
	models.InitDB(connectionString)

	router := httprouter.New()

	router.POST("/public/metrics", routes.AddMetric)
	router.GET("/public/metrics/:host/:name/:days", routes.GetMetrics)
	router.GET("/public/hosts", routes.GetHosts)

	// Add CORS support (Cross Origin Resource Sharing)
	handler := cors.Default().Handler(router)

	log.Println("Starting web server")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), handler))
}
