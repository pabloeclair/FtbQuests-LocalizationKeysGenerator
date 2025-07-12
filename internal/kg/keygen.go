package kg

import (
	"fmt"
	"strings"
)

// type Quest struct {
// 		Number       int
// 		ModpackName  string
// 		Chapter      string
// 		Id           string
// 		Title        string
// 		Subtitle     string
// 		TaskTitles   map[int]string
// 		Description  []string
// 		OriginalText string
// }

// todo checking fields
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

func (q *Quest) GenerateMapPart(l string) (string, error) {
	lang := Lang_array[l]
	if lang.String() == "unknown" {
		return "", fmt.Errorf("unknown language %s; please, check correct codes in internal/kg/lang.go", l)
	}

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
