package clientinfo

import "testing"

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		header    string
		locale    string
		preferred string
	}{
		{"en-US,en;q=0.9", "en", "en_US"},
		{"ru", "ru", "ru"},
		{"fr-CA,fr;q=0.8,en;q=0.6", "fr", "fr_CA"},
		{"", defaultLocale, defaultFullLocale},
	}

	for _, tt := range tests {
		locale, preferred := ParseAcceptLanguage(tt.header)
		if locale != tt.locale || preferred != tt.preferred {
			t.Fatalf("ParseAcceptLanguage(%q) = (%s,%s), want (%s,%s)", tt.header, locale, preferred, tt.locale, tt.preferred)
		}
	}
}
