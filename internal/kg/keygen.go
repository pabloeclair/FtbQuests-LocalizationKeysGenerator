package kg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseQuestsAndGenerateKeys(modpackName string, fileName string) ([]Quest, string, error) {
	resultQuests := []Quest{}
	resultKeys := ""

	filePath := filepath.Join("ftbquests", "quests", "chapters", fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return resultQuests, resultKeys, fmt.Errorf("opening file %s error: %v", filePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	chapter := strings.TrimSuffix(fileName, ".snbt")

	startArray := false
	startQuest := false
	questLines := ""

	for scanner.Scan() {

		// Waiting for the quests array to start
		if !startArray {
			resultKeys += scanner.Text() + "\n"
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
			quest, err := SnbtToQuest(len(resultQuests), modpackName, chapter, questLines)
			if err != nil {
				return resultQuests, resultKeys, fmt.Errorf("parsing quest %s error: %v", questLines, err)
			}
			resultKeys += quest.GenerateKeys()
			questLines = ""
			resultQuests = append(resultQuests, quest)
			startQuest = false
			continue
		}

		if startArray && !startQuest {
			resultKeys += scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		return resultQuests, resultKeys, fmt.Errorf("reading file %s error: %v", filePath, err)
	}

	return resultQuests, resultKeys, nil
}

func GenerateMap(lang string, questsMap map[string][]Quest) (string, error) {
	l := Lang_array[lang]
	if l.String() == "unknown" {
		return "", fmt.Errorf("unknown language %s; please, check correct codes in internal/kg/lang.go", l)
	}

	result := "{\n"

	i := 0
	for _, quests := range questsMap {
		for ii, quest := range quests {
			if ii == len(quests)-1 && i == len(questsMap)-1 {
				result += quest.GenerateMapPart() + "\n"
			} else {
				result += quest.GenerateMapPart() + ",\n"
			}
		}
		i++
	}
	result += "}\n"

	return result, nil
}
