package times

import (
	"testing"
	"time"
)

import (
	"github.com/go-playground/assert/v2"
)

func TestTime(t *testing.T) {
	t.Log(Millisecond())
	t.Log(Second())

	assert.Equal(t, Time13ToDateTime(1666260660000), "2022-10-20 18:11:00")
	assert.Equal(t, Time13ToDate(1666260660000), "2022-10-20")
	assert.Equal(t, Time13ToDateTs17(1666260660000), "20221020181100000")

	assert.Equal(t, DateTs17ToTime13("20221020181100000"), int64(1666260660000))
	assert.Equal(t, Time13ToDateTs14(1666260660000), "20221020181100")
	assert.Equal(t, Time13ToTS8(1666260660000), "20221020")
	assert.Equal(t, DateTs14ToTime13("20221020181100"), int64(1666260660000))
	assert.Equal(t, Date19ToTime13("2022-10-20 18:11:00"), int64(1666260660000))

	t.Log(Time13ToTime(1666260660000))
	t.Log(TimeToTime13(time.Now()))

	assert.Equal(t, Time10ToDate(1666260660), "2022-10-20 18:11:00")
	t.Log(Time10ToTime(1666260660))

	t.Log(CurrentDateTime())
	t.Log(CurrentDate())

	t.Log(AddDatesToTime10(1))
	t.Log(AddDatesToTime13(1))
	t.Log(AddDatesToDateTime(1))
	t.Log(AddDatesToDate(1))
	t.Log(AddDaysToZeroTime10(1))

	t.Log(TimeZoneParse(3600))

	rFC3339ToDateTime, _ := RFC3339ToDateTime("2022-10-02T15:04:05Z")
	assert.Equal(t, rFC3339ToDateTime, "2022-10-02 15:04:05")

	assert.Equal(t, Tim13ToYearMonthSlice(1634724600000, 1666260660000),
		[]int{202110, 202111, 202112, 202201, 202202, 202203, 202204, 202205, 202206, 202207, 202208, 202209, 202210})
}
