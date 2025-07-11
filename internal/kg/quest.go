package kg

import "strings"

type Lang string

const (
	Ru Lang = "ru_ru"
	En Lang = "en_us"
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
