package ticket

import (
	"database/sql"
	"fmt"
)

type TicketPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *TicketPostgresRepository {
	return &TicketPostgresRepository{DB: db}
}

// func (repo *TicketPostgresRepository) GetAllFlights() ([]*Flight, error) {
// 	flights := make([]*Flight, 0)
// 	rows, err := repo.DB.Query("SELECT * FROM flight;")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute the query: %w", err)
// 	}

// 	for rows.Next() {
// 		f := &Flight{}
// 		if err := rows.Scan(&f.ID, &f.FlightNumber, &f.Date, &f.FromAirportId, &f.ToAirportId, &f.Price); err != nil {
// 			return nil, fmt.Errorf("failed to execute the query: %w", err)
// 		}
// 		flights = append(flights, f)
// 	}
// 	defer rows.Close()

// 	return flights, nil
// }

// func (repo *TicketPostgresRepository) GetFlightByNumber(flightNumber string) (*Flight, error) {
// 	flight := &Flight{}

// 	err := repo.DB.
// 		QueryRow("SELECT * FROM flight WHERE flight_number = $1;", flightNumber).
// 		Scan(
// 			&flight.ID,
// 			&flight.FlightNumber,
// 			&flight.Date,
// 			&flight.FromAirportId,
// 			&flight.ToAirportId,
// 			&flight.Price,
// 		)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return flight, err
// 		}
// 	}

// 	return flight, nil
// }

func (repo *TicketPostgresRepository) GetByUsername(username string) ([]*Ticket, error) {
	tickets := make([]*Ticket, 0)
	rows, err := repo.DB.Query("SELECT id, ticket_uid, username, flight_number, price, status FROM ticket WHERE username = $1;", username)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		ticket := new(Ticket)
		rows.Scan(
			&ticket.ID,
			&ticket.TicketUID,
			&ticket.Username,
			&ticket.FlightNumber,
			&ticket.Price,
			&ticket.Status)

		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %s", err)
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (repo *TicketPostgresRepository) Add(ticket *Ticket) error {
	_, err := repo.DB.Query(
		"INSERT INTO ticket (ticket_uid, username, flight_number, price, status) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		ticket.TicketUID,
		ticket.Username,
		ticket.FlightNumber,
		ticket.Price,
		ticket.Status,
	)

	return err
}
