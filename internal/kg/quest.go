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
	RewardTitles map[int]string
	Description  []string
	OriginalText string
}

// todo checking fields
func SnbtToQuest(num int, modpackName string, chapter string, originalText string) (*Quest, error) {
	quest := Quest{
		Number:       num,
		ModpackName:  modpackName,
		Chapter:      chapter,
		TaskTitles:   map[int]string{},
		RewardTitles: map[int]string{},
		OriginalText: originalText,
	}

	isSavingDescription := false
	numTask := 0
	numReward := 0
	isTasks := false
	isReward := false
	for _, line := range strings.Split(originalText, "\n") {

		// id
		if after, ok := strings.CutPrefix(line, "\t\t\tid: "); ok {
			id := strings.Trim(after, "\"")
			quest.Id = id
			continue
		}

		// title
		if after, ok := strings.CutPrefix(line, "\t\t\ttitle: "); ok {
			title := strings.Trim(after, "\"")
			if title[len(title)-1] == '\\' {
				title += "\""
			}
			quest.Title = title
			continue
		}

		// subtitle
		if after, ok := strings.CutPrefix(line, "\t\t\tsubtitle: "); ok {
			subtitle := strings.Trim(after, "\"")
			if subtitle[len(subtitle)-1] == '\\' {
				subtitle += "\""
			}
			quest.Subtitle = subtitle
			continue
		}

		// description
		stripedLine := strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "\t\t\tdescription: "); ok {
			if strings.HasSuffix(stripedLine, "]") {
				description := strings.Trim(after, "[]\"")
				if description[len(description)-1] == '\\' {
					description += "\""
				}
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
			if len(descriptionPart) > 1 && descriptionPart[len(descriptionPart)-1] == '\\' {
				descriptionPart += "\""
			}
			quest.Description = append(quest.Description, descriptionPart)
			continue
		}

		// task titles
		if strings.HasPrefix(line, "\t\t\ttasks:") {
			isTasks = true
			continue
		}

		if after, ok := strings.CutPrefix(stripedLine, "title: "); ok && isTasks {
			taskTitle := strings.Trim(after, "\"")
			if taskTitle[len(taskTitle)-1] == '\\' {
				taskTitle += "\""
			}
			quest.TaskTitles[numTask] = taskTitle
			numTask++
		}

		if isTasks && (strings.HasPrefix(line, "\t\t\t]") || strings.HasPrefix(line, "\t\t\t}]")) {
			isTasks = false
			continue
		}

		// reward titles
		if strings.HasPrefix(line, "\t\t\trewards:") {
			isReward = true
			continue
		}

		if after, ok := strings.CutPrefix(stripedLine, "title: "); ok && isReward {
			rewardTitle := strings.Trim(after, "\"")
			if rewardTitle[len(rewardTitle)-1] == '\\' {
				rewardTitle += "\""
			}
			quest.RewardTitles[numReward] = rewardTitle
			numReward++
		}

		if isReward && (strings.HasPrefix(line, "\t\t\t]") || strings.HasPrefix(line, "\t\t\t}]")) {
			isReward = false
			continue
		}
	}

	return &quest, nil
}

func (q *Quest) GenerateKeys() string {
	result := ""

	isSavingDescription := false
	numDesctription := 0
	numArrayDescription := 0

	var taskNums []int
	numTask := 0
	for key := range q.TaskTitles {
		taskNums = append(taskNums, key)
	}

	var rewardNums []int
	numReward := 0
	for key := range q.RewardTitles {
		rewardNums = append(rewardNums, key)
	}

	isTasks := false
	isReward := false
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
		if strings.HasPrefix(line, "\t\t\ttasks:") {
			isTasks = true
		}

		if _, ok := strings.CutPrefix(strings.TrimSpace(line), "title:"); ok && isTasks {
			parts := strings.Split(line, "title:")
			result += parts[0] + fmt.Sprintf("title: \"{%s.%s.%s.quest%d.task%d.title}\"\n",
				q.ModpackName, q.Chapter, q.Id, q.Number, taskNums[numTask])
			numTask++
			continue
		}

		if isTasks && (strings.HasPrefix(line, "\t\t\t]") || strings.HasPrefix(line, "\t\t\t}]")) {
			isTasks = false
		}

		// reward titles
		if strings.HasPrefix(line, "\t\t\trewards:") {
			isReward = true
		}

		if _, ok := strings.CutPrefix(strings.TrimSpace(line), "title:"); ok && isReward {
			parts := strings.Split(line, "title:")
			result += parts[0] + fmt.Sprintf("title: \"{%s.%s.%s.quest%d.reward%d.title}\"\n",
				q.ModpackName, q.Chapter, q.Id, q.Number, rewardNums[numReward])
			numReward++
			continue
		}

		if isReward && (strings.HasPrefix(line, "\t\t\t]") || strings.HasPrefix(line, "\t\t\t}]")) {
			isReward = false
		}

		result += line + "\n"
	}
	return result
}

func (q *Quest) GenerateMapPart() string {
	numDesctription := 0
	result := ""

	// title
	if q.Title != "" {
		result += fmt.Sprintf(",\n\t\"%s.%s.%s.quest%d.title\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, q.Title)
	}

	// subtitle
	if q.Subtitle != "" {
		result += fmt.Sprintf(",\n\t\"%s.%s.%s.quest%d.subtitle\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, q.Subtitle)
	}

	// description
	for _, description := range q.Description {
		if description != "" {
			result += fmt.Sprintf(",\n\t\"%s.%s.%s.quest%d.description%d\": \"%s\"",
				q.ModpackName, q.Chapter, q.Id, q.Number, numDesctription, description)
			numDesctription++
		}
	}

	// task titles
	for i, titleTask := range q.TaskTitles {
		result += fmt.Sprintf(",\n\t\"%s.%s.%s.quest%d.task%d.title\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, i, titleTask)
	}

	// reward titles
	for i, titleReward := range q.RewardTitles {
		result += fmt.Sprintf(",\n\t\"%s.%s.%s.quest%d.reward%d.title\": \"%s\"",
			q.ModpackName, q.Chapter, q.Id, q.Number, i, titleReward)
	}

	result = strings.TrimPrefix(result, ",\n")

	return result
}
