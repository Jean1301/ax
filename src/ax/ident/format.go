package ident

// FormatID returns a compact rendering of ID for use in log messages and other
// human-readable strings.
func FormatID(id string) string {
	if looksLikeUUID(id) {
		return id[:uuidSep1] + id[uuidLen:]
	}

	return id
}

func looksLikeUUID(id string) bool {
	if len(id) < uuidLen {
		return false
	}

	return id[uuidSep1] == '-' &&
		id[uuidSep2] == '-' &&
		id[uuidSep3] == '-' &&
		id[uuidSep4] == '-'
}

const (
	uuidLen  = 36
	uuidSep1 = 8
	uuidSep2 = 13
	uuidSep3 = 18
	uuidSep4 = 23
)