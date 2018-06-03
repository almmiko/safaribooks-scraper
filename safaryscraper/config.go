package safaryscraper

type Config struct {
	Url        string
	RawCookies string
}

func NewConfig(url string, rawCookies string) *Config {
	return &Config{
		Url:        url,
		RawCookies: rawCookies,
	}
}
