package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	cfg "about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
)

var GeocodeClient = http.DefaultClient

func Geocode(city string) (Coordinates, error) {
	dummy := Coordinates{}

	if !(len(city) > 0) {
		err := errors.New("an empty string has been provided") // case of the string is not a problem, but spaces are
		l.SimpleLogger.Error(err.Error())
		return dummy, err
	}
	city = strings.ReplaceAll(city, " ", "%20")

	l.SimpleLogger.Info("making a request to the Geocoding API")

	link := fmt.Sprint(cfg.Global.Weather.GeocodingAPIUrl, "q=", city, LimitKey, APIDelim, cfg.Global.Weather.APIKey)

	resp, err := GeocodeClient.Get(link)
	if err != nil {
		reqErr := "could not make a GET request: " + err.Error()
		l.SimpleLogger.Error(reqErr)
		return dummy, errors.New(reqErr)
	}
	if resp.StatusCode != 200 {
		clientErr := fmt.Sprintf("the GET request went through, but the status code is not an OK (200): %v", resp.StatusCode)
		l.SimpleLogger.Error(clientErr)
		return dummy, errors.New(clientErr)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // probably harder to deliberately make it fail here than to get it to fail in real world
	if err != nil {
		readErr := "could not read the response body: " + err.Error()
		l.SimpleLogger.Error(readErr)
		return dummy, errors.New(readErr)
	}

	var coordData []Coordinates
	if err = json.Unmarshal(body, &coordData); err != nil {
		decodeErr := "could not decode the response body: " + err.Error()
		l.SimpleLogger.Error(decodeErr)
		return dummy, errors.New(decodeErr)
	}

	if !(len(coordData) > 0) {
		l.SimpleLogger.Error(ErrNoMatchesFound.Error())
		return dummy, ErrNoMatchesFound
	}

	return coordData[0], nil
}
