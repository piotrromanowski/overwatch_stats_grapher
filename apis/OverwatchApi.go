package apis

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/piotrromanowski/overwatch_stats_grapher/models"
)

func UpdateHeros() error {
	resp, err := http.Get("https://api.lootbox.eu/pc/us/PiotrJS-1914/quick-play/heroes")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var heroes []models.Hero

	if err := json.NewDecoder(resp.Body).Decode(&heroes); err != nil {
		return err
	}

	var heroesAlternateNames =  map[string]string {
		"Soldier: 76" : "Soldier76",
		"Torbj&#xF6;rn" : "Torbjoern",
		"L&#xFA;cio" : "Lucio",
	}

	// get top five heroes
	heroes = heroes[:5]

	for _, hero := range heroes {
		// Swap out any hero alternate name
		if name := heroesAlternateNames[hero.Name]; len(name) != 0 {
			hero.Name = name
		}

		// Check to make sure record exists in db, if not, add it
		hero.Id, err = models.FindByName(hero.Name)

		if err != nil {
			return err
		}

		// If hero does not exist, save him/her to the db.
		if hero.Id == 0 {
			err := hero.Save()

			if err != nil {
				return err
			}
		}

		fmt.Println(hero)

		err = updateStats(hero)

		if err != nil {
			return err
		}
	}
	return nil
}

func updateStats(hero models.Hero) error {
	url := fmt.Sprintf("https://api.lootbox.eu/pc/us/PiotrJS-1914/quick-play/hero/%s/", hero.Name)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var dirtyStats map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&dirtyStats); err != nil {
		return err
	}

	fmt.Printf("Dirty Stats: %s \n", dirtyStats)

	var stats models.HeroStats

	err = stats.Decode(dirtyStats)

	if err != nil {
		return err
	}

	stats.HeroId = hero.Id

	err = stats.Save()

	if err != nil {
		return err
		fmt.Printf("Error Saving Hero %s \n", hero.Name)
	}

	return nil
}
