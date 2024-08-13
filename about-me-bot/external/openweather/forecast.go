package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	cfg "about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
)

var ForecastClient = http.DefaultClient

func Forecast(coords Coordinates) (text string, iconLink string, err error) {
	// No real way to check if Coordinates has been initialized
	l.SimpleLogger.Info("making a request to the OpenWeather One Call API")

	link := fmt.Sprintf("%s%s%.2f%s%.2f%s",
		cfg.Global.Weather.ForecastAPIUrl,
		LatKey, coords.Latitude, LonKey, coords.Longitude,
		fmt.Sprint(ExcludeKey, UnitsKey, LangKey,
			APIDelim, cfg.WeatherCfg.APIKey))

	resp, err := ForecastClient.Get(link)
	if err != nil {
		reqErr := "could not make a GET request: " + err.Error()
		l.SimpleLogger.Error(reqErr)
		return "", "", errors.New(reqErr)
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			l.SimpleLogger.Error(ErrNoData.Error())
			return "", "", ErrNoData
		}

		clientErr := fmt.Sprintf("the GET request went through, but the status code is not an OK (200) nor a NotFound (404): %v", resp.StatusCode)
		l.SimpleLogger.Error(clientErr)
		return "", "", errors.New(clientErr)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // probably harder to deliberately make it fail here than to get it to fail in real world
	if err != nil {
		readErr := "could not read the response body: " + err.Error()
		l.SimpleLogger.Error(readErr)
		return "", "", errors.New(readErr)
	}

	var weatherForecast WeatherForecast
	if err = json.Unmarshal(body, &weatherForecast); err != nil {
		decodeErr := "could not decode the response body: " + err.Error()
		l.SimpleLogger.Error(decodeErr)
		return "", "", errors.New(decodeErr)
	}

	text, iconLink = PrettyWeather(coords, weatherForecast)
	return text, iconLink, nil
}
