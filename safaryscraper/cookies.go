package safaryscraper

import (
	"net/http"
	"strings"
)

func newCookiesList(rawCookies string) []*http.Cookie {
	var cookies []*http.Cookie

	cookiesList := strings.Split(rawCookies, ";")

	for _, list := range cookiesList {

		values := strings.Split(list, "=")

		cookie := &http.Cookie{
			Name:  values[0],
			Value: values[1],
		}

		cookies = append(cookies, cookie)
	}

	return cookies
}
