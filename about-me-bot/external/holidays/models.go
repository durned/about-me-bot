package holidays

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HolidayData struct {
	Location string `json:"location"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	Name     string `json:"name"`
}

const (
	CountryKey string = "&country="
)

var (
	SupportedCountries = map[string]struct{ string }{
		"LV": {"Latvia"},
		"NL": {"The Netherlands"},
		"UA": {"Ukraine"},
		"BR": {"Brazil"},
		"FJ": {"Fiji"},
		"JP": {"Japan"},
		"CU": {"Cuba"},
		"TD": {"Chad"},
	}
	DateKey = fmt.Sprintf("&year=%s&month=%s&day=%s", time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"))
)

var LatvianHolidaysClient = &http.Client{
	Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	}),
}
