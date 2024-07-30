package holidays

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	cfg "about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
)

var MyClient = http.DefaultClient

func Holiday(country string) (string, error) {
	if len(country) != 2 {
		err := errors.New("country code provided is not 2 symbols long") // case of the string is not a problem
		l.SimpleLogger.Error(err.Error())
		return "", err
	}

	l.SimpleLogger.Info("making a request to the Holidays API")

	link := cfg.BotCfg.Holidays.Endpoint + CountryKey + country + DateKey

	// Multiple holidays on the same day:
	// link = cfg.BotCfg.HolidaysEndpoint + CountryKey + "LV" + "&year=2024&month=05&day=01"

	resp, err := MyClient.Get(link)
	if err != nil {
		reqErr := "could not make a GET request: " + err.Error()
		l.SimpleLogger.Error(reqErr)
		return "", errors.New(reqErr)
	}
	if resp.StatusCode != 200 {
		clientErr := fmt.Sprintf("the GET request went through, but the status code is not an OK (200): %v", resp.StatusCode)
		l.SimpleLogger.Error(clientErr)
		return "", errors.New(clientErr)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // probably harder to deliberately make it than to get in real world
	if err != nil {
		readErr := "could not read the response body: " + err.Error()
		l.SimpleLogger.Error(readErr)
		return "", errors.New(readErr)
	}

	var hData []HolidayData
	if err = json.Unmarshal(body, &hData); err != nil {
		decodeErr := "could not decode the response body: " + err.Error()
		l.SimpleLogger.Error(decodeErr)
		return "", errors.New(decodeErr)
	}

	if !(len(hData) > 0) {
		return fmt.Sprintf("Turns out there are no holidays in %s today..", SupportedCountries[country].string), nil
	}

	var res string
	for i := range hData {
		res += fmt.Sprintf("On %s in %s the Holiday is: %s %s\n", hData[i].Date, hData[i].Location, hData[i].Type, hData[i].Name)
	}
	return res, nil
}
