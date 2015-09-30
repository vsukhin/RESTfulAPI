/* Models package provides data structures for both business logic and database layers */

package models

import (
	"errors"
	"math"
	"time"
)

const (
	FORMAT_DATETIME = "2006-01-02 15:04:05 -0700 MST"
	FORMAT_DATE     = "2006-01-02"
)

var DateLayouts = []string{"2006 Jan 02", "2006-Jan-02", "2006/Jan/02", "2006.Jan.02", "2006-01-02", "2006/01/02", "2006.01.02",
	"02 Jan 2006", "02-Jan-2006", "02/Jan/2006", "02.Jan.2006", "02-01-2006", "02/01/2006", "02.01.2006",
	"02 Jan 06", "02-Jan-06", "02/Jan/06", "02.Jan.06"}

func ParseDate(value string) (date time.Time, err error) {
	for _, layout := range DateLayouts {
		date, err = time.Parse(layout, value)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.New("Date is wrong")
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow

	return newVal
}
