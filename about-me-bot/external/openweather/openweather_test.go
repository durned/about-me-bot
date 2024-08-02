package openweather

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

// Not every local name is provided in this test, those are skipped anyways and the string would be too long
var RigaCoordinatesJSON string = `
[
  {
    "name": "Riga",
    "local_names": {
      "ie": "Riga",
      "kl": "Riga",
      "cs": "Riga"
    },
    "lat": 56.9493977,
    "lon": 24.1051846,
    "country": "LV",
    "state": "Vidzeme"
  }
]
`

var RigaCoordsDecoded = Coordinates{Name: "Riga", Latitude: 56.9493977, Longitude: 24.1051846, Country: "LV"}

var RigaWeatherDecoded WeatherForecast = WeatherForecast{
	Current: CurrentForecast{
		Temperature: 18.47,
		Weather:     []CurrentWeather{{"clear sky"}},
	},
	Daily: []DailyForecast{
		{
			Summary: "Expect a day of partly cloudy with clear spells",
			Temp:    DailyTemp{Min: 13.64, Max: 21.92},
			Weather: []DailyWeather{{Icon: "03d"}},
		},
	},
}

func TestGeocoder(t *testing.T) {
	assertMulti := assert.New(t)

	dummy := Coordinates{}
	myError := errors.New("I am an error!")

	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    Coordinates
		wantErr error
		client  *http.Client
	}{
		{
			"empty city string",
			args{""},
			dummy,
			errors.New("an empty string has been provided"),
			nil,
		},
		{
			"no GET request made",
			args{"Riga"},
			dummy,
			fmt.Errorf("could not make a GET request: Get \"%s\": %s", fmt.Sprint("q=", "Riga", LimitKey, APIDelim), myError),
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
			args{"Riga"},
			dummy,
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
			args{"Riga"},
			dummy,
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
			"no matches",
			args{"["},
			dummy,
			ErrNoMatchesFound,
			NoMatchesFoundClient,
		},
		{
			"coordinates of Riga", // no multiple responses possible: only the first element is returned
			args{"Riga"},
			RigaCoordsDecoded,
			nil,
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(RigaCoordinatesJSON)),
					}, nil
				}),
			},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			GeocodeClient = tt.client
			coords, err := Geocode(tt.args.city)
			assertMulti.Equal(tt.want, coords)
			assertMulti.Equal(tt.wantErr, err)
		})
	}
}

func TestForecast(t *testing.T) {
	assertMulti := assert.New(t)

	dummyIn := Coordinates{}
	myError := errors.New("I am an error!")

	prettyRigaWeather, iconLink := PrettyWeather(RigaCoordsDecoded, RigaWeatherDecoded)

	type args struct {
		coords Coordinates
	}
	tests := []struct {
		name        string
		args        args
		wantMsgText string
		wantLink    string
		wantErr     error
		client      *http.Client
	}{
		{
			"no GET request made",
			args{dummyIn},
			"",
			"",
			fmt.Errorf("could not make a GET request: Get \"%s\": %s",
				fmt.Sprintf("%s%.2f%s%.2f%s", LatKey, 0.0, LonKey, 0.0, fmt.Sprint(ExcludeKey, UnitsKey, LangKey, APIDelim)),
				myError),
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, myError
				}),
			},
		},
		{
			"response code != 200 && != 404",
			args{dummyIn},
			"",
			"",
			fmt.Errorf("the GET request went through, but the status code is not an OK (200) nor a NotFound (404): %v", http.StatusTooManyRequests),
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
					}, nil
				}),
			},
		},
		{
			"response code = 404",
			args{dummyIn},
			"",
			"",
			ErrNoData,
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
					}, nil
				}),
			},
		},
		{
			"error decoding",
			args{RigaCoordsDecoded},
			"",
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
			"weather in Riga",
			args{RigaCoordsDecoded},
			prettyRigaWeather,
			iconLink,
			nil,
			&http.Client{
				Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(RigaWeatherJSONTest)),
					}, nil
				}),
			},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ForecastClient = tt.client
			text, link, err := Forecast(tt.args.coords)
			assertMulti.Equal(tt.wantMsgText, text)
			assertMulti.Equal(tt.wantLink, link)
			assertMulti.Equal(tt.wantErr, err)
		})
	}
}

var RigaWeatherJSONTest string = `
{
  "lat": 56.95,
  "lon": 24.11,
  "timezone": "Europe/Riga",
  "timezone_offset": 10800,
  "current": {
    "dt": 1722532264,
    "sunrise": 1722478988,
    "sunset": 1722537389,
    "temp": 18.47,
    "feels_like": 18.3,
    "pressure": 1008,
    "humidity": 74,
    "dew_point": 13.75,
    "uvi": 0.2,
    "clouds": 0,
    "visibility": 10000,
    "wind_speed": 5.66,
    "wind_deg": 340,
    "weather": [
      {
        "id": 800,
        "main": "Clear",
        "description": "clear sky",
        "icon": "01d"
      }
    ]
  },
  "daily": [
    {
      "dt": 1722506400,
      "sunrise": 1722478988,
      "sunset": 1722537389,
      "moonrise": 1722462420,
      "moonset": 1722535020,
      "moon_phase": 0.9,
      "summary": "Expect a day of partly cloudy with clear spells",
      "temp": {
        "day": 21.64,
        "min": 13.64,
        "max": 21.92,
        "night": 16.83,
        "eve": 18.47,
        "morn": 15.94
      },
      "feels_like": {
        "day": 21.22,
        "night": 16.55,
        "eve": 18.28,
        "morn": 15.78
      },
      "pressure": 1007,
      "humidity": 52,
      "dew_point": 11.57,
      "wind_speed": 5.65,
      "wind_deg": 330,
      "wind_gust": 10.27,
      "weather": [
        {
          "id": 802,
          "main": "Clouds",
          "description": "scattered clouds",
          "icon": "03d"
        }
      ],
      "clouds": 32,
      "pop": 0.05,
      "uvi": 4.77
    },
    {
      "dt": 1722592800,
      "sunrise": 1722565507,
      "sunset": 1722623663,
      "moonrise": 1722552540,
      "moonset": 1722623520,
      "moon_phase": 0.93,
      "summary": "You can expect partly cloudy in the morning, with clearing in the afternoon",
      "temp": {
        "day": 21.81,
        "min": 14.46,
        "max": 22.29,
        "night": 15.75,
        "eve": 19.79,
        "morn": 15.33
      },
      "feels_like": {
        "day": 21.2,
        "night": 15.44,
        "eve": 19.39,
        "morn": 15.01
      },
      "pressure": 1008,
      "humidity": 44,
      "dew_point": 9.2,
      "wind_speed": 4.48,
      "wind_deg": 334,
      "wind_gust": 8.55,
      "weather": [
        {
          "id": 803,
          "main": "Clouds",
          "description": "broken clouds",
          "icon": "04d"
        }
      ],
      "clouds": 65,
      "pop": 0,
      "uvi": 4.05
    },
    {
      "dt": 1722679200,
      "sunrise": 1722652027,
      "sunset": 1722709936,
      "moonrise": 1722643800,
      "moonset": 1722711060,
      "moon_phase": 0.96,
      "summary": "Expect a day of partly cloudy with rain",
      "temp": {
        "day": 21.12,
        "min": 13.48,
        "max": 22.62,
        "night": 15.87,
        "eve": 19.47,
        "morn": 14.61
      },
      "feels_like": {
        "day": 20.67,
        "night": 15.71,
        "eve": 19.04,
        "morn": 14.48
      },
      "pressure": 1007,
      "humidity": 53,
      "dew_point": 11.31,
      "wind_speed": 3.9,
      "wind_deg": 307,
      "wind_gust": 7.19,
      "weather": [
        {
          "id": 500,
          "main": "Rain",
          "description": "light rain",
          "icon": "10d"
        }
      ],
      "clouds": 26,
      "pop": 0.65,
      "rain": 0.54,
      "uvi": 4.67
    },
    {
      "dt": 1722765600,
      "sunrise": 1722738547,
      "sunset": 1722796206,
      "moonrise": 1722735480,
      "moonset": 1722798060,
      "moon_phase": 0,
      "summary": "Expect a day of partly cloudy with rain",
      "temp": {
        "day": 21.99,
        "min": 13.75,
        "max": 23.48,
        "night": 18.61,
        "eve": 21.47,
        "morn": 14.81
      },
      "feels_like": {
        "day": 21.68,
        "night": 18.46,
        "eve": 21.27,
        "morn": 14.62
      },
      "pressure": 1008,
      "humidity": 55,
      "dew_point": 12.8,
      "wind_speed": 2.7,
      "wind_deg": 297,
      "wind_gust": 6.26,
      "weather": [
        {
          "id": 500,
          "main": "Rain",
          "description": "light rain",
          "icon": "10d"
        }
      ],
      "clouds": 99,
      "pop": 0.62,
      "rain": 0.71,
      "uvi": 4.8
    },
    {
      "dt": 1722852000,
      "sunrise": 1722825068,
      "sunset": 1722882475,
      "moonrise": 1722827100,
      "moonset": 1722884880,
      "moon_phase": 0.03,
      "summary": "You can expect rain in the morning, with partly cloudy in the afternoon",
      "temp": {
        "day": 15.34,
        "min": 14.88,
        "max": 19.28,
        "night": 16.64,
        "eve": 19.28,
        "morn": 14.99
      },
      "feels_like": {
        "day": 15.41,
        "night": 16.55,
        "eve": 19.12,
        "morn": 15.1
      },
      "pressure": 1008,
      "humidity": 95,
      "dew_point": 14.67,
      "wind_speed": 2.03,
      "wind_deg": 204,
      "wind_gust": 3.89,
      "weather": [
        {
          "id": 500,
          "main": "Rain",
          "description": "light rain",
          "icon": "10d"
        }
      ],
      "clouds": 100,
      "pop": 1,
      "rain": 4.87,
      "uvi": 4.87
    },
    {
      "dt": 1722938400,
      "sunrise": 1722911590,
      "sunset": 1722968742,
      "moonrise": 1722918480,
      "moonset": 1722971580,
      "moon_phase": 0.06,
      "summary": "Expect a day of partly cloudy with rain",
      "temp": {
        "day": 21.16,
        "min": 15.76,
        "max": 22.13,
        "night": 17.82,
        "eve": 21.25,
        "morn": 15.76
      },
      "feels_like": {
        "day": 21.03,
        "night": 17.8,
        "eve": 21.02,
        "morn": 15.92
      },
      "pressure": 1012,
      "humidity": 65,
      "dew_point": 14.29,
      "wind_speed": 3.97,
      "wind_deg": 322,
      "wind_gust": 8.01,
      "weather": [
        {
          "id": 500,
          "main": "Rain",
          "description": "light rain",
          "icon": "10d"
        }
      ],
      "clouds": 9,
      "pop": 0.26,
      "rain": 0.26,
      "uvi": 5
    },
    {
      "dt": 1723024800,
      "sunrise": 1722998112,
      "sunset": 1723055007,
      "moonrise": 1723009680,
      "moonset": 1723058160,
      "moon_phase": 0.09,
      "summary": "There will be partly cloudy today",
      "temp": {
        "day": 22.56,
        "min": 15.61,
        "max": 23.58,
        "night": 18.9,
        "eve": 22.67,
        "morn": 15.8
      },
      "feels_like": {
        "day": 22.54,
        "night": 18.91,
        "eve": 22.56,
        "morn": 15.94
      },
      "pressure": 1017,
      "humidity": 64,
      "dew_point": 15.51,
      "wind_speed": 2.93,
      "wind_deg": 335,
      "wind_gust": 5.13,
      "weather": [
        {
          "id": 802,
          "main": "Clouds",
          "description": "scattered clouds",
          "icon": "03d"
        }
      ],
      "clouds": 34,
      "pop": 0,
      "uvi": 5
    },
    {
      "dt": 1723111200,
      "sunrise": 1723084634,
      "sunset": 1723141271,
      "moonrise": 1723100700,
      "moonset": 1723144740,
      "moon_phase": 0.12,
      "summary": "There will be partly cloudy today",
      "temp": {
        "day": 25.29,
        "min": 15.67,
        "max": 27.62,
        "night": 22.58,
        "eve": 27.39,
        "morn": 15.67
      },
      "feels_like": {
        "day": 25.21,
        "night": 22.83,
        "eve": 27.66,
        "morn": 15.62
      },
      "pressure": 1017,
      "humidity": 51,
      "dew_point": 14.6,
      "wind_speed": 3.56,
      "wind_deg": 205,
      "wind_gust": 4.95,
      "weather": [
        {
          "id": 802,
          "main": "Clouds",
          "description": "scattered clouds",
          "icon": "03d"
        }
      ],
      "clouds": 31,
      "pop": 0,
      "uvi": 5
    }
  ]
}
`
