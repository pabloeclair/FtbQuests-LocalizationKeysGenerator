package main

import (
	"bufio"
	"errors"
	"keysgen/internal/kg"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("beginning_game.snbt")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	chapter := strings.TrimSuffix(f.Name(), ".snbt")

	scanner := bufio.NewScanner(f)

	startArray := false
	startQuest := false
	result := ""
	questLines := ""
	questArray := []kg.Quest{}

	for scanner.Scan() {

		// Waiting for the quests array to start
		if !startArray {
			result += scanner.Text() + "\n"
			if strings.TrimSpace(scanner.Text()) == "quests: [" {
				startArray = true
				continue
			} else {
				continue
			}
		}

		// Start reading quest
		if scanner.Text() == "\t\t{" {
			questLines = scanner.Text() + "\n"
			startQuest = true
			continue
		}

		if startQuest {
			questLines += scanner.Text() + "\n"
		}

		// End reading quest
		if scanner.Text() == "\t\t}" {

			// Creating keys
			quest, err := kg.SnbtToQuest(len(questArray), "homestead", chapter, questLines)
			if err != nil {
				log.Fatal("Error parsing quest:", err)
			}
			result += quest.GenerateKeys()
			questLines = ""
			questArray = append(questArray, quest)
			startQuest = false
			continue
		}

		if startArray && !startQuest {
			result += scanner.Text() + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}

	if err = os.Mkdir("output", 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("Error creating output directory:", err)
		}
	}

	newF, err := os.Create("output/beginning_game.snbt")
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("Error creating output file:", err)
		}
	}

	if _, err = newF.WriteString(result); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("Error writing to output file:", err)
		}
	}
}
