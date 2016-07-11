package models

import (
	"fmt"
	"database/sql"
)

type Hero struct {
	Id int
	Name string `json:"name"`
	Image string `json:"image"`
}

// Save Hero. If saving fails at any point return the error.
func (h *Hero) Save() error{
	db, err := sql.Open("postgres", databaseInfo)

	if err != nil {
		return err
	}

	defer db.Close()

	err = db.QueryRow("INSERT INTO hero(name, image) VALUES($1, $2) returning id;", h.Name, h.Image).Scan(&h.Id)

	if err != nil {
		return err
	}

	return nil
}

func FindByName(name string) (id int, err error) {
	query := fmt.Sprintf("SELECT * FROM hero where name = '%s'", name)

	db, err := sql.Open("postgres", databaseInfo)

	if err != nil {
		return 0, err
	}

	defer db.Close()

	rows, err := db.Query(query)

	if err != nil {
		return 0, err
	}

	var hero = new(Hero)

	for rows.Next() {
		err := rows.Scan(&hero.Id, &hero.Name, &hero.Image)

		if err != nil {
			return 0, err
		}

		return hero.Id, nil
	}

	return 0, nil
}