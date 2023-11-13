package config

import "time"

var CacheTimeFiveSecond time.Duration
var CacheTimeTenSecond time.Duration
var CacheTimeHalfMinute time.Duration
var CacheTimeOneMinute time.Duration
var CacheTimeTwoMinute time.Duration
var CacheTimeFiveMinute time.Duration
var CacheTimeTenMinute time.Duration
var CacheTimeTwentyMinute time.Duration
var CacheTimeHalfHour time.Duration
var CacheTimeOneHour time.Duration
var CacheTimeSixHour time.Duration
var CacheTimeHalfDay time.Duration
var CacheTimeOneDay time.Duration
var CacheTimeOneWeek time.Duration
var CacheTimeOneMonth time.Duration
var CacheTimeOneYear time.Duration

func init() {
	CacheTimeFiveSecond = 5
	CacheTimeTenSecond = 10
	CacheTimeHalfMinute = 30
	CacheTimeOneMinute = 60
	CacheTimeTwoMinute = 120
	CacheTimeFiveMinute = 300
	CacheTimeTenMinute = 600
	CacheTimeTwentyMinute = 1200
	CacheTimeHalfHour = 1800
	CacheTimeOneHour = 3600
	CacheTimeSixHour = 21600
	CacheTimeHalfDay = 43200
	CacheTimeOneDay = 86400
	CacheTimeOneWeek = 604800
	CacheTimeOneMonth = 2592000
	CacheTimeOneYear = 31536000
}
