package privilege

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (repo *PostgresRepository) CreateHistoryRecord(record *PrivilegeHistory) error {
	_, err := repo.DB.Query(
		"INSERT INTO privilege_history (privilege_id, ticket_uid, datetime, balance_diff, operation_type) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		record.PrivilegeID,
		record.TicketUID,
		record.Date,
		record.BalanceDiff,
		record.OperationType,
	)

	return err
}

func (repo *PostgresRepository) CreatePrivilege(record *Privilege) error {
	_, err := repo.DB.Query(
		"INSERT INTO privilege (username, balance) VALUES ($1, $2) RETURNING id;",
		record.Username,
		record.Balance,
	)
	return err
}

func (repo *PostgresRepository) GetPrvilegeByUsername(username string) (*Privilege, error) {
	var privilege Privilege

	log.Printf(">>>> username: %s", username)
	row := repo.DB.QueryRow("SELECT * FROM privilege where username = $1;", username)
	err := row.Scan(&privilege.ID, &privilege.Username, &privilege.Status, &privilege.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &privilege, err
		}
	}

	return &privilege, nil
}

func (repo *PostgresRepository) GetHistoryById(privilegeID string) ([]*PrivilegeHistory, error) {
	var history []*PrivilegeHistory
	rows, err := repo.DB.Query("SELECT * FROM privilege_history where privilege_id = $1;", privilegeID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %s", err)
	}

	for rows.Next() {
		row := new(PrivilegeHistory)
		rows.Scan(
			&row.ID,
			&row.PrivilegeID,
			&row.TicketUID,
			&row.Date,
			&row.BalanceDiff,
			&row.OperationType,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %s", err)
		}

		history = append(history, row)
	}

	return history, nil
}

// const createHistory = `INSERT INTO privilege_history (privilege_id, ticket_uid, datetime, balance_diff, operation_type) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

// func (repo *PostgresRepository) CreateHistoryRecord(record *PrivilegeHistory) error {

// 	_, err := repo.DB.Query(
// 		createHistory,
// 		record.PrivilegeID,
// 		record.TicketUID,
// 		record.Date,
// 		record.BalanceDiff,
// 		record.OperationType,
// 	)

// 	return err
// }

// const createPrivilege = `INSERT INTO privilege (username, balance) VALUES ($1, $2) RETURNING id;`

// func (repo *PostgresRepository) CreatePrivilege(record *Privilege) error {
// 	repo.DB = db.CreateConnection()
// 	defer repo.DB.Close()

// 	_, err := repo.DB.Query(
// 		createPrivilege,
// 		record.Username,
// 		record.Balance,
// 	)

// 	return err
// }

// const selectPrivilegeByUsername = `SELECT * FROM privilege where username = $1;`

// func (repo *PostgresRepository) GetPrvilegeByUsername(username string) (*Privilege, error) {
// 	repo.DB = db.CreateConnection()
// 	defer repo.DB.Close()

// 	var privilege Privilege

// 	log.Printf(">>>> username: %s", username)
// 	row := repo.DB.QueryRow(selectPrivilegeByUsername, username)
// 	err := row.Scan(&privilege.ID, &privilege.Username, &privilege.Status, &privilege.Balance)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &privilege, err
// 		}
// 	}

// 	return &privilege, nil
// }

// const selectHistoryByID = `SELECT * FROM privilege_history where privilege_id = $1;`

// func (repo *PostgresRepository) GetHistoryById(privilegeID string) ([]*PrivilegeHistory, error) {
// 	repo.DB = db.CreateConnection()
// 	defer repo.DB.Close()

// 	var history []*PrivilegeHistory
// 	rows, err := repo.DB.Query(selectHistoryByID, privilegeID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute the query: %w", err)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to execute the query: %s", err)
// 	}

// 	for rows.Next() {
// 		row := new(PrivilegeHistory)
// 		rows.Scan(
// 			&row.ID,
// 			&row.PrivilegeID,
// 			&row.TicketUID,
// 			&row.Date,
// 			&row.BalanceDiff,
// 			&row.OperationType,
// 		)

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to execute the query: %s", err)
// 		}

// 		history = append(history, row)
// 	}

// 	return history, nil
// }
