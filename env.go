package tangerino

// Environment holds the API base URLs for a deployment target.
// Most resources live under employerBaseURL; punch clock endpoints use punchesBaseURL.
type Environment struct {
	employerBaseURL string
	punchesBaseURL  string
}

var prodEnv = Environment{
	employerBaseURL: "https://employer.tangerino.com.br",
	punchesBaseURL:  "https://apis.tangerino.com.br",
}

var stagingEnv = Environment{
	employerBaseURL: "https://employer.tangerino.com.br",
	punchesBaseURL:  "https://apis.tangerino.com.br",
}
