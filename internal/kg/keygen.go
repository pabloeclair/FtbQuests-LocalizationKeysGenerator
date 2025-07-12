package kg

import "fmt"

func ParseQuests(modpackName string, chapter string, text string) ([]Quest, error) {
	result := []Quest{}
	return result, nil
}

func GenerateKeys(quests []Quest) string {
	result := ""
	return result
}

func GenerateMap(lang string, quests []Quest) (string, error) {
	l := Lang_array[lang]
	if l.String() == "unknown" {
		return "", fmt.Errorf("unknown language %s; please, check correct codes in internal/kg/lang.go", l)
	}

	return "", nil
}
