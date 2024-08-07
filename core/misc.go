package core

import (
	"encoding/json"
	"fmt"
	"os"

	"go-abbreviation/models"
)

func InsertFromJson() {
	jsonFile, err := os.Open("static/data.json")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Successfully opened json file")
	defer jsonFile.Close()

	// Decode json file
	var list models.List
	json.NewDecoder(jsonFile).Decode(&list)
	for i := 0; i < len(list.List); i++ {
		fmt.Println(list.List[i].Short, ".....", list.List[i].Long)
	}
}

// Defined for slice only
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
