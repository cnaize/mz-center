package server

const (
	maxSearchLen = 64
)

func ParseSearchText(text string) string {
	if len(text) > maxSearchLen {
		text = text[:maxSearchLen]
	}

	return text
}
