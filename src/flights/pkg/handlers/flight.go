package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/models/airport"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/models/flight"
	"github.com/rsoi-2022-lab2-microservices-DrStarland/src/flights/pkg/myjson"
	"go.uber.org/zap"
)

type FlightsHandler struct {
	Logger      *zap.SugaredLogger
	FlightRepo  flight.Repository
	AirportRepo airport.Repository
}

func (h *FlightsHandler) GetAllFlight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flight, err := h.FlightRepo.GetAllFlights()
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(flight); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FlightsHandler) GetFlight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	number := ps.ByName("flightNumber")
	log.Println(number)
	flight, err := h.FlightRepo.GetFlightByNumber(number)

	if err != nil {
		log.Printf("Failed to get flights: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(flight); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FlightsHandler) GetAirport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("airportID")
	log.Println(id)
	airport, err := h.AirportRepo.GetAirportByID(id)
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	myjson.JsonResponce(w, http.StatusOK, airport)
}
