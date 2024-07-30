package holidays

import (
	"net/http"
)

type MyFakeService func(*http.Request) (*http.Response, error)

func (s MyFakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}
