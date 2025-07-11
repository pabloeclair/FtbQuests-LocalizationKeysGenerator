package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Quest struct {
	Id           string
	Number       int
	Chapter      string
	Title        string
	Subtitle     string
	TaskTitles   []string
	Description  []string
	OriginalText string
}

// todo checking fields
// todo adding titles of tasksgit
func (q *Quest) CreateKeys() string {
	result := ""

	isSavingDescription := false
	numDesctription := 0
	for _, line := range strings.Split(q.OriginalText, "\n") {

		// title
		if strings.HasPrefix(line, "\t\t\ttitle: ") {
			result += fmt.Sprintf("\t\t\ttitle: \"{homestead.%s.%s.quest%d.title}\"\n",
				q.Chapter, q.Id, q.Number)
			continue
		}

		// subtitle
		if strings.HasPrefix(line, "\t\t\tsubtitle: ") {
			result += fmt.Sprintf("\t\t\tsubtitle: \"{homestead.%s.%s.quest%d.subtitle}\"\n",
				q.Chapter, q.Id, q.Number)
			continue
		}

		// description
		if strings.HasPrefix(line, "\t\t\tdescription: ") {
			if len(q.Description) == 1 {
				result += fmt.Sprintf("\t\t\tdescription: [\"{homestead.%s.%s.quest%d.description0}\"]\n",
					q.Chapter, q.Id, q.Number)
				continue
			}
			result += "\t\t\tdescription: [\n"
			isSavingDescription = true
			continue
		}

		if isSavingDescription {
			if q.Description[numDesctription] == "" {
				result += "\t\t\t\t\"\"\n"
			} else {
				result += fmt.Sprintf("\t\t\t\t\"{homestead.%s.%s.quest%d.description%d}\"\n",
					q.Chapter, q.Id, q.Number, numDesctription)
			}
			numDesctription++
			if numDesctription == len(q.Description) {
				isSavingDescription = false
			}
			continue
		}

		result += line + "\n"
	}
	return result
}

// todo
func (q *Quest) CreateEnFields() string {
	return ""
}

// todo checking fields
// todo adding titles of tasks
func SnbtToQuest(s string, num int, chapter string) (Quest, error) {
	quest := Quest{Number: num, Chapter: chapter, OriginalText: s}

	isSavingDescription := false
	for _, line := range strings.Split(s, "\n") {

		// id
		if after, ok := strings.CutPrefix(line, "\t\t\tid: "); ok {
			id := strings.Trim(after, "\"")
			quest.Id = id
		}

		// title
		if after, ok := strings.CutPrefix(line, "\t\t\ttitle: "); ok {
			title := strings.Trim(after, "\"")
			quest.Title = title
		}

		// subtitle
		if after, ok := strings.CutPrefix(line, "\t\t\tsubtitle: "); ok {
			subtitle := strings.Trim(after, "\"")
			quest.Subtitle = subtitle
		}

		// description
		stripedLine := strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "\t\t\tdescription: "); ok {
			if strings.HasSuffix(stripedLine, "]") {
				description := strings.Trim(after, "[]\"")
				quest.Description = append(quest.Description, description)
				continue
			}
			isSavingDescription = true
			continue
		}

		if isSavingDescription {
			if stripedLine == "]" {
				isSavingDescription = false
				continue
			}

			descriptionPart := strings.Trim(stripedLine, "\"")
			quest.Description = append(quest.Description, descriptionPart)
		}
	}

	return quest, nil
}

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
	questArray := []Quest{}

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
			quest, err := SnbtToQuest(questLines, len(questArray), chapter)
			if err != nil {
				log.Fatal("Error parsing quest:", err)
			}
			result += quest.CreateKeys()
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
