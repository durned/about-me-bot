package holidays

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	_ "net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cubanHolidayJSON = `
[
  {
    "name": "Revolution Anniversary Celebration",
    "name_local": "",
    "language": "",
    "description": "",
    "country": "CU",
    "location": "Cuba",
    "type": "National",
    "date": "07/27/2024",
    "date_year": "2024",
    "date_month": "07",
    "date_day": "27",
    "week_day": "Saturday"
  }
]
`

var multipleHolidaysJSON = `
[
  {
    "name": "Labour Day/May Day",
    "name_local": "",
    "language": "",
    "description": "",
    "country": "LV",
    "location": "Latvia",
    "type": "National",
    "date": "05/01/2024",
    "date_year": "2024",
    "date_month": "05",
    "date_day": "01",
    "week_day": "Wednesday"
  },
  {
    "name": "Constituent Assembly Convocation Day",
    "name_local": "",
    "language": "",
    "description": "",
    "country": "LV",
    "location": "Latvia",
    "type": "National",
    "date": "05/01/2024",
    "date_year": "2024",
    "date_month": "05",
    "date_day": "01",
    "week_day": "Wednesday"
  }
]
`

func TestHoliday(t *testing.T) {
	assertMulti := assert.New(t)
	myError := errors.New("I am an error!")

	type args struct {
		country string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
		client  *http.Client
	}{
		{
			"not a country code",
			args{"COUNTRY"},
			"",
			errors.New("country code provided is not 2 symbols long"),
			nil,
		},
		{
			"no GET request made",
			args{"LV"},
			"",
			fmt.Errorf("could not make a GET request: Get \"%s\": %s", CountryKey+"LV"+DateKey, myError),
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, myError
				}),
			},
		},
		{
			"response code != 200",
			args{"LV"},
			"",
			fmt.Errorf("the GET request went through, but the status code is not an OK (200): %v", http.StatusTooManyRequests),
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
					}, nil
				}),
			},
		},
		{
			"error decoding",
			args{"LV"},
			"",
			fmt.Errorf("could not decode the response body: %s", "invalid character 'p' looking for beginning of value"),
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("probably not JSON")),
					}, nil
				}),
			},
		},
		{
			"no holidays =(",
			args{"LV"},
			fmt.Sprintf("Turns out there are no holidays in %s today..", "Latvia"),
			nil,
			LatvianHolidaysClient,
		},
		{
			"a Cuban holiday!",
			args{"CU"},
			fmt.Sprintf("On %s in %s the Holiday is: %s %s\n", "07/27/2024", "Cuba", "National", "Revolution Anniversary Celebration"),
			nil,
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(cubanHolidayJSON)),
					}, nil
				}),
			},
		},
		{
			"multiple holidays",
			args{"LV"},
			fmt.Sprintf("On %s in %s the Holiday is: %s %s\n", "05/01/2024", "Latvia", "National", "Labour Day/May Day") +
				fmt.Sprintf("On %s in %s the Holiday is: %s %s\n", "05/01/2024", "Latvia", "National", "Constituent Assembly Convocation Day"),
			nil,
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(multipleHolidaysJSON)),
					}, nil
				}),
			},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			MyClient = tt.client
			holiday, err := Holiday(tt.args.country)
			assertMulti.Equal(tt.want, holiday)
			assertMulti.Equal(tt.wantErr, err)
		})
	}
}
