package main

import (
	"github.com/piotrromanowski/overwatch_stats_grapher/models"
	"fmt"
	"github.com/piotrromanowski/overwatch_stats_grapher/apis"
)

func main() {
	//apis.GetAllHeroes()

	update, err := models.ShouldUpdate()

	if err != nil {
		fmt.Println(err.Error())
	}
	if update {
		fmt.Println("Should Update")
		err := apis.UpdateHeros()

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}