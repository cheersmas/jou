package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cheersmas/jou/domains"
)

const createJournalTable string = `
	CREATE TABLE IF NOT EXISTS journals (
	id INTEGER NOT NULL PRIMARY KEY,
	content TEXT NOT NULL,
	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`

type journalRepository struct {
	db *sql.DB

	// queries
	readJournalQuery    *sql.Stmt
	existsJournalQuery  *sql.Stmt
	insertJournalQuery  *sql.Stmt
	deleteJournalQuery  *sql.Stmt
	updateJournalQuery  *sql.Stmt
	listAllJournalQuery *sql.Stmt
}

func (jr *journalRepository) Create(ctx context.Context, content domains.Journal) (int, error) {
	// Use Go's time.Now() to ensure consistent timezone handling
	now := time.Now()
	res, err := jr.insertJournalQuery.ExecContext(ctx, content.Content, now)
	if err != nil {
		log.Printf("ERROR: failed to create a journal entry: %v", err)
		return -1, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, err
	}
	return int(id), nil
}

func (jr *journalRepository) Read(ctx context.Context, journalId int) (domains.Journal, error) {
	var journal domains.Journal
	if err := jr.readJournalQuery.QueryRowContext(ctx, journalId).Scan(&journal.Id, &journal.Content, &journal.CreatedAt); err == sql.ErrNoRows {
		return journal, err
	}
	return journal, nil
}

func (jr *journalRepository) ListAll(ctx context.Context) ([]domains.Journal, error) {
	rows, err := jr.listAllJournalQuery.QueryContext(ctx)
	if err != nil {
		log.Printf("ERROR: failed to query journals: %v", err)
		return nil, err
	}
	defer rows.Close()

	var journals []domains.Journal
	for rows.Next() {
		var journal domains.Journal
		if err := rows.Scan(&journal.Id, &journal.Content, &journal.CreatedAt); err != nil {
			log.Printf("ERROR: failed to scan journal row: %v", err)
			return nil, err
		}
		journals = append(journals, journal)
	}

	if err = rows.Err(); err != nil {
		log.Printf("ERROR: error after scanning rows: %v", err)
		return nil, err
	}

	return journals, nil
}

func (jr *journalRepository) Update(ctx context.Context, id int, content string) (int, error) {
	res, err := jr.updateJournalQuery.ExecContext(ctx, content, id)
	if err != nil {
		return -1, nil
	}
	rowsEffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowsEffected == 0 {
		return -1, fmt.Errorf("no journal found with id %d", id)
	}

	return int(id), nil
}

func (jr *journalRepository) Delete(ctx context.Context, id int) (int, error) {
	_, err := jr.deleteJournalQuery.ExecContext(ctx, id)
	if err != nil {
		return -1, err
	}
	return id, err
}

func NewJournalRepository(ctx context.Context, db *sql.DB) (*journalRepository, error) {
	// Create tables if they don't exist
	if _, err := db.ExecContext(ctx, createJournalTable); err != nil {
		log.Printf("%q %s\n", err, createJournalTable)
		return nil, err
	}

	readJournalQuery, err := db.PrepareContext(ctx, "SELECT id, content, createdAt FROM journals WHERE id = ?")
	if err != nil {
		return nil, err
	}
	existsJournalQuery, err := db.PrepareContext(ctx, "SELECT EXISTS(SELECT 1 FROM journals WHERE id = ?)")
	if err != nil {
		return nil, err
	}
	// Updated to include createdAt parameter
	insertJournalQuery, err := db.PrepareContext(ctx, "INSERT INTO journals(content, createdAt) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	deleteJournalQuery, err := db.PrepareContext(ctx, "DELETE FROM journals WHERE id = ?")
	if err != nil {
		return nil, err
	}
	updateJournalQuery, err := db.PrepareContext(ctx, "UPDATE journals SET content = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	listAllJournalQuery, err := db.PrepareContext(ctx, "SELECT id, content, createdAt FROM journals ORDER BY createdAt DESC")
	if err != nil {
		return nil, err
	}

	return &journalRepository{
		db:                  db,
		readJournalQuery:    readJournalQuery,
		existsJournalQuery:  existsJournalQuery,
		insertJournalQuery:  insertJournalQuery,
		deleteJournalQuery:  deleteJournalQuery,
		updateJournalQuery:  updateJournalQuery,
		listAllJournalQuery: listAllJournalQuery,
	}, nil
}
