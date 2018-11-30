package calendar

import (
	"time"
	"fmt"
)

var (
	daysInMonth =[12]int {31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

type SyncTime struct {
	StartTime string
	EndTime string
	nextSync string
}

func GetYestoday(num int)  string{

	year:=time.Now().Year()
	month:=time.Now().Month()
	day:=time.Now().Day()

	if year%4 == 0 && year%100 != 0 {
		daysInMonth[2] = 29
	}

	if num >= day {
		if month ==1 {
			year = year-1
			month = 12
		}else {
			month = month-1
		}
		day = daysInMonth[month-1]+day-num
	}else{
		day = day-num
	}

	return fmt.Sprintf("%d-%d-%d",year,month,day)
}

func GetDays() int{
	weekday := time.Now().Weekday()
	if weekday == 1{
		return 3
	}else {
		return 1
	}
}


func TimeEqual(eqT string) bool{
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	var t string
	if minute <10 {
		t = fmt.Sprintf("%d:0%d",hour,minute)
	} else{
		t = fmt.Sprintf("%d:%d",hour,minute)
	}
	fmt.Println(eqT,t)
	if eqT == t{
		return true
	}else{
		return false
	}
}

func GetCurrMinuteTime() string{
	hour := time.Now().Hour()
	minute := time.Now().Minute()

	if minute <10 {
		return fmt.Sprintf("%d:0%d",hour,minute)
	} else{
		return fmt.Sprintf("%d:%d",hour,minute)
	}
}


