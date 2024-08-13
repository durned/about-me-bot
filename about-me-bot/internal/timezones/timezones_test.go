package timezones

import (
	"reflect"
	"testing"
	"time"
)

func TestHHMMStringToUTC(t *testing.T) {
	utcLoc, _ := time.LoadLocation("UTC")
	year, month, day := time.Now().In(utcLoc).Date()

	parsedTime, _ := time.Parse("15:04", "06:20")

	type args struct {
		hhMM string
		loc  *time.Location
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			"time.Parse() err",
			args{
				"TIME",
				&time.Location{},
			},
			time.Time{},
			true,
		},
		{
			"06:20 UTC",
			args{
				"06:20",
				utcLoc,
			},
			time.Date(year, month, day, parsedTime.Hour(), parsedTime.Minute(), 0, 0, utcLoc),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HHMMStringToUTC(tt.args.hhMM, tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("HHMMStringToUTC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HHMMStringToUTC() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatLonToTz(t *testing.T) {
	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Amster",
			args{
				52.371807, 4.896029,
			},
			"Europe/Amsterdam",
		},
		{
			"New York",
			args{
				40.730610, -73.935242,
			},
			"America/New_York",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LatLonToTz(tt.args.lat, tt.args.lon); got != tt.want {
				t.Errorf("LatLonToTz() = %v, want %v", got, tt.want)
			}
		})
	}
}
