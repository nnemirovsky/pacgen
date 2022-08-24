package regexp

import (
	"fmt"
	"regexp"
)

func DomainAndSubdomains(domain string) string {
	return fmt.Sprintf(`(?:^|\.)%s$`, regexp.QuoteMeta(domain))
}

func Domain(domain string) string {
	return fmt.Sprintf(`^%s$`, regexp.QuoteMeta(domain))
}
