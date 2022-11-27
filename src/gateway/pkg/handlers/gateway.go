package handlers

import (
	"encoding/json"
	"gateway/pkg/models/airport"
	"gateway/pkg/myjson"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type GatewayHandler struct {
	TicketServiceAddress string
	FlightServiceAddress string
	BonusServiceAddress  string
	Logger               *zap.SugaredLogger
}

func (h *GatewayHandler) GetAllFlights(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	params := r.URL.Query()

	flights, err := h.GetAllFlightsInfo(h.FlightServiceAddress)
	if err != nil {
		h.Logger.Errorln("failed to get response from flighst service: " + err.Error())
		myjson.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pageParam := params.Get("page")
	if pageParam == "" {
		log.Println("invalid query parameter: (page)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		log.Printf("unable to convert the string into int: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		log.Println("invalid query parameter (size)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		log.Printf("unable to convert the string into int:  %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	right := page * size
	if len(*flights) < right {
		right = len(*flights)
	}

	flightsStripped := (*flights)[(page-1)*size : right]
	cars := airport.FlightsLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(flightsStripped),
		Items:         &flightsStripped,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		log.Printf("failed to encode response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) GetUserTickets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketsInfo, err := h.UserTicketsController(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketsInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) CancelTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := CalncelTicketController(
		h.TicketServiceAddress,
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GatewayHandler) GetUserTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketUID := ps.ByName("ticketUid")

	ticketsInfo, err := h.UserTicketsController(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ticketInfo *airport.TicketInfo
	for _, ticket := range *ticketsInfo {
		if ticket.TicketUID == ticketUID {
			ticketInfo = &ticket
		}
	}

	if ticketInfo == nil {
		log.Printf("Ticket not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) BuyTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ticketInfo airport.BuyTicketInfo
	err := json.NewDecoder(r.Body).Decode(&ticketInfo)
	if err != nil {
		log.Printf("failed to decode post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tickets, err := h.BuyTicketController(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		h.BonusServiceAddress,
		username,
		&ticketInfo,
	)

	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tickets)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, err := h.UserInfoController(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) GetPrivilege(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privilegeInfo, err := UserPrivilegeController(
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(privilegeInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
