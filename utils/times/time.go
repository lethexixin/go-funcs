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
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

// Millisecond 生成13位时间戳(毫秒)
func Millisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// Second 生成10位时间戳(秒)
func Second() int64 {
	return time.Now().UnixNano() / 1e9
}

// Time13ToDateTime 13位时间戳转日期和时间 2006-01-02 15:04:05
func Time13ToDateTime(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.Format(FormatDateTime)
}

// Time13ToDate 13位时间戳转日期 2006-01-02
func Time13ToDate(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.Format(FormatDate)
}

// Time13ToTS17 13位时间戳转日期Ts17位,例如20161114143116001
func Time13ToDateTs17(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.Format(FormatDateTS17)
}

// DateTs17ToTime13 日期Ts17位转13位时间戳,例如20161114143116001转13位时间戳
func DateTs17ToTime13(t string) int64 {
	tm, _ := time.ParseInLocation(FormatDateTS17, t, LOC)
	return tm.UnixNano() / 1e6
}

// Time13ToTS14 13位时间戳转日期Ts14位,例如20161114143116
func Time13ToDateTs14(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.Format(FormatDateTS14)
}

// Time13ToTS8 13位时间戳转日期Ts8位,例如20161114
func Time13ToTS8(t int64) string {
	unix := time.Unix(t/1e3, 0)
	return unix.Format(FormatDateTS8)
}

// DateTs14ToTime13 日期Ts14位转13位时间戳,例如20161114143116转13位时间戳
func DateTs14ToTime13(t string) int64 {
	tm, _ := time.ParseInLocation(FormatDateTS14, t, LOC)
	return tm.UnixNano() / 1e6
}

// Date19ToTime13 日期19位转13位时间戳,例如2016-11-14 14:31:16转13位时间戳
func Date19ToTime13(t string) int64 {
	tm, _ := time.Parse(FormatDateTime, t)
	return tm.UnixNano() / 1e6
}

// Time13ToTime 13位时间戳转Time
func Time13ToTime(t int64) time.Time {
	return time.Unix(t/1e3, 0)
}

// TimeToTime13 Time转13位时间戳
func TimeToTime13(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// Time10ToDate 10位时间戳转日期
func Time10ToDate(t int64) string {
	unix := time.Unix(t, 0)
	return unix.Format(FormatDateTime)
}

// Time10ToTime 10位时间戳转Time
func Time10ToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

// TimeToTime10 Time转10位时间戳
func TimeToTime10(t time.Time) int64 {
	return t.UnixNano() / 1e9
}

// CurrentDateTime 获取当前日期与时间,("2006-01-02 15:04:05")格式
func CurrentDateTime() string {
	return time.Now().Format(FormatDateTime)
}

// CurrentDate 获取当前日期,("2006-01-02")格式
func CurrentDate() string {
	return time.Now().Format(FormatDate)
}

// AddDatesToTime10 通过增加days返回days天后的10位时间戳,days可以是负数
func AddDatesToTime10(days int) int64 {
	return time.Now().AddDate(0, 0, days).UnixNano() / 1e9
}

// AddDatesToTime13 通过增加days返回days天后的13位时间戳,days可以是负数
func AddDatesToTime13(days int) int64 {
	return time.Now().AddDate(0, 0, days).UnixNano() / 1e6
}

// AddDatesToDateTime 通过增加days返回days天后的日期与时间,days可以是负数
func AddDatesToDateTime(days int) string {
	return time.Now().AddDate(0, 0, days).Format(FormatDateTime)
}

// AddDatesToDate 通过增加days返回days天后的日期与时间,days可以是负数
func AddDatesToDate(days int) string {
	return time.Now().AddDate(0, 0, days).Format(FormatDate)
}

// AddDaysToZeroTime10 通过增加days返回days天后 00:00:00 时的10位时间戳,days可以是负数
func AddDaysToZeroTime10(days int) int64 {
	tm := time.Now().AddDate(0, 0, days)
	hour := int64(tm.Hour())
	min := int64(tm.Minute())
	sec := int64(tm.Second())
	return tm.Unix() - hour*3600 - min*60 - sec
}

// TimeZoneParse, 返回 1970-01-01 08:00:00 后加 seconds秒之后的日期与时间
func TimeZoneParse(seconds int64) string {
	return time.Unix(seconds, 0).Format("2006-01-02 15:04:05")
}

// 回退当日自定时时间(秒)
func Tm11FmtTz(sec int64, Fmt string) int64 {
	date := TimeZoneParse(sec)
	date = date[:11] + Fmt
	time1, _ := time.Parse(FormatDateTime, date)
	time1 = time1.UTC()
	tm := time1.UnixNano()/1e9 - 8*3600
	return tm
}

// RFC3339ToDateTime RFC3339格式数据转日期与时间
func RFC3339ToDateTime(t string) (string, error) {
	a, e := time.Parse(time.RFC3339, t)
	if e != nil {
		return "", e
	}
	return a.Format("2006-01-02 15:04:05"), nil
}

// Tim13ToYearMonthSlice 输入开始和结束的13位时间戳,输出其包含的年份和月份数组
func Tim13ToYearMonthSlice(startTime int64, endTime int64) []int {
	start, _ := strconv.Atoi(time.Unix(startTime/1e3, 0).Format("200601"))
	end, _ := strconv.Atoi(time.Unix(endTime/1e3, 0).Format("200601"))
	var month []int
	for i := start; i <= end; i++ {
		str := strconv.Itoa(i)
		sAtoi, _ := strconv.Atoi(str[len(str)-2 : len(str)])
		if sAtoi > 12 {
			i += 87
			continue
		}
		strAtoi, _ := strconv.Atoi(str)
		month = append(month, strAtoi)
	}
	return month
}
