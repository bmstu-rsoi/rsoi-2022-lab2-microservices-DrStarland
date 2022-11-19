package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/database"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/handlers"
	mid "github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/middleware"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/models/airport"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/models/flight"
	"go.uber.org/zap"
)

func HealthOK(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	db, err := database.CreateConnection()
	if err != nil {
		log.Panicln(err.Error())
	}
	defer db.Close()

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	repoFlight := flight.NewPostgresRepo(db)
	repoAirport := airport.NewPostgresRepo(db)

	allHandler := &handlers.FlightsHandler{
		Logger:      logger,
		FlightRepo:  repoFlight,
		AirportRepo: repoAirport,
	}

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok == true {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	router.GET("/api/v1/flights", mid.AccessLog(allHandler.GetAllFlight, logger))
	router.GET("/api/v1/flight/:flightNumber", mid.AccessLog(allHandler.GetFlight, logger))
	router.GET("/api/v1/airport/:airportID", mid.AccessLog(allHandler.GetAirport, logger))
	router.GET("/manage/health", HealthOK)

	ServerAddress := ":8080" // os.Getenv("PORT")

	logger.Infow("starting server",
		"type", "START",
		"addr", ServerAddress,
	)
	err = http.ListenAndServe(ServerAddress, router)
	if err != nil {
		log.Panicln(err.Error())
	}
}
