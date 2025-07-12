package kg

type Lang int

const (
	Lang_RU_RU Lang = iota
	Lang_EN_US
)

var (
	Lang_name = map[Lang]string{
		Lang_RU_RU: "ru_ru",
		Lang_EN_US: "en_us",
	}

	// correct language codes
	Lang_array = map[string]Lang{
		"ru_ru": Lang_RU_RU,
		"en_us": Lang_EN_US,
	}
)

func (l Lang) String() string {
	if name, ok := Lang_name[l]; ok {
		return name
	} else {
		return "unknown"
	}
}
