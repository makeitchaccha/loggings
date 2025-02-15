package settings

import (
	"strings"
	"time"
)

func SanitizeFormat(format string) string {
	replaces := map[string]string{
		"{timestamp}":  time.Now().UTC().Format(time.RFC3339),
		"{avatar url}": "https://cdn.discordapp.com/embed/avatars/0.png",
	}

	for k, v := range replaces {
		format = strings.ReplaceAll(format, k, v)
	}

	return format
}
