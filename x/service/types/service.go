package types

import (
	"fmt"
	"net/url"
	"regexp"
)

var (
	errEmptySchemeFmt = "empty scheme in endpoint URL: %s"
	errEmptyHostFmt   = "empty host in endpoint URL: %s"
	errEmptyPortFmt   = "empty port in endpoint URL: %s"
	// NB: limiting the length of the URL scheme to 25 characters to mitigate
	// regex-based DoS attack vectors.
	urlSchemePresenceRegex = regexp.MustCompile(`^\w{0,25}://`)
)

// TODO_INCOMPLETE: Discuss what / how much validation we want to do here.
func (m *ServiceConfig) ValidateEndpoints() error {
	for _, endpoint := range m.Endpoints {
		// Ensure that endpoint URLs contain a scheme to avoid ambiguity when
		// parsing. (See: https://pkg.go.dev/net/url#Parse)
		if !urlSchemePresenceRegex.Match([]byte(endpoint.Url)) {
			return fmt.Errorf(errEmptySchemeFmt, endpoint.Url)
		}

		endpointURL, err := url.Parse(endpoint.Url)
		if err != nil {
			// TODO_CONSIDERATION: accumulate all errors and return at the end.
			// Rationale: save operators time by not having to fix one error at
			// a time.
			return err
		}

		if endpointURL.Host == "" {
			return fmt.Errorf(errEmptyHostFmt, endpoint.Url)
		}
		if endpointURL.Port() == "" {
			return fmt.Errorf(errEmptyPortFmt, endpoint.Url)
		}
	}
	return nil
}
