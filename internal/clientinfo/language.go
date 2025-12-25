package clientinfo

import "strings"

const (
	defaultLocale     = "en"
	defaultFullLocale = "en"
)

// ParseAcceptLanguage derives locale and preferred language from Accept-Language header.
func ParseAcceptLanguage(header string) (string, string) {
	header = strings.TrimSpace(header)
	if header == "" {
		return defaultLocale, defaultFullLocale
	}

	parts := strings.Split(header, ",")
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}

		token = strings.Split(token, ";")[0]

		primary, formatted := normalizeLangToken(token)
		if primary != "" {
			return primary, formatted
		}
	}

	return defaultLocale, defaultFullLocale
}

func normalizeLangToken(token string) (string, string) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ""
	}

	segments := strings.Split(token, "-")

	primary := strings.ToLower(segments[0])
	if primary == "" {
		return "", ""
	}

	if len(segments) == 1 {
		return primary, primary
	}

	builders := make([]string, 0, len(segments))
	builders = append(builders, primary)

	for _, seg := range segments[1:] {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}

		builders = append(builders, strings.ToUpper(seg))
	}

	if len(builders) == 1 {
		return primary, primary
	}

	return primary, strings.Join(builders, "_")
}
