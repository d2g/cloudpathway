package web

type Configuration struct {
	Port          string
	Files         string
	DeveloperMode bool
	Sessions      struct {
		Timeout int64
	}
}
