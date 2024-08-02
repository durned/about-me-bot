package openweather

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"about-me-bot/internal/config"
)

func PrettyWeather(c Coordinates, wf WeatherForecast) (text string, iconLink string) {
	return fmt.Sprintf(`⏲️ Right now, outside in %s, %s is:
%s,
🌡️ With a temperature of: <b>%v°C</b>.
The forecast for today is such:
%s,
🌡️ lowest: <b>%v°C</b>, highest: <b>%v°C</b>.
Today's weather can be described with this <i>icon</i>:`,
			c.Name, c.Country, wf.Current.Weather[0].Description, wf.Current.Temperature, wf.Daily[0].Summary, wf.Daily[0].Temp.Min, wf.Daily[0].Temp.Max),
		fmt.Sprint(config.BotCfg.Weather.IconsUrl, wf.Daily[0].Weather[0].Icon, "@2x.png")
}

type Coordinates struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
}

type WeatherForecast struct {
	Current CurrentForecast `json:"current"`
	Daily   []DailyForecast `json:"daily"`
}

type CurrentForecast struct {
	Temperature float64          `json:"temp"`
	Weather     []CurrentWeather `json:"weather"`
}

type CurrentWeather struct {
	Description string `json:"description"`
}

type DailyForecast struct {
	Summary string         `json:"summary"`
	Temp    DailyTemp      `json:"temp"`
	Weather []DailyWeather `json:"weather"`
}

type DailyTemp struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type DailyWeather struct {
	Icon string `json:"icon"`
}

const (
	LimitKey   string = "&limit=1"
	APIDelim   string = "&appid="
	LatKey     string = "lat="
	LonKey     string = "&lon="
	ExcludeKey string = "&exclude=minutely,hourly,alerts"
	UnitsKey   string = "&units=metric"
	LangKey    string = "&lang=en"
)

var (
	ErrNoMatchesFound    error = errors.New("no such city found, try again")
	ErrNoData            error = errors.New("there seems to be no records for weather in this city at the moment")
	NoMatchesFoundClient       = &http.Client{
		Transport: MyFakeService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`[]`)),
			}, nil
		}),
	}
)
