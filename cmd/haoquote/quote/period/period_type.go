package period

import "time"

type PeriodType string

const (
	PERIOD_M1  PeriodType = "m1"
	PERIOD_M3  PeriodType = "m3"
	PERIOD_M5  PeriodType = "m5"
	PERIOD_M15 PeriodType = "m15"
	PERIOD_M30 PeriodType = "m30"
	PERIOD_H1  PeriodType = "h1"
	PERIOD_H2  PeriodType = "h2"
	PERIOD_H4  PeriodType = "h4"
	PERIOD_H6  PeriodType = "h6"
	PERIOD_H8  PeriodType = "h8"
	PERIOD_H12 PeriodType = "h12"
	PERIOD_D1  PeriodType = "d1"
	PERIOD_D3  PeriodType = "d3"
	PERIOD_W1  PeriodType = "w1"
	PERIOD_MN  PeriodType = "mn"
)

func Periods() []PeriodType {
	return []PeriodType{
		PERIOD_M1, PERIOD_M3, PERIOD_M5, PERIOD_M15, PERIOD_M30,
		PERIOD_H1, PERIOD_H2, PERIOD_H4, PERIOD_H6, PERIOD_H8, PERIOD_H12,
		PERIOD_D1, PERIOD_D3,
		PERIOD_W1,
		PERIOD_MN,
	}
}

func parse_start_end_time(at time.Time, pt PeriodType) (start, end time.Time) {
	switch pt {
	case PERIOD_M1:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute(), 0, 0, time.Local)
		end = start.Add(time.Duration(1) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M3:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%3, 0, 0, time.Local)
		end = start.Add(time.Duration(3) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M5:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%5, 0, 0, time.Local)
		end = start.Add(time.Duration(5) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M15:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%15, 0, 0, time.Local)
		end = start.Add(time.Duration(15) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_M30:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute()-at.Minute()%30, 0, 0, time.Local)
		end = start.Add(time.Duration(30) * time.Minute).Add(time.Duration(-1) * time.Second)
	case PERIOD_H1:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), 0, 0, 0, time.Local)
		end = start.Add(time.Duration(1) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H2:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%2, 0, 0, 0, time.Local)
		end = start.Add(time.Duration(2) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H4:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%4, 0, 0, 0, time.Local)
		end = start.Add(time.Duration(4) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H6:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%6, 0, 0, 0, time.Local)
		end = start.Add(time.Duration(6) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H8:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%8, 0, 0, 0, time.Local)
		end = start.Add(time.Duration(8) * time.Hour).Add(time.Duration(-1) * time.Second)
	case PERIOD_H12:
		start = time.Date(at.Year(), at.Month(), at.Day(), at.Hour()-at.Hour()%12, 0, 0, 0, time.Local)
		end = start.Add(time.Duration(12) * time.Hour).Add(time.Duration(-1) * time.Second)

	case PERIOD_D1:
		start = time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 0, 1).Add(time.Duration(-1) * time.Second)
	case PERIOD_D3:
		start = time.Date(at.Year(), at.Month(), at.Day()-at.Day()%3, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 0, 3).Add(time.Duration(-1) * time.Second)

	case PERIOD_W1:
		offsend := int(time.Monday - at.Weekday())
		if offsend > 0 {
			offsend = -6
		}
		weekStart := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offsend)
		start = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 0, 7).Add(time.Duration(-1) * time.Second)
	case PERIOD_MN:
		start = time.Date(at.Year(), at.Month(), 1, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 1, 0).Add(time.Duration(-1) * time.Second)
	}
	return start, end
}
