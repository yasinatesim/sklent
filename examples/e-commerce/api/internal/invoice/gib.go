package invoice

import "net/url"

// allowedHosts is the GIB e-Arşiv allowlist; the proxy forwards only to these hosts (no open proxy).
var allowedHosts = map[string]struct{}{
	"earsivportal.efatura.gov.tr":     {},
	"earsivportaltest.efatura.gov.tr": {},
}

// IsGibTarget reports whether target is an https URL pointing at an allowed GIB host.
func IsGibTarget(target string) bool {
	parsed, err := url.Parse(target)
	if err != nil {
		return false
	}
	if parsed.Scheme != "https" {
		return false
	}
	_, ok := allowedHosts[parsed.Host]
	return ok
}
