package timezones

import (
	"time"

	"about-me-bot/internal/logger"

	"github.com/zsefvlol/timezonemapper"
)

func HHMMStringToUTC(hhMM string, loc *time.Location) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", hhMM)
	if err != nil {
		logger.SimpleLogger.Error("provided time isn't of the right fmt: " + err.Error())
		return time.Time{}, err
	}

	// get user's dd/mm/yy
	// we only use hh:mm in the end
	year, month, day := time.Now().In(loc).Date()

	// compose a time.Time
	userTime := time.Date(year, month, day, parsedTime.Hour(), parsedTime.Minute(), 0, 0, loc)

	return userTime.UTC(), nil
}

func LatLonToTz(lat float64, lon float64) string {
	return timezonemapper.LatLngToTimezoneString(lat, lon)
}
