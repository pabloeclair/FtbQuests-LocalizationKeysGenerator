package kg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func GenerateQuestsAndKeys(modpackName string, fileName string) ([]*Quest, string, error) {
	resultQuests := []*Quest{}
	resultKeys := ""

	filePath := filepath.Join("ftbquests", "quests", "chapters", fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, resultKeys, fmt.Errorf("opening file %s error: %v", filePath, err)
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
				return nil, resultKeys, fmt.Errorf("parsing quest %s error: %v", questLines, err)
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
		return nil, resultKeys, fmt.Errorf("reading file %s error: %v", filePath, err)
	}

	return resultQuests, resultKeys, nil
}

func GenerateMap(lang string, questsMap map[string][]*Quest) (string, error) {
	l := Lang_array[lang]
	if l.String() == "unknown" {
		return "", fmt.Errorf("unknown language %s; please, check correct codes in internal/kg/lang.go", l)
	}

	keys := make([]string, 0, len(questsMap))
	for k := range questsMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	result := "{\n"

	i := 0
	for _, key := range keys {
		for ii, quest := range questsMap[key] {
			if ii == len(questsMap[key])-1 && i == len(questsMap)-1 {
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
