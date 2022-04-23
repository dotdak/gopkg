package unidecode

// From: https://github.com/mozillazg/go-unidecode

import (
	"strings"
	"unicode"

	"github.com/dotdak/gopkg/unidecode/table"
)

// Version return version
func Version() string {
	return "0.1.1"
}

// Unidecode implements transliterate Unicode text into plain 7-bit ASCII.
// e.g. Unidecode("kožušček") => "kozuscek"
func Unidecode(s string) string {
	return unidecode(s)
}

func UnidecodeTrimNonLetter(s string) string {
	return unidecodeTrimNonLetter(s)
}

func unidecode(s string) string {
	ret := []string{}
	for _, r := range s {
		if r < unicode.MaxASCII {
			v := string(r)
			ret = append(ret, v)
			continue
		}
		if r > 0xeffff {
			continue
		}

		section := r >> 8   // Chop off the last two hex digits
		position := r % 256 // Last two hex digits
		if tb, ok := table.Tables[section]; ok {
			if len(tb) > int(position) {
				ret = append(ret, tb[position])
			}
		}
	}
	return strings.Join(ret, "")
}

func unidecodeTrimNonLetter(s string) string {
	ret := []string{}
	for _, r := range s {
		if ok := unicode.IsLetter(r); !ok && r != ' ' {
			continue
		}
		if r < unicode.MaxASCII {
			ret = append(ret, string(r))
			continue
		}
		if r > 0xeffff {
			continue
		}

		section := r >> 8   // Chop off the last two hex digits
		position := r % 256 // Last two hex digits
		if tb, ok := table.Tables[section]; ok {
			if len(tb) > int(position) {
				ret = append(ret, tb[position])
			}
		}
	}
	return strings.Join(ret, "")
}
