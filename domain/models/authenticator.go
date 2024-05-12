package models

import (
	"encoding/json"
	"log"
	"os"
)

var Authenticators map[string]Authenticator

type Authenticator struct {
	Name      string `json:"name"`
	IconDark  string `json:"icon_dark"`
	IconLight string `json:"icon_light"`
}

func init() {
	data, err := os.ReadFile("aaguids.json")
	if err != nil {
		log.Fatalf("Failed to read authenticators file: %v", err)
	}

	err = json.Unmarshal(data, &Authenticators)
	if err != nil {
		log.Fatalf("Failed to parse authenticators file: %v", err)
	}
}
