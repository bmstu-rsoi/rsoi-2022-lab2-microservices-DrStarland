package handlers

import (
	"fmt"
	"gateway/pkg/models/airport"
)

func CalncelTicketController(ticketServiceAddress, bonusServiceAddress, username string) error {
	return nil
}

func (h *GatewayHandler) UserTicketsController(ticketServiceAddress, flightServiceAddress, username string) (*[]airport.TicketInfo, error) {
	tickets, err := GetUserTickets(ticketServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s\n", err)
	}

	ticketsInfo := make([]airport.TicketInfo, 0)
	for _, ticket := range *tickets {
		flight, err := GetFlight(flightServiceAddress, ticket.FlightNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to get flights: %s", err)
		}

		airportFrom, err := h.GetAirport(flightServiceAddress, flight.FromAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		airportTo, err := h.GetAirport(flightServiceAddress, flight.ToAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		ticketInfo := airport.TicketInfo{
			TicketUID:    ticket.TicketUID,
			FlightNumber: ticket.FlightNumber,
			FromAirport:  fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
			ToAirport:    fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
			Date:         flight.Date,
			Price:        ticket.Price,
			Status:       ticket.Status,
		}

		ticketsInfo = append(ticketsInfo, ticketInfo)
	}

	return &ticketsInfo, nil
}

func (h *GatewayHandler) UserInfoController(ticketServiceAddress, flightServiceAddress, bonusServiceAddress, username string) (*airport.UserInfo, error) {
	ticketsInfo, err := h.UserTicketsController(ticketServiceAddress, flightServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s", err)
	}

	privilege, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get privilege info: %s", err)
	}

	userInfo := &airport.UserInfo{
		TicketsInfo: ticketsInfo,
		Privilege: &airport.PrivilegeShortInfo{
			Status:  privilege.Status,
			Balance: privilege.Balance,
		},
	}

	return userInfo, nil
}

func UserPrivilegeController(bonusServiceAddress, username string) (*airport.PrivilegeInfo, error) {
	privilegeShortInfo, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s", err)
	}

	privilegeHistory, err := GetPrivilegeHistory(bonusServiceAddress, privilegeShortInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get privilege info: %s", err)
	}

	privilegeInfo := &airport.PrivilegeInfo{
		Status:  privilegeShortInfo.Status,
		Balance: privilegeShortInfo.Balance,
		History: privilegeHistory,
	}

	return privilegeInfo, nil
}

func (h *GatewayHandler) BuyTicketController(tAddr, fAddr, bAddr, username string, info *airport.BuyTicketInfo) (*airport.PurchaseTicketInfo, error) {
	flight, err := GetFlight(fAddr, info.FlightNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to get flights: %s", err)
	}

	airportFrom, err := h.GetAirport(fAddr, flight.FromAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	airportTo, err := h.GetAirport(fAddr, flight.ToAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	moneyPaid := flight.Price
	bonusesPaid := 0
	diff := int(float32(info.Price) * 0.1)
	optype := "FILL_IN_BALANCE"

	if info.PaidFromBalance {
		if info.Price > 0 {
			bonusesPaid = 0
		} else {
			bonusesPaid = info.Price
		}

		moneyPaid = moneyPaid - bonusesPaid
		diff = -bonusesPaid
		optype = "DEBIT_THE_ACCOUNT"
	}

	uid, err := CreateTicket(tAddr, username, info.FlightNumber, flight.Price)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ticket: %s", err)
	}

	if !info.PaidFromBalance {
		if err := CreatePrivilege(bAddr, username, diff); err != nil {
			return nil, fmt.Errorf("Failed to get privilege info: %s", err)
		}
	}

	err = CreatePrivilegeHistoryRecord(bAddr, uid, flight.Date, optype, 1, diff)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bonus history record: %s", err)
	}

	purchaseInfo := airport.PurchaseTicketInfo{
		TicketUID:     uid,
		FlightNumber:  info.FlightNumber,
		FromAirport:   fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
		ToAirport:     fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
		Date:          flight.Date,
		Price:         flight.Price,
		PaidByMoney:   moneyPaid,
		PaidByBonuses: bonusesPaid,
		Status:        "PAID",
		Privilege: &airport.PrivilegeShortInfo{
			Balance: diff,
			Status:  "GOLD",
		},
	}

	return &purchaseInfo, nil
}
