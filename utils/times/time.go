package times

import (
	"strconv"
	"time"
)

const (
	FormatDateTime = "2006-01-02 15:04:05"
	FormatDate     = "2006-01-02"
	FormatDateTS17 = "20060102150405000"
	FormatDateTS14 = "20060102150405"
	FormatDateTS8  = "20060102"
	FormatDateTS6  = "200601"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

// Millisecond time.Now() 生成13位时间戳(毫秒)
func Millisecond() int64 {
	return time.Now().In(LOC).UnixMilli()
}

// Second time.Now() 生成10位时间戳(秒)
func Second() int64 {
	return time.Now().In(LOC).Unix()
}

// Time13ToDateTime 13位时间戳转日期和时间 2006-01-02 15:04:05 格式
func Time13ToDateTime(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.In(LOC).Format(FormatDateTime)
}

// Time13ToDate 13位时间戳转日期 2006-01-02 格式
func Time13ToDate(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.In(LOC).Format(FormatDate)
}

// Time13ToDateTs17 13位时间戳转日期Ts17位,例如 20161114143116001 格式
func Time13ToDateTs17(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.In(LOC).Format(FormatDateTS17)
}

// DateTs17ToTime13 日期Ts17位转13位时间戳,例如 20161114143116001 转13位时间戳
func DateTs17ToTime13(t string) int64 {
	tm, _ := time.ParseInLocation(FormatDateTS17, t, LOC)
	return tm.UnixMilli()
}

// Time13ToDateTs14 13位时间戳转日期Ts14位,例如 20161114143116 格式
func Time13ToDateTs14(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.In(LOC).Format(FormatDateTS14)
}

// Time13ToTS8 13位时间戳转日期Ts8位,例如 20161114 格式
func Time13ToTS8(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.In(LOC).Format(FormatDateTS8)
}

// DateTs14ToTime13 日期Ts14位转13位时间戳,例如 20161114143116 转13位时间戳
func DateTs14ToTime13(t string) int64 {
	tm, _ := time.ParseInLocation(FormatDateTS14, t, LOC)
	return tm.UnixMilli()
}

// Date19ToTime13 日期19位转13位时间戳,例如 2016-11-14 14:31:16 转13位时间戳
func Date19ToTime13(t string) int64 {
	tm, _ := time.ParseInLocation(FormatDateTime, t, LOC)
	return tm.UnixMilli()
}

// Time13ToTime 13位时间戳转Time
func Time13ToTime(t int64) time.Time {
	return time.Unix(t/1e3, 0).In(LOC)
}

// TimeToTime13 Time转13位时间戳
func TimeToTime13(t time.Time) int64 {
	return t.In(LOC).UnixMilli()
}

// Time10ToDate 10位时间戳转日期与时间 2006-01-02 15:04:05 格式
func Time10ToDate(t int64) string {
	unix := time.Unix(t, 0)
	return unix.In(LOC).Format(FormatDateTime)
}

// Time10ToTime 10位时间戳转Time
func Time10ToTime(t int64) time.Time {
	return time.Unix(t, 0).In(LOC)
}

// CurrentDateTime time.Now() 获取当前日期与时间 2006-01-02 15:04:05 格式
func CurrentDateTime() string {
	return time.Now().In(LOC).Format(FormatDateTime)
}

// CurrentDate time.Now() 获取当前日期 2006-01-02 格式
func CurrentDate() string {
	return time.Now().In(LOC).Format(FormatDate)
}

// AddDatesToTime10 time.Now() 通过增加days返回days天后的10位时间戳,days可以是负数
func AddDatesToTime10(days int) int64 {
	return time.Now().AddDate(0, 0, days).In(LOC).Unix()
}

// AddDatesToTime13 time.Now() 通过增加days返回days天后的13位时间戳,days可以是负数
func AddDatesToTime13(days int) int64 {
	return time.Now().AddDate(0, 0, days).In(LOC).UnixMilli()
}

// AddDatesToDateTime time.Now() 通过增加days返回days天后的日期与时间 2006-01-02 15:04:05 格式,days可以是负数
func AddDatesToDateTime(days int) string {
	return time.Now().AddDate(0, 0, days).In(LOC).Format(FormatDateTime)
}

// AddDatesToDate time.Now() 通过增加days返回days天后的日期 2006-01-02 格式,days可以是负数
func AddDatesToDate(days int) string {
	return time.Now().AddDate(0, 0, days).In(LOC).Format(FormatDate)
}

// AddDaysToZeroTime10 time.Now() 通过增加days返回days天后 00:00:00 时的10位时间戳,days可以是负数
func AddDaysToZeroTime10(days int) int64 {
	tm := time.Now().AddDate(0, 0, days).In(LOC)
	hour := int64(tm.Hour())
	min := int64(tm.Minute())
	sec := int64(tm.Second())
	return tm.Unix() - hour*3600 - min*60 - sec
}

// TimeZoneParse TimeZoneParse, 返回 1970-01-01 08:00:00 后加 seconds 秒之后的日期与时间 2006-01-02 15:04:05 格式
func TimeZoneParse(seconds int64) string {
	return time.Unix(seconds, 0).In(LOC).Format(FormatDateTime)
}

// RFC3339ToDateTime RFC3339格式数据 2006-01-02T15:04:05Z 转日期与时间 2006-01-02 15:04:05 格式
func RFC3339ToDateTime(t string) (string, error) {
	a, e := time.ParseInLocation(time.RFC3339, t, LOC)
	if e != nil {
		return "", e
	}
	return a.Format(FormatDateTime), nil
}

// Tim13ToYearMonthSlice 输入开始和结束的13位时间戳,输出其包含的年份和月份数组,例如 [202209, 202210]
func Tim13ToYearMonthSlice(startTime int64, endTime int64) []int {
	start, _ := strconv.Atoi(time.Unix(startTime/1e3, 0).In(LOC).Format(FormatDateTS6))
	end, _ := strconv.Atoi(time.Unix(endTime/1e3, 0).In(LOC).Format(FormatDateTS6))
	var month []int
	for i := start; i <= end; i++ {
		str := strconv.Itoa(i)
		n, _ := strconv.Atoi(str[len(str)-2 : len(str)])
		if n > 12 {
			i += 87
			continue
		}
		s, _ := strconv.Atoi(str)
		month = append(month, s)
	}
	return month
}
