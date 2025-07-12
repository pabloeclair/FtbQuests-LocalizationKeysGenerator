package kg

import (
	"fmt"
	"strings"
)

type Lang string

const (
	Ru Lang = "ru_ru"
	En Lang = "en_us"
)

// todo checking fields
func (q Quest) GenerateKeys() string {
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
			if q.Description[numArrayDescription] == "" {
				result += "\t\t\t\t\"\"\n"
			} else {
				result += fmt.Sprintf("\t\t\t\t\"{homestead.%s.%s.quest%d.description%d}\"\n",
					q.Chapter, q.Id, q.Number, numDesctription)
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
			result += fmt.Sprintf("\t\t\t\t\ttitle: \"{homestead.%s.%s.quest%d.task%d.title}\"\n",
				q.Chapter, q.Id, q.Number, keys[numKey])
			numKey++
			continue
		}

		result += line + "\n"
	}
	return result
}

// todo
func (q *Quest) GenerateMap(lang Lang) string {
	return ""
}
