package types

import "regexp"

var (
	validSlug      = regexp.MustCompile(`^[a-z\-]+$`).MatchString
	validYear      = regexp.MustCompile(`^[\-0-9]+$`).MatchString
	validEventName = regexp.MustCompile(`^^[A-Za-z'0-9\s/&]+$`).MatchString
)
