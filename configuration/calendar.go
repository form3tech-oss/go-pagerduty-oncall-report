package configuration

import (
	"log"
	"time"

	"os"

	"strings"

	"strconv"

	"github.com/GeertJohan/go.rice"
	"gopkg.in/yaml.v2"
)

type Day struct {
	Time *time.Time
}

func (d *Day) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var timeStr string
	if err := unmarshal(&timeStr); err != nil {
		return err
	}

	// Parse the string to produce a proper time.Time struct.
	pt, err := time.Parse("02/01/2006", timeStr)
	if err != nil {
		return err
	}
	d.Time = &pt
	return nil
}

func (d *Day) ToHashKey() string {
	return d.Time.Format("02/01/2006")
}

type BankHoliday struct {
	Name string `yaml:"title,omitempty"`
	Date Day    `yaml:"date,omitempty"`
}

type BHCalendar struct {
	DaysMaps map[string]BankHoliday // map[date 02-01-2006 format]
}

func (b *BHCalendar) IsDateBankHoliday(date time.Time) bool {
	_, present := b.DaysMaps[date.Format("02/01/2006")]
	return present
}

func (b *BHCalendar) IsWeekend(date time.Time) bool {
	return date.Weekday() == 6 || date.Weekday() == 0
}

type BHCalendars map[string]BHCalendar // map[calendar_name-year]

var BankHolidaysCalendars BHCalendars

func LoadCalendars(year int) {
	log.Printf("Loading calendars for year: %d", year)

	riceConf := rice.Config{
		LocateOrder: []rice.LocateMethod{
			rice.LocateEmbedded,
			rice.LocateAppended,
			rice.LocateFS,
			rice.LocateWorkingDirectory,
		},
	}
	calendarsLocation, err := riceConf.FindBox("./../_assets")
	if err != nil {
		log.Fatalf("cannot find box '_assets': %s", err.Error())
	}

	BankHolidaysCalendars = BHCalendars{}
	err = calendarsLocation.Walk("calendars", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		split := strings.Split(f.Name(), ".")
		if parsedYear, _ := strconv.Atoi(split[2]); parsedYear != year {
			return nil
		}

		key := strings.Join(split[1:len(split)-1], "-")
		fileBytes, e := calendarsLocation.Bytes(path)
		if e != nil {
			panic(e)
		}

		var bankHolidays []BankHoliday
		err = yaml.Unmarshal(fileBytes, &bankHolidays)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		calendar := map[string]BankHoliday{}
		for _, bh := range bankHolidays {
			calendar[bh.Date.ToHashKey()] = bh
		}

		log.Printf("Loaded calendar: '%s'", key)
		BankHolidaysCalendars[key] = BHCalendar{
			DaysMaps: calendar,
		}

		return nil
	})

	if err != nil {
		log.Fatalf("error going through calendars directory: %s", err.Error())
	}
}
