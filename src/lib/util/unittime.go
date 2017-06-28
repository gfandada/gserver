package util

/*************************************************
				sec   0-59
				min   0-59
				hour  0-23
				day   1-31
				month 1-12
				week  0-6
*************************************************/

// 定时器配置
type TimerCfg struct {
	Daily   []DailyTimer
	Weekly  []WeekTimer
	Monthly []MonthlyTimer
}

// 每日定时器
type DailyTimer struct {
	Hour uint
	Min  uint
	Sec  uint
}

// 每周定时器
type WeekTimer struct {
	Week uint
	Hour uint
	Min  uint
	Sec  uint
}

// 每月定时器
type MonthlyTimer struct {
	Month uint8
	Day   uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}
