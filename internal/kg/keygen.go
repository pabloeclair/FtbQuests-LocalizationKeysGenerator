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
// todo adding titles of tasks
func (q *Quest) GenerateKeys() string {
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
func (q *Quest) GenerateMap(lang Lang) string {
	return ""
}
