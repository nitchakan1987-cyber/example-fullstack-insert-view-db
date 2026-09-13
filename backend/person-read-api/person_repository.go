package main

import (
	"context"
	"database/sql"
	"time"
)

type personRepository struct {
	db *sql.DB
}

func (repository *personRepository) list(ctx context.Context) ([]Person, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT Id, FirstName, LastName, BirthDate, Address
		FROM dbo.Persons
		ORDER BY Id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]Person, 0)
	for rows.Next() {
		person, err := scanPerson(rows)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (repository *personRepository) get(ctx context.Context, id int) (Person, error) {
	var row personRow
	err := repository.db.QueryRowContext(ctx, `
		SELECT Id, FirstName, LastName, BirthDate, Address
		FROM dbo.Persons
		WHERE Id = @Id
	`, sql.Named("Id", id)).Scan(
		&row.id,
		&row.firstName,
		&row.lastName,
		&row.birthDate,
		&row.address,
	)
	if err != nil {
		return Person{}, err
	}
	return row.person(), nil
}

type personRow struct {
	id        int
	firstName string
	lastName  string
	birthDate time.Time
	address   string
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPerson(scanner rowScanner) (Person, error) {
	var row personRow
	if err := scanner.Scan(
		&row.id,
		&row.firstName,
		&row.lastName,
		&row.birthDate,
		&row.address,
	); err != nil {
		return Person{}, err
	}
	return row.person(), nil
}

func (row personRow) person() Person {
	thailand := time.FixedZone("UTC+7", 7*60*60)
	return Person{
		ID:        row.id,
		FirstName: row.firstName,
		LastName:  row.lastName,
		BirthDate: row.birthDate.Format("2006-01-02"),
		Age:       time.Now().In(thailand).Year() - row.birthDate.Year(),
		Address:   row.address,
	}
}