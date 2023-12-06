package utils

import (
	"digi-model-engine/models"
	"digi-model-engine/utils/constants"
	"time"
)

type Date struct{}

// Returns curr datetime
func (d *Date) CurrentDate() time.Time {
	return time.Now()
}

// Converts date to string in the specified format
func (d *Date) DateToString(date time.Time, format string) string {
	return date.Format(format)
}

// Converts date string to time object using the specified format
func (d *Date) StringToDate(dateString, format string) (time.Time, error) {
	return time.Parse(format, dateString)
}

// Calculates difference between two dates in seconds
func (d *Date) DateDiffInSeconds(startDate, endDate time.Time) int64 {
	diff := endDate.Sub(startDate)
	return int64(diff.Seconds())
}

// Calculates the difference between two dates in milliseconds
func (d *Date) DateDiffInMilliseconds(startDate, endDate time.Time) int64 {
	diff := endDate.Sub(startDate)
	elapsedMs := int64(diff.Milliseconds())
	return elapsedMs
}

// Calculates the difference between two dates in hours
func (d *Date) DateDiffInHours(startDate, endDate time.Time) int {
	daysToHours := int64(d.DateDiffInDays(startDate, endDate)) * 24
	diffBtwTwoTimes := int64(d.DateDiffInSeconds(startDate, endDate)) / 3600
	return int(daysToHours + diffBtwTwoTimes)
}

// Calculates the difference between two dates in days
func (d *Date) DateDiffInDays(startDate, endDate time.Time) int {
	return int(endDate.Sub(startDate).Hours() / 24)
}

// Converts given date into ISO
func (d *Date) GetISODate(dateTime time.Time) time.Time {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return dateTime.In(location)
}

// Returns current date and time truncated to seconds
func (d *Date) CurrentDateTillSeconds() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
}

func GetDateInfo(currTime time.Time) models.DateInfoModel {

	loc, _ := time.LoadLocation("Asia/Kolkata")

	// Today's date
	toDate := currTime.Format(constants.DateTimeFormat)
	// today's ISO Date
	toISODate := currTime.In(loc)
	// Current month's 1st day
	toMonthStartDate := time.Date(currTime.Year(), currTime.Month(), 1, 0, 0, 0, 0, loc).Format(constants.DateTimeFormat)
	// Current month's 1st day ISO
	toMonthStartISODate := time.Date(currTime.Year(), currTime.Month(), 1, 0, 0, 0, 0, loc).In(loc)
	// today's date till midnight-1sec
	toMonthEndDate := time.Date(currTime.Year(), currTime.Month(), daysInMonth(currTime), 23, 59, 59, 0, loc).Format(constants.DateTimeFormat)
	// today's date till midnight-1sec ISO
	toMonthEndISODate := time.Date(currTime.Year(), currTime.Month(), daysInMonth(currTime), 23, 59, 59, 0, loc).In(loc)

	// 2 weeks past ISO date
	fromTwoWeeksISODate := toISODate.AddDate(0, 0, -14)
	// 2 weeks past string date
	fromTwoWeeksDate := fromTwoWeeksISODate.Format(constants.DateTimeFormat)
	// 1 months before today's date
	from1MonthDate := currTime.AddDate(0, -1, 0).Format(constants.DateTimeFormat)
	// 2 months before today's date
	from2MonthDate := currTime.AddDate(0, -2, 0).Format(constants.DateTimeFormat)
	// 3 months before today's date
	from3MonthDate := currTime.AddDate(0, -3, 0).Format(constants.DateTimeFormat)
	// 3 months ISO Date before today's date
	from3MonthISODate := currTime.AddDate(0, -3, 0).In(loc)
	// 3 months ISO Midnight Date before today's date
	from3MonthISOMidnight := time.Date(from3MonthISODate.Year(), from3MonthISODate.Month(), from3MonthISODate.Day(), 0, 0, 0, 0, loc)
	// 2 month's ago 1st day
	from3MStartDate := time.Date(currTime.Year(), currTime.Month()-2, 1, 0, 0, 0, 0, loc).Format(constants.DateTimeFormat)
	// 6 months before today's date
	from6MonthDate := currTime.AddDate(0, -6, 0).Format(constants.DateTimeFormat)
	// 5 month's ago 1st day
	from6MStartDate := time.Date(currTime.Year(), currTime.Month()-5, 1, 0, 0, 0, 0, loc).Format(constants.DateTimeFormat)
	// 11 month's ago 1st day
	from12MStartDate := time.Date(currTime.Year(), currTime.Month()-11, 1, 0, 0, 0, 0, loc).Format(constants.DateTimeFormat)
	// One year before today's date
	fromDate := currTime.AddDate(-1, 0, 0).Format(constants.DateTimeFormat)
	// One year before today's Date ISO
	fromISODate := currTime.AddDate(-1, 0, 0).In(loc)

	months := make([]time.Time, 0)
	EMIMonths := make([]string, 0)
	tempDate := currTime

	for i := 0; i < 12; i++ {
		months = append(months, tempDate)
		EMITempDate := tempDate.Format("01-2006")
		EMITempDate = EMITempDate[1:]
		EMIMonths = append(EMIMonths, EMITempDate)
		tempDate = tempDate.AddDate(0, -1, 0)
	}

	return models.DateInfoModel{
		Month1:  months[0].Format(constants.MonthFormat),
		Month2:  months[1].Format(constants.MonthFormat),
		Month3:  months[2].Format(constants.MonthFormat),
		Month4:  months[3].Format(constants.MonthFormat),
		Month5:  months[4].Format(constants.MonthFormat),
		Month6:  months[5].Format(constants.MonthFormat),
		Month7:  months[6].Format(constants.MonthFormat),
		Month8:  months[7].Format(constants.MonthFormat),
		Month9:  months[8].Format(constants.MonthFormat),
		Month10: months[9].Format(constants.MonthFormat),
		Month11: months[10].Format(constants.MonthFormat),
		Month12: months[11].Format(constants.MonthFormat),

		EMIMonth1:  EMIMonths[0],
		EMIMonth2:  EMIMonths[1],
		EMIMonth3:  EMIMonths[2],
		EMIMonth4:  EMIMonths[3],
		EMIMonth5:  EMIMonths[4],
		EMIMonth6:  EMIMonths[5],
		EMIMonth7:  EMIMonths[6],
		EMIMonth8:  EMIMonths[7],
		EMIMonth9:  EMIMonths[8],
		EMIMonth10: EMIMonths[9],
		EMIMonth11: EMIMonths[10],
		EMIMonth12: EMIMonths[11],

		ToDate:              toDate,
		ToISODate:           toISODate,
		ToMonthEndDate:      toMonthEndDate,
		ToMonthEndISODate:   toMonthEndISODate,
		ToMonthStartDate:    toMonthStartDate,
		ToMonthStartISODate: toMonthStartISODate,

		FromTwoWeeksISODate:   fromTwoWeeksISODate,
		FromTwoWeeksDate:      fromTwoWeeksDate,
		From1MDate:            from1MonthDate,
		From2MDate:            from2MonthDate,
		From3MDate:            from3MonthDate,
		From3MISODate:         from3MonthISODate,
		From3MStartDate:       from3MStartDate,
		From3MISOMidnightDate: from3MonthISOMidnight,
		From6MDate:            from6MonthDate,
		From6MStartDate:       from6MStartDate,
		From12MStartDate:      from12MStartDate,
		FromDate:              fromDate,
		FromISODate:           fromISODate,
	}
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
