package kg

import "strings"

type Quest struct {
	Id           string
	Number       int
	Chapter      string
	Title        string
	Subtitle     string
	TaskTitles   map[int]string
	Description  []string
	OriginalText string
}

// todo checking fields
// todo adding titles of tasks
func SnbtToQuest(s string, num int, chapter string) (Quest, error) {
	quest := Quest{
		Number:       num,
		Chapter:      chapter,
		TaskTitles:   map[int]string{},
		OriginalText: s,
	}

	isSavingDescription := false
	numTask := 0
	isTasks := false
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
