package kg

import (
	"fmt"
	"strings"
)

type Quest struct {
	Number       int
	ModpackName  string
	Chapter      string
	Id           string
	Title        string
	Subtitle     string
	TaskTitles   map[int]string
	Description  []string
	OriginalText string
}

// todo checking fields
func SnbtToQuest(num int, modpackName string, chapter string, originalText string) (Quest, error) {
	quest := Quest{
		Number:       num,
		ModpackName:  modpackName,
		Chapter:      chapter,
		TaskTitles:   map[int]string{},
		OriginalText: originalText,
	}

	isSavingDescription := false
	numTask := 0
	isTasks := false
	for _, line := range strings.Split(originalText, "\n") {

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
			continue
		}

		// task titles
		if strings.HasPrefix(line, "\t\t\ttasks:") {
			isTasks = true
		}
		if isTasks && strings.HasPrefix(line, "\t\t\t\t}") {
			numTask++
		}
		if after, ok := strings.CutPrefix(stripedLine, "title: "); ok && isTasks {
			taskTitle := strings.Trim(after, "\"")
			quest.TaskTitles[numTask] = taskTitle
		}
	}

	return quest, nil
}

func (q *Quest) GenerateKeys() string {
	result := ""

	isSavingDescription := false
	numDesctription := 0
	numArrayDescription := 0

	var keys []int
	numKey := 0
	for key := range q.TaskTitles {
		keys = append(keys, key)
	}

	for _, line := range strings.Split(q.OriginalText, "\n") {

		// title
		if strings.HasPrefix(line, "\t\t\ttitle: ") {
			result += fmt.Sprintf("\t\t\ttitle: \"{%s.%s.%s.quest%d.title}\"\n",
				q.ModpackName, q.Chapter, q.Id, q.Number)
			continue
		}

		// subtitle
		if strings.HasPrefix(line, "\t\t\tsubtitle: ") {
			result += fmt.Sprintf("\t\t\tsubtitle: \"{%s.%s.%s.quest%d.subtitle}\"\n",
				q.ModpackName, q.Chapter, q.Id, q.Number)
			continue
		}

		// description
		if strings.HasPrefix(line, "\t\t\tdescription: ") {
			if len(q.Description) == 1 {
				result += fmt.Sprintf("\t\t\tdescription: [\"{%s.%s.%s.quest%d.description0}\"]\n",
					q.ModpackName, q.Chapter, q.Id, q.Number)
				continue
			}
			result += "\t\t\tdescription: [\n"
			isSavingDescription = true
			continue
		}

		if isSavingDescription {
			if q.Description[numArrayDescription] == "" {
				result += "\t\t\t\t\"\"\n"
			} else {
				result += fmt.Sprintf("\t\t\t\t\"{%s.%s.%s.quest%d.description%d}\"\n",
					q.ModpackName, q.Chapter, q.Id, q.Number, numDesctription)
				numDesctription++
			}
			numArrayDescription++
			if numArrayDescription == len(q.Description) {
				isSavingDescription = false
			}
			continue
		}

		// task titles
		if strings.HasPrefix(line, "\t\t\t\t\ttitle: ") {
			result += fmt.Sprintf("\t\t\t\t\ttitle: \"{%s.%s.%s.quest%d.task%d.title}\"\n",
				q.ModpackName, q.Chapter, q.Id, q.Number, keys[numKey])
			numKey++
			continue
		}

		result += line + "\n"
	}
	return result
}

func (q *Quest) GenerateMapPart() (string, error) {
	numDesctription := 0
	result := ""

	// title
	if q.Title != "" {
		result += fmt.Sprintf(", \"%s.%s.%s.quest%d.title\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, q.Title)
	}

	// subtitle
	if q.Subtitle != "" {
		result += fmt.Sprintf(", \"%s.%s.%s.quest%d.subtitle\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, q.Subtitle)
	}

	// description
	for _, description := range q.Description {
		if description != "" {
			result += fmt.Sprintf(", \"%s.%s.%s.quest%d.description%d\": \"%s\"",
				q.ModpackName, q.Chapter, q.Id, q.Number, numDesctription, description)
			numDesctription++
		}
	}

	// task titles
	for i, titleTask := range q.TaskTitles {
		result += fmt.Sprintf(", \"%s.%s.%s.quest%d.task%d.title\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, i, titleTask)
	}

	if q.Number == 0 {
		result = strings.TrimPrefix(result, ", ")
	}

	return result, nil
}
