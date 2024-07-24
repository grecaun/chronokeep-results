package types

import "regexp"

var (
	validSlug      = regexp.MustCompile(`^[A-Za-z\-0-9]+$`).MatchString
	validYear      = regexp.MustCompile(`^[\-0-9a-zA-Z]+$`).MatchString
	validEventName = regexp.MustCompile(`^^[A-Za-z'0-9\s/&]+$`).MatchString
)
