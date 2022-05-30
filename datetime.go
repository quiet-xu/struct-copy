package struct_copy

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	TIME_INDEX_YEAR = iota
	TIME_INDEX_MONTH
	TIME_INDEX_DAY
	TIME_INDEX_HOURS
	TIME_INDEX_MINUTE
	TIME_INDEX_SEC
	TIME_INDEX_MSEC
	TIME_INDEX_USEC
	TIME_INDEX_NSEC
)

const (
	TIME_SECONDS   = 1
	TIME_MINUTES   = TIME_SECONDS * 60
	TIME_HOURS     = TIME_MINUTES * 60
	TIME_DAY       = TIME_HOURS * 24
	TIME_WEEK      = TIME_DAY * 7
	TIME_HALFMONTH = TIME_DAY * 15
	TIME_FULLMONTH = TIME_DAY * 30
	TIME_QUARTER   = TIME_FULLMONTH * 3
	TIME_HALFYEAR  = TIME_QUARTER * 2
	TIME_FULLYEAR  = TIME_DAY * 365
)

type TimeFormatMode int

const (
	TIME_JOIN TimeFormatMode = iota
	TIME_RELATIVE
	TIME_ABSOLUTE
)

type Datetime time.Time

func (j Datetime) GobEncode() ([]byte, error) {
	return time.Time(j).MarshalBinary()
}

func (j *Datetime) GobDecode(data []byte) error {
	var t time.Time
	if e := t.UnmarshalBinary(data); nil != e {
		return e
	}
	*j = Datetime(t)
	return nil
}

func (j Datetime) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

func (j *Datetime) UnmarshalXLSX(buf []byte) error {
	if len(buf) < 4 {
		*j = Datetime{}
		return nil
	}
	var arr [9]int
	ParseDateTime(string(buf), TIME_ABSOLUTE, arr[:])
	*j = Datetime(ToTime(arr[:]))
	return nil
}

func (j *Datetime) UnmarshalJSON(buf []byte) (err error) {
	if len(buf) < 6 {
		*j = Datetime{}
		return
	}
	var arr [9]int
	ParseDateTime(string(buf[1:len(buf)-1]), TIME_ABSOLUTE, arr[:])
	*j = Datetime(ToTime(arr[:]))
	return
}

func (j Datetime) String() string {
	return time.Time(j).Format("2006-01-02 15:04:05")
}

func (j Datetime) StringEx() string {
	return time.Time(j).Format("20060102150405") + strconv.Itoa(int(time.Now().UnixNano()%100000))
}

func (j Datetime) AsTime() time.Time {
	return time.Time(j)
}
func (j Datetime) FromString(dateString string) (Datetime, error) {
	if len(dateString) <= 0 {
		return Datetime{}, nil
	}
	t, err := time.ParseInLocation("2006-1-2 15:04:05", dateString, time.Local)
	if err != nil {
		return Datetime{}, err
	}
	return Datetime(t), nil
}
func StringToDatetime(dateString string) Datetime {
	//Datetime{}的零值好像不带时区
	var t time.Time
	var err error
	if len(dateString) == 0 {
		goto Here
	}
	t, err = time.ParseInLocation("2006-1-2 15:04:05", dateString, time.Local)
	if err != nil {
		goto Here
	} else {
		return Datetime(t)
	}
Here:
	t, _ = time.ParseInLocation("2006-1-2 15:04:05", "0001-01-01 00:00:00", time.Local)
	return Datetime(t)
}

type Date time.Time

func (j Date) FromString(dateString string) (Date, error) {
	if len(dateString) <= 0 {
		return Date{}, nil
	}
	t, err := time.ParseInLocation("2006-1-2", dateString, time.Local)
	if err != nil {
		return Date{}, err
	}
	return Date(t), nil
}

func (j Date) GobEncode() ([]byte, error) {
	return time.Time(j).MarshalBinary()
}

func (j *Date) GobDecode(data []byte) error {
	var t time.Time
	if e := t.UnmarshalBinary(data); nil != e {
		return e
	}
	*j = Date(t)
	return nil
}

func (j Date) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(j).Format("2006-01-02") + `"`), nil
}

func (j *Date) UnmarshalXLSX(buf []byte) error {
	if len(buf) < 4 {
		*j = Date{}
		return nil
	}
	var arr [9]int
	ParseDateTime(string(buf), TIME_ABSOLUTE, arr[:])
	*j = Date(ToTime(arr[:]))
	return nil
}

func (j *Date) UnmarshalJSON(buf []byte) (err error) {
	if len(buf) < 6 {
		*j = Date{}
		return
	}
	var arr [9]int
	ParseDateTime(string(buf[1:len(buf)-1]), TIME_ABSOLUTE, arr[:])
	*j = Date(ToTime(arr[:]))
	return
}

func (j Date) String() string {
	return time.Time(j).Format("2006-01-02")
}

func (j Date) StringEx() string {
	return time.Time(j).Format("20060102150405") + strconv.Itoa(int(time.Now().UnixNano()%100000))
}

func (j Date) AsTime() time.Time {
	return time.Time(j)
}

type Timestamp int64

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t < 0 {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Unix(int64(t), 0).Format("2006-01-02 15:04:05") + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(buf []byte) (err error) {
	if len(buf) < 2 {
		return errors.New("illegal time format")
	}
	var arr [9]int
	ParseDateTime(string(buf[1:len(buf)-1]), TIME_ABSOLUTE, arr[:])
	*t = Timestamp(ToTime(arr[:]).Unix())
	return
}

func (t Timestamp) String() string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func AsDate(op1 time.Time) time.Time {
	y, m, d := op1.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, op1.Location())
}

func AsTime(op1 time.Time) time.Time {
	h, m, s := op1.Clock()
	return time.Date(0, 0, 0, h, m, s, 0, op1.Location())
}

func AsNowDateEx(args ...int) time.Time {
	y, m, d := time.Now().Date()
	var hh, mm, ss int

	switch len(args) {
	case 3:
		ss = args[2]
	case 2:
		mm = args[1]
	case 1:
		hh = args[0]
	}

	return time.Date(y, m, d, hh, mm, ss, 0, time.Local)
}

func AsNowDate(op1 time.Time) time.Time {
	y, m, d := time.Now().Date()
	hh, mm, ss := op1.Clock()
	return time.Date(y, m, d, hh, mm, ss, 0, op1.Location())
}

func AsNowTime(op1 time.Time) time.Time {
	y, m, d := op1.Date()
	hh, mm, ss := time.Now().Clock()
	return time.Date(y, m, d, hh, mm, ss, 0, op1.Location())
}

func EqualDate(op1 time.Time, op2 time.Time) bool {
	y1, m1, d1 := op1.Date()
	y2, m2, d2 := op2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ToTime(datetime []int) time.Time {
	return time.Date(
		datetime[TIME_INDEX_YEAR],
		time.Month(datetime[TIME_INDEX_MONTH]),
		datetime[TIME_INDEX_DAY],
		datetime[TIME_INDEX_HOURS],
		datetime[TIME_INDEX_MINUTE],
		datetime[TIME_INDEX_SEC]+datetime[TIME_INDEX_MSEC]/1000+datetime[TIME_INDEX_USEC]/1000000,
		datetime[TIME_INDEX_NSEC], time.Local)
}

func ToDuration(datetime []int) (d time.Duration) {
	d += time.Duration(datetime[TIME_INDEX_YEAR]) * TIME_FULLYEAR * time.Second
	d += time.Duration(datetime[TIME_INDEX_MONTH]) * TIME_FULLMONTH * time.Second
	d += time.Duration(datetime[TIME_INDEX_DAY]) * TIME_DAY * time.Second
	d += time.Duration(datetime[TIME_INDEX_HOURS]) * TIME_HOURS * time.Second
	d += time.Duration(datetime[TIME_INDEX_MINUTE]) * TIME_MINUTES * time.Second

	d += time.Duration(datetime[TIME_INDEX_SEC]) * time.Second
	d += time.Duration(datetime[TIME_INDEX_MSEC]) * time.Millisecond
	d += time.Duration(datetime[TIME_INDEX_USEC]) * time.Microsecond
	d += time.Duration(datetime[TIME_INDEX_NSEC]) * time.Nanosecond
	return
}

func TimeOf(datetime []int, tm time.Time) []int {
	if len(datetime) < 9 {
		datetime = make([]int, 9, 9)
	}
	cur := time.Now()
	var month time.Month
	datetime[TIME_INDEX_YEAR], month, datetime[TIME_INDEX_DAY] = cur.Date()
	datetime[TIME_INDEX_MONTH] = int(month)
	datetime[TIME_INDEX_HOURS], datetime[TIME_INDEX_MINUTE], datetime[TIME_INDEX_SEC] = cur.Clock()
	datetime[TIME_INDEX_NSEC] = 0 //cur.Nanosecond()
	return datetime
}

func NowTime(datetime []int) []int {
	return TimeOf(datetime, time.Now())
}

func ParseTime(str string, dir string, datetime []int) (count int) {
	if "" == str {
		return
	}
	arr := strings.Split(str, dir)

	n := len(datetime)
	if len(datetime) > len(arr) {
		n = len(arr)
	}

	for i := 0; i < n; i++ {
		if len(arr[i]) > 0 {
			datetime[i] = int(ParseInteger(arr[i], 0))
			count++
		}
	}
	return
}

func parseDateTime(str string, mode TimeFormatMode, datetime []int, now func([]int) []int) {
	pos := strings.IndexByte(str, 32)
	var d, t string

	if pos > 2 {
		if d = str[:pos]; pos < len(str) {
			t = str[pos+1:]
		}
	} else if strings.IndexByte(str, 45) < 0 {
		t = str
	} else {
		d = str
	}

	if TIME_JOIN == mode {
		now(datetime[:])
	}
	ParseTime(d, "-", datetime[:3])
	ParseTime(t, ":", datetime[3:])

	if TIME_RELATIVE == mode {
		cur := now(nil)
		for k, v := range cur {
			datetime[k] += v
		}
	}
	return
}

func ParseDateTime(str string, mode TimeFormatMode, datetime []int) {
	parseDateTime(str, mode, datetime, NowTime)
}

func AddTime(str string, tm time.Time, args ...TimeFormatMode) time.Time {
	var arr [9]int
	var mode TimeFormatMode
	if len(args) > 0 {
		mode = args[0]
	} else {
		mode = TIME_RELATIVE
	}

	parseDateTime(str, mode, arr[:], func(src []int) []int {
		return TimeOf(src, tm)
	})

	return ToTime(arr[:])
}

func ToDateTime(str string, args ...TimeFormatMode) time.Time {
	var arr [9]int
	var mode TimeFormatMode
	if len(args) > 0 {
		mode = args[0]
	} else {
		mode = TIME_ABSOLUTE
	}

	ParseDateTime(str, mode, arr[:])
	return ToTime(arr[:])
}

func ToDate(str string) time.Time {
	var arr [3]int
	if pos := strings.IndexByte(str, 32); pos > 2 {
		str = str[:pos]
	}
	ParseTime(str, "-", arr[:])
	return time.Date(
		arr[TIME_INDEX_YEAR],
		time.Month(arr[TIME_INDEX_MONTH]),
		arr[TIME_INDEX_DAY],
		0, 0, 0, 0, time.Local)
}

func (j Datetime) MonthDifferent(k Datetime) int {
	var highOne time.Time
	var lowOne time.Time
	switch {
	case j.AsTime().Before(k.AsTime()):
		highOne = k.AsTime()
		lowOne = j.AsTime()
	case j.AsTime().After(k.AsTime()):
		highOne = j.AsTime()
		lowOne = k.AsTime()
	default:
		return 0
	}
	return (highOne.Year()-lowOne.Year())*12 + int(highOne.Month()) - int(lowOne.Month())
}

func (j Datetime) YearDifferent(k Datetime) int {
	var highOne time.Time
	var lowOne time.Time
	switch {
	case j.AsTime().Before(k.AsTime()):
		highOne = k.AsTime()
		lowOne = j.AsTime()
	case j.AsTime().After(k.AsTime()):
		highOne = j.AsTime()
		lowOne = k.AsTime()
	default:
		return 0
	}
	var yearByMonth = float64(highOne.Month()-lowOne.Month()) / 12.0
	return (highOne.Year() - lowOne.Year()) + int(math.Round(yearByMonth))
}

//DayBegin 日初丨日末丨月初丨月末丨年初丨年末
func (j Datetime) DayBegin() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year(), xTime.Month(), xTime.Day(),
		0, 0, 0, 0, xTime.Location()))
}
func (j Datetime) DayEnd() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year(), xTime.Month(), xTime.Day(),
		23, 59, 59, 999999999, xTime.Location()))
}
func (j Datetime) MonthBegin() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year(), xTime.Month(), 1,
		0, 0, 0, 0, xTime.Location()))
}
func (j Datetime) MonthEnd() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year(), xTime.Month()+1, 0,
		23, 59, 59, 999999999, xTime.Location()))
}
func (j Datetime) YearBegin() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year(), 1, 1,
		0, 0, 0, 0, xTime.Location()))
}
func (j Datetime) YearEnd() Datetime {
	xTime := j.AsTime()
	return Datetime(time.Date(xTime.Year()+1, 1, 0,
		23, 59, 59, 999999999, xTime.Location()))
}

//IsZero 是否零值（综合）
func (j Datetime) IsZero() (result bool) {
	return j.AsTime().IsZero() ||
		j.String() == "0001-01-01 00:00:00" ||
		j.String() == "0001-01-01 08:00:00" || //数据库是0001-01-01 00:00:00 时  映射过来可能会变成 0001-01-01 08:00:00
		j.String() == "0000-00-00 00:00:00"
}
func (j Datetime) NotZero() (result bool) {
	return !j.IsZero()
}

func NowNowNow() Datetime {
	return Datetime(time.Now())
}

//StringReplaceWithEmpty 时间转string  如果是空的则出来的也是空的string
func (j Datetime) StringReplaceWithEmpty() (res string) {
	if j.IsZero() {
		return ""
	}
	res = j.String()
	return
}

func ParseInteger(s string, v int64) int64 {
	if len(s) > 0 && uint(s[0])-48 < 10 {
		v = 0
		for _, c := range s {
			u := uint(c) - 48
			if u > uint(9) {
				break
			}
			v = v*10 + int64(u)
		}
	}
	return v
}
