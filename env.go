package tangerino

// Environment holds the API base URL for a deployment target.
type Environment struct {baseURL string}

var prodEnv    = Environment{baseURL: "https://employer.tangerino.com.br"}
var	stagingEnv = Environment{baseURL: "https://employer.tangerino.com.br"}

