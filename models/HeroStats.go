package models

import (
	"time"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
	"strconv"
)

const (
	DB_USER = ""
	DB_PASSWORD = ""
	DB_NAME = "overwatch"
	UPDATE_INTERVAL = 24
)

var databaseInfo = fmt.Sprintf(
	"user=%s password=%s dbname=%s sslmode=disable",
	DB_USER,
	DB_PASSWORD,
	DB_NAME,
)

type HeroStats struct {
	Id int
	HeroId int
	Eliminations int `json:",string,omitempty"`
	Deaths int `json:",string,omitempty"`
	WeaponAccuracy string
	EliminationsAverage int `json:"Eliminations-Average"`
	TimePlayed string
	GamesPlayed int `json:",string,omitempty"`
	GamesWon int `json:",string,omitempty"`
	WinPercentage string
	UpdateTime time.Time
}

// Function should return a true if the last update in the db is over 24 hours ago, or if there
// aren't any records in the db.
func ShouldUpdate() (bool, error) {
	db, err := sql.Open("postgres", databaseInfo)

	if err != nil {
		return false, err
	}

	var count int

	// Check count of records in the db first to prevent nil from being stored in a time field
	db.QueryRow("SELECT COUNT(id) FROM hero_stats").Scan(&count)

	if count == 0 {
		return true, nil
	}

	rows, err := db.Query("SELECT MAX(update_time) FROM hero_stats")

	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		stats := new(HeroStats)
		err = rows.Scan(&stats.UpdateTime)


		if err != nil {
			return false, err
		}

		// This time is not accurate because the db stores time without a time zone
		// but, it's close enough for now
		timeSince := time.Since(stats.UpdateTime).Hours()

		// If it's been less than the specified update interval, do not update
		if timeSince < UPDATE_INTERVAL {
			return false, nil
		}
	}

	return true, nil
}

func (h HeroStats) Save() (err error) {
	db, err := sql.Open("postgres", databaseInfo)

	if err != nil {
		return err
	}

	defer db.Close()

	fmt.Println(h)
	fmt.Println("Inserting into db")
	fmt.Println(h.Deaths)

	err = db.QueryRow(
		"INSERT INTO hero_stats(hero_id, eliminations, deaths, weapon_accuracy, eliminations_average, time_played, games_played, games_won, win_percentage, update_time)" +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;",
		h.HeroId,
		h.Eliminations,
		h.Deaths,
		h.WeaponAccuracy,
		h.EliminationsAverage,
		h.TimePlayed,
		h.GamesPlayed,
		h.GamesWon,
		h.WinPercentage,
		time.Now(),
	).Scan(&h.Id)

	if err != nil {
		return err
	}

	return nil
}

// Hackish Method that maps deserialized json to the stats. This was needed because the numbers returned
// from the api include commas which prevent clean deserialization to a struct with ints
func (h *HeroStats) Decode(dirty map[string]string) (err error) {
	if deaths := dirty["Deaths"]; len(deaths) > 0 {
		h.Deaths, err = strconv.Atoi(strings.Replace(deaths, ",", "", -1))
	}
	if eliminations := dirty["Eliminations"]; len(eliminations) > 0 {
		h.Eliminations, err = strconv.Atoi(strings.Replace(eliminations, ",", "", -1))
	}
	if gamesPlayed := dirty["GamesPlayed"]; len(gamesPlayed) > 0 {
		h.GamesPlayed, err = strconv.Atoi(strings.Replace(gamesPlayed, ",", "", -1))
	}
	if gamesWon := dirty["GamesWon"]; len(gamesWon) > 0 {
		h.GamesWon, err = strconv.Atoi(strings.Replace(gamesWon, ",", "", -1))
	}
	if eliminationsAvg := dirty["Eliminations-Average"]; len(eliminationsAvg) > 0 {
		// TODO: Maybe make this a float?
		eliminationsAvg = strings.Split(eliminationsAvg, ".")[0]

		h.EliminationsAverage, err = strconv.Atoi(strings.Replace(eliminationsAvg, ",", "", -1))
	}
	h.TimePlayed = dirty["TimePlayed"]
	h.WinPercentage = dirty["WinPercentage"]
	h.WeaponAccuracy = dirty["WeaponAccuracy"]

	return
}