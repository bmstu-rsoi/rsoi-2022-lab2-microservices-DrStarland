package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTicketsByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	ticketRepo := repository.TicketRepository{}
	params := mux.Vars(r)
	tickets, err := ticketRepo.GetTicketsByUsername(params["username"])
	if err != nil {
		log.Printf("Failed to get ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tickets); err != nil {
		log.Printf("Failed to encode ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func BuyTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket

	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketRepo := repository.TicketRepository{}
	if err := ticketRepo.CreateTicket(&ticket); err != nil {
		log.Printf("Failed to create ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
