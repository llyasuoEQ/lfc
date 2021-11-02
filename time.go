package lfc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MILLI  = 1e6
	SECOND = 1e9

	PRECISION             = MILLI                             // frequency control supported to time granularity
	RULE_PERIOD_PRECISION = MILLI                             // sliding window supported to milliseconds
	RATE                  = RULE_PERIOD_PRECISION / PRECISION // scale
)

func TimeStamp(times ...time.Time) (result int64) {
	if len(times) > 0 {
		result = times[0].UnixNano() / int64(PRECISION)
	} else {
		result = time.Now().UnixNano() / int64(PRECISION)
	}
	return
}

func getLowerBound(upperBound int64, period int64) (result int64) {
	result = upperBound - convertTimeAccuracy(period)
	return
}

func convertTimeAccuracy(period int64) int64 {
	return period * RATE
}

// timeStringToMilliSecond...
func timeStringToMilliSecond(ts string) (m int64, err error) {
	// 字符串拆分
	var numList []string
	var unitList []string

	numRe := regexp.MustCompile("[0-9]+")
	numList = numRe.FindAllString(ts, -1)

	unitRe := regexp.MustCompile("[^0-9]+")
	unitList = unitRe.FindAllString(ts, -1)

	if len(numList) == 0 || len(unitList) == 0 || len(numList) != len(unitList) {
		err = errors.New("time format to millisecond error")
		return
	}

	for i, unit := range unitList {
		var timeNum int64
		timeNum, err = strconv.ParseInt(numList[i], 10, 64)
		if err != nil {
			return
		}
		switch {
		case strings.EqualFold(unit, "ms"):
			m += timeNum
		case strings.EqualFold(unit, "s"):
			m += timeNum * 1000
		case strings.EqualFold(unit, "min"):
			m += timeNum * 60000
		case strings.EqualFold(unit, "h"):
			m += timeNum * 3600000
		case strings.EqualFold(unit, "d"):
			m += timeNum * 86400000
		default:
			err = errors.New("time format to millisecond error")
			return
		}
	}

	return
}
