package utils

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// GetMonthStartAndEnd returns a month's start date and end date when given
// a "yearMonth" string, i.e. 20185 or 201812
func GetMonthStartAndEnd(yearMonth string) (startDte time.Time, endDte time.Time) {
	year := time.Now().Year()
	month := int(time.Now().Month())

	if len(yearMonth) > 4 {
		year = ParseInt(yearMonth[0:4], 0)
		month = ParseInt(yearMonth[4:], 0)
	}
	startDte = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDte = startDte.AddDate(0, 1, -1)

	return startDte, endDte
}

// GetMonthEnd returns the end of the month when fed, e.g., 20187 => 7/31/2018
func GetMonthEnd(yearMonth string) (endDate time.Time) {

	year := time.Now().Year()
	month := int(time.Now().Month())

	if len(yearMonth) > 4 {
		year = ParseInt(yearMonth[0:4], 0)
		month = ParseInt(yearMonth[4:], 0)
	}

	startDte := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return startDte.AddDate(0, 1, -1)
}

// GetMonthEnd returns the end of the month when fed, e.g., 20187 => 7/31/2018
func GetMonthEndIn(yearMonth string, location *time.Location) (endDate time.Time) {

	year := time.Now().In(location).Year()
	month := int(time.Now().In(location).Month())

	if len(yearMonth) > 4 {
		year = ParseInt(yearMonth[0:4], 0)
		month = ParseInt(yearMonth[4:], 0)
	}

	startDte := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location)
	return startDte.AddDate(0, 1, -1)
}

// DaysIn returns the number of days in a month for a given year.
func DaysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetNearestWeekday returns the most recent weekday before today.
// Used to provide a date for to get a mutual fund price. It has
// to be before today and not a weekend to have a price.
func GetNearestWeekday(dte time.Time) time.Time {

	if dte.Weekday() == 1 { //This is Monday, return last Friday
		return dte.Add(-72 * time.Hour) //Three days ago
	}

	if dte.Weekday() == 7 { //This is Sunday, return last Friday
		return dte.Add(-48 * time.Hour) //Three days ago
	}

	return dte.Add(-24 * time.Hour) // return yesterday

}

// DateEqual compares date equality regardless of the time
func DateEqual(t1 time.Time, t2 time.Time) bool {
	if t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay() {
		return true
	}
	return false
}

// ParseDateUS tries three formats mm/dd/yyyy, m/dd/yyyy, m/d/yyyy, and mm/d/yyyy
func ParseDateUS(candidate string, defaultResult time.Time) (time.Time, error) {
	candidate = strings.Replace(candidate, "-", "/", -1)

	dateSplit := strings.Split(candidate, "/")
	if len(dateSplit) == 2 {
		//in this case, the date looks like, e.g., 1/15, without the year
		//so append a slash and the current year, like 1/15/2018
		candidate = candidate + fmt.Sprintf("/%v", time.Now().Year())
	}

	t, err := dateparse.ParseAny(candidate)
	if err != nil {
		return defaultResult, err
	}

	return t, err

}

// ParseDate is a fast parse for date []byte formatted as
// yyyy-mm-dd
func ParseDate(date []byte, location *time.Location) (time.Time, error) {
	if len(string(date)) != 10 {
		return time.Time{}, fmt.Errorf(`date "%s" not in right format`, string(date))
	}

	year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
	month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
	day := (int(date[8])-'0')*10 + int(date[9]) - '0'
	return time.Date(year, month, day, 0, 0, 0, 0, location), nil
}

// ParseDateTime3 is a fast parse for date-time []byte formatted as
// yyyy-mm-dd hh:mm:ss
func ParseDateTime3(date []byte, location *time.Location) (time.Time, error) {
	if len(string(date)) != 19 {
		return time.Time{}, fmt.Errorf(`date "%s" not in right format`, string(date))
	}

	year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
	month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
	day := (int(date[8])-'0')*10 + int(date[9]) - '0'
	hour := (int(date[11])-'0')*10 + int(date[12]) - '0'
	minute := (int(date[14])-'0')*10 + int(date[15]) - '0'
	second := (int(date[17])-'0')*10 + int(date[18]) - '0'
	return time.Date(year, month, day, hour, minute, second, 0, location), nil
}

// ParseDateTime4 is a fast parse for date-time string formatted as
// yyyy-mm-dd hh:mm:ss
func ParseDateTime4(date string, location *time.Location) (time.Time, error) {
	if len(date) != 19 {
		return time.Time{}, fmt.Errorf(`date "%s" not in right format`, date)
	}

	year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
	month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
	day := (int(date[8])-'0')*10 + int(date[9]) - '0'
	hour := (int(date[11])-'0')*10 + int(date[12]) - '0'
	minute := (int(date[14])-'0')*10 + int(date[15]) - '0'
	second := (int(date[17])-'0')*10 + int(date[18]) - '0'
	return time.Date(year, month, day, hour, minute, second, 0, location), nil
}

// ParseDateTime5 is a fast parse for date-time string formatted as
// yyyy-mm-dd hh:mm
func ParseDateTime5(date string, location *time.Location) (time.Time, error) {
	if len(date) != 16 {
		return time.Time{}, fmt.Errorf("date string not 16 characters")
	}
	year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
	month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
	day := (int(date[8])-'0')*10 + int(date[9]) - '0'
	hour := (int(date[11])-'0')*10 + int(date[12]) - '0'
	minute := (int(date[14])-'0')*10 + int(date[15]) - '0'

	return time.Date(year, month, day, hour, minute, 0, 0, location), nil
}

// ParseDateTime parses the suppied string in location America/New York.
// This parse is VERY SLOW.
func ParseDateTime(candidate string, defaultResult time.Time, location *time.Location) (time.Time, error) {

	tryCandidate, timeErr := time.ParseInLocation("2006-01-02 15:04:05", candidate, location)

	if timeErr != nil {
		return defaultResult, timeErr
	}

	return tryCandidate, nil

}

// WeeksInMonth returns the number of weeks in the month given
// by now.
func WeeksInMonth(now time.Time) int {
	return 0
}

// WeeksInMonth2 ...
func WeeksInMonth2(now time.Time, location *time.Location) int {
	dayThreshold := []int{5, 1, 5, 6, 5, 6, 5, 5, 6, 5, 6, 5}
	currentYear, currentMonth, _ := now.Date()

	firstDay := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, location)
	baseWeeks := 5
	if int(currentMonth) == 2 {
		// only February can fit in 4 weeks
		baseWeeks = 4
	}
	// TODO: account for leap years

	// add an extra week if the month starts beyond the threshold day.
	adjustThreshold := 0

	if int(firstDay.Weekday()) >= dayThreshold[int(currentMonth)] {
		adjustThreshold = 1
	}

	return baseWeeks + adjustThreshold
}

// TimeFromFloat takes a number like 13.5 and pairs it with dateVal to make a date-time value
// If dateVal is "nil" (time.Time{}), the func will pair hourMinute with today's date
// et = current time in eastern time zone
func TimeFromFloat(hourMinute float64, dateVal time.Time, et time.Time) (result time.Time) {

	hour, fractionOfHour := math.Modf(hourMinute)
	minute := float64(60) * fractionOfHour

	if dateVal.After(time.Time{}) { //If we have a non-nil time
		result = time.Date(dateVal.Year(), dateVal.Month(), dateVal.Day(), int(hour), int(minute), 0, 0, et.Location())
	} else { //Use today's date as the date value
		result = time.Date(et.Year(), et.Month(), et.Day(), int(hour), int(minute), 0, 0, et.Location())
	}

	return
}

// MinutesFromFloat ...
func MinutesFromFloat(val float64) int {
	hour, fractionOfHour := math.Modf(val)
	minute := float64(60) * fractionOfHour
	minute = minute + (hour * 60)

	return int(minute)
}

// TimeToFloat takes time (e.g., 10:30) and converts it
// into a float 10.5
func TimeToFloat(val string) (float64, error) {
	//Split time on colon
	timeArray := strings.Split(val, ":")

	if len(timeArray) != 2 {
		return 0, fmt.Errorf("cannot parse time into float")
	}

	hour := ParseFloat64(timeArray[0])
	minute := (ParseFloat64(timeArray[1]) / 60)

	return hour + minute, nil
}

// TimeToInt64 takes time (e.g., 10:30) and converts it
// into an int64, 1030
func TimeToInt64(val string) int64 {
	return ParseInt64(strings.Replace(val, ":", "", -1), 0)
}

// EasternTime returns current Eastern Time (Princeton),
// including DST as appropriate
func EasternTime() time.Time {
	ET, err := time.LoadLocation("America/New_York")

	if err != nil {
		return time.Time{}
	}
	return time.Now().In(ET)
}

// CurrentAcademicYear returns Now()'s
// academic year
func CurrentAcademicYear(monthYearEnd time.Month) int64 {
	return AcademicYear(time.Now(), monthYearEnd)
}

// AcademicYearView returns 2019-20 when given 2020.
func AcademicYearView(academicYear int64) string {
	start := academicYear - 1

	return fmt.Sprintf("%d-%d", start, academicYear)
}

// AcademicYear returns the academic year of *dte*. If monthYearEnd is 6 (June, in North America), and
// the current date is anywhere between July 1 2024 and June 20 2025, the returned academic year is 2025
func AcademicYear(dte time.Time, monthYearEnd time.Month) int64 {

	ay := int64(dte.Year())

	m := dte.Month()

	if m > monthYearEnd {
		ay++
	}
	return ay
}

// EndOfAcademicYear is
func EndAcademicYear(dte time.Time, monthYearEnd time.Month, loc *time.Location) time.Time {

	firstDateOfNextYear := BeginAcademicYear(dte, monthYearEnd, loc).AddDate(1, 0, 0)
	return firstDateOfNextYear.AddDate(0, 0, -1).In(loc)

}

// BeginAcademicYear is
func BeginAcademicYear(dte time.Time, monthYearEnd time.Month, loc *time.Location) time.Time {

	yearBeginMonth := monthYearEnd + 1
	if monthYearEnd == 12 {
		yearBeginMonth = 1
	}

	currentAcademicYear := int(AcademicYear(dte, monthYearEnd))

	//Calculate the first day of the academic year
	return time.Date(currentAcademicYear-1, yearBeginMonth, 01, 00, 00, 00, 000, loc)

}

// AcademicYears returns a slice of ints
// from *start* to dte
func AcademicYears(start int64, dte time.Time, monthYearEnd time.Month) (result []int64) {

	for i := AcademicYear(dte, monthYearEnd); i > start; i-- {
		result = append(result, i)
	}

	return result
}

// BeginningOfDay returns 12AM of the supplied time
func BeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// BeginningOfHour returns :00 of the supplied time
func BeginningOfHour(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), 0, 0, 0, t.Location())
}

// EndOfDay returns 12AM of the supplied time
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 1000, t.Location())
}

// Int64ToAmPm takes *hour* like 13 and
// returns 1p
func Int64ToAmPm(hour int64) string {
	remainder := hour % 12

	if remainder == 0 {
		remainder = 12
	}

	display := fmt.Sprintf("%v", remainder)

	if hour >= 12 {
		display += "p"
	} else {
		display += "a"
	}

	return display
}

// To8601Format produces a time in 8601 format for iCalendar
// and other functions.
func To8601Format(val time.Time) string {
	return val.UTC().Format("20060102T150405Z")
}

// FirstDayOfISOWeek returns time.Time when fed a year and week
func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday -- ISO weeks begin on Monday by definition
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		_, isoWeek = date.ISOWeek()
	}

	return date
}

// MakeTimeFromTimeField takes the value from
// an HTML time field, combines
// it with *date* (can be IsZero, which then uses today's date),
// and returns a time.Time object
func MakeTimeFromTimeField(atTime string, date time.Time, loc *time.Location) (time.Time, error) {

	var (
		hour   int
		minute int
	)

	//Try splitting the time at a colon, e.g, 14:15
	split := strings.Split(atTime, ":")
	if len(split) == 2 {
		hour = ParseInt(split[0], 0)
		minute = ParseInt(split[1], 0)
	} else if len(atTime) == 4 {
		hour = ParseInt(atTime[0:2], 0)
		minute = ParseInt(atTime[3:], 0)
	} else {
		return time.Time{}, fmt.Errorf("error parsing date time string supplied (%s); must have 4 characters", atTime)
	}

	if date.IsZero() {
		date = time.Now().In(loc)
	}
	//We just need to store the time, easy enough to do in a time/date field,
	//so use today's date for the year, month, and day values
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc), nil

}

// GetLocationFromTZ returns location from tz unless it's blank or unless
// there's an error, and otherwise returns the defaultLoc
func GetLocationFromTZ(tz string, defaultLoc *time.Location) *time.Location {

	if tz == "" {
		return defaultLoc
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		return defaultLoc
	}
	return location
}

// IsToday determins whether supplied date is today
func IsToday(dte time.Time) bool {
	yearNow, monthNow, dayNow := time.Now().Date()
	yearDte, monthDte, dayDte := dte.Date()

	if yearNow == yearDte && monthNow == monthDte && dayNow == dayDte {
		return true
	}

	return false

}

// AcademicYearsAgo returns now.Year - n + 1
// to account for full academic years
func AcademicYearsAgo(n int64, monthYearEnd time.Month) int64 {

	thisYear := AcademicYear(time.Now(), monthYearEnd)

	//Why + 1. Suppose we're in the 2019 academic year. Ten years
	//ago is 2009, but not quite, because
	//we don't want to count year 0 (2009). If we add one to year 0,
	//we get... 2010, 2011...,2019, or 10 complete years.
	return thisYear - n + 1
}

// InFuture checks whether the provided dte is after now
func InFuture(dte time.Time, loc *time.Location) bool {
	now := time.Now().In(loc)
	return dte.After(now)
}

// InPast checks whether the provided dte is before now
func InPast(dte time.Time, loc *time.Location) bool {
	now := time.Now().In(loc)
	return dte.Before(now)
}
