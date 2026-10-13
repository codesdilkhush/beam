package backend

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConvertResult struct {
	Success   bool   `json:"success"`
	Query     string `json:"query"`
	FromValue string `json:"fromValue"`
	FromUnit  string `json:"fromUnit"`
	ToValue   string `json:"toValue"`
	ToUnit    string `json:"toUnit"`
	Result    string `json:"result"`
	Category  string `json:"category"` // "unit", "currency", "color", "time", "css"
	Formatted string `json:"formatted"`
}

// Length in meters
var lengthRates = map[string]float64{
	"km": 1000.0, "kilometer": 1000.0, "kilometers": 1000.0,
	"m": 1.0, "meter": 1.0, "meters": 1.0,
	"cm": 0.01, "centimeter": 0.01, "centimeters": 0.01,
	"mm": 0.001, "millimeter": 0.001, "millimeters": 0.001,
	"mi": 1609.344, "mile": 1609.344, "miles": 1609.344,
	"yd": 0.9144, "yard": 0.9144, "yards": 0.9144,
	"ft": 0.3048, "foot": 0.3048, "feet": 0.3048,
	"in": 0.0254, "inch": 0.0254, "inches": 0.0254,
}

// Mass in grams
var massRates = map[string]float64{
	"kg": 1000.0, "kilogram": 1000.0, "kilograms": 1000.0,
	"g": 1.0, "gram": 1.0, "grams": 1.0,
	"mg": 0.001, "milligram": 0.001, "milligrams": 0.001,
	"lb": 453.59237, "pound": 453.59237, "pounds": 453.59237, "lbs": 453.59237,
	"oz": 28.349523, "ounce": 28.349523, "ounces": 28.349523,
	"ton": 1000000.0, "tonne": 1000000.0, "tons": 1000000.0,
}

// Data storage in bytes (base 1024)
var dataRates = map[string]float64{
	"b": 1.0, "byte": 1.0, "bytes": 1.0,
	"kb": 1024.0, "kilobyte": 1024.0, "kilobytes": 1024.0,
	"mb": 1048576.0, "megabyte": 1048576.0, "megabytes": 1048576.0,
	"gb": 1073741824.0, "gigabyte": 1073741824.0, "gigabytes": 1073741824.0,
	"tb": 1099511627776.0, "terabyte": 1099511627776.0, "terabytes": 1099511627776.0,
	"pb": 1125899906842624.0, "petabyte": 1125899906842624.0,
}

// Speed in meters/sec
var speedRates = map[string]float64{
	"m/s": 1.0, "mps": 1.0,
	"km/h": 0.277778, "kmh": 0.277778, "kph": 0.277778,
	"mph":  0.44704,
	"knot": 0.514444, "knots": 0.514444,
}

// Time in seconds
var timeRates = map[string]float64{
	"s": 1.0, "sec": 1.0, "second": 1.0, "seconds": 1.0,
	"m": 60.0, "min": 60.0, "minute": 60.0, "minutes": 60.0,
	"h": 3600.0, "hr": 3600.0, "hour": 3600.0, "hours": 3600.0,
	"d": 86400.0, "day": 86400.0, "days": 86400.0,
	"wk": 604800.0, "week": 604800.0, "weeks": 604800.0,
	"mo": 2592000.0, "month": 2592000.0, "months": 2592000.0,
	"yr": 31536000.0, "year": 31536000.0, "years": 31536000.0,
}

// Major currencies relative to 1 USD
var currencyRates = map[string]float64{
	"usd": 1.0,
	"eur": 0.92,
	"gbp": 0.78,
	"jpy": 152.0,
	"cad": 1.38,
	"aud": 1.51,
	"inr": 83.5,
	"chf": 0.88,
	"cny": 7.24,
	"cbr": 1.0,
}

var (
	convertRegex = regexp.MustCompile(`^(?i)(?:convert\s+)?([+-]?\d+(?:\.\d+)?)\s*([a-zA-Z/°%]+)\s+(?:to|in|as)\s+([a-zA-Z/°%]+)$`)
	hexRegex     = regexp.MustCompile(`^(?i)#?([0-9a-f]{3}|[0-9a-f]{6})$`)
	rgbRegex     = regexp.MustCompile(`^(?i)rgb\s*\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$`)
)

func Convert(raw string) ConvertResult {
	query := strings.TrimSpace(raw)
	if query == "" {
		return ConvertResult{Success: false}
	}

	// 1. Color Hex -> RGB / HSL
	if strings.HasPrefix(query, "#") || (len(query) == 6 && isHex(query)) || (len(query) == 3 && isHex(query)) {
		if res, ok := convertHexColor(query); ok {
			return res
		}
	}

	// 2. Color RGB -> Hex / HSL
	if strings.HasPrefix(strings.ToLower(query), "rgb(") {
		if res, ok := convertRGBColor(query); ok {
			return res
		}
	}

	// 3. Calendar & Date queries (days in march, days in current month, what is the weekday, etc.)
	if res, ok := convertCalendarOrDate(query); ok {
		return res
	}

	// 4. Epoch/Unix time conversion
	if res, ok := convertEpochTime(query); ok {
		return res
	}

	// 4. Standard "<value> <unit1> to <unit2>" pattern
	matches := convertRegex.FindStringSubmatch(query)
	if len(matches) == 4 {
		valStr := matches[1]
		fromU := strings.ToLower(strings.TrimSpace(matches[2]))
		toU := strings.ToLower(strings.TrimSpace(matches[3]))

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return ConvertResult{Success: false}
		}

		// CSS Units: px <-> rem (16px base)
		if (fromU == "px" && toU == "rem") || (fromU == "rem" && toU == "px") {
			return convertCSSUnits(val, fromU, toU)
		}

		// Temperature: C, F, K
		if isTempUnit(fromU) && isTempUnit(toU) {
			return convertTemperature(val, fromU, toU)
		}

		// Length
		if r1, ok1 := lengthRates[fromU]; ok1 {
			if r2, ok2 := lengthRates[toU]; ok2 {
				resVal := (val * r1) / r2
				return formatStandardResult(val, fromU, resVal, toU, "length", "unit")
			}
		}

		// Mass
		if r1, ok1 := massRates[fromU]; ok1 {
			if r2, ok2 := massRates[toU]; ok2 {
				resVal := (val * r1) / r2
				return formatStandardResult(val, fromU, resVal, toU, "mass", "unit")
			}
		}

		// Data
		if r1, ok1 := dataRates[fromU]; ok1 {
			if r2, ok2 := dataRates[toU]; ok2 {
				resVal := (val * r1) / r2
				return formatStandardResult(val, fromU, resVal, toU, "data", "unit")
			}
		}

		// Speed
		if r1, ok1 := speedRates[fromU]; ok1 {
			if r2, ok2 := speedRates[toU]; ok2 {
				resVal := (val * r1) / r2
				return formatStandardResult(val, fromU, resVal, toU, "speed", "unit")
			}
		}

		// Time Duration
		if r1, ok1 := timeRates[fromU]; ok1 {
			if r2, ok2 := timeRates[toU]; ok2 {
				resVal := (val * r1) / r2
				return formatStandardResult(val, fromU, resVal, toU, "time", "unit")
			}
		}

		// Currency
		if r1, ok1 := currencyRates[fromU]; ok1 {
			if r2, ok2 := currencyRates[toU]; ok2 {
				// (val in USD) = val / r1; then in toU = (val / r1) * r2
				resVal := (val / r1) * r2
				return formatCurrencyResult(val, strings.ToUpper(fromU), resVal, strings.ToUpper(toU))
			}
		}
	}

	return ConvertResult{Success: false}
}

func isTempUnit(u string) bool {
	u = strings.TrimPrefix(u, "°")
	return u == "c" || u == "celsius" || u == "f" || u == "fahrenheit" || u == "k" || u == "kelvin"
}

func convertTemperature(val float64, fromU, toU string) ConvertResult {
	fromU = strings.TrimPrefix(fromU, "°")
	toU = strings.TrimPrefix(toU, "°")

	// Convert fromU to Celsius
	var celsius float64
	switch fromU {
	case "c", "celsius":
		celsius = val
	case "f", "fahrenheit":
		celsius = (val - 32.0) * (5.0 / 9.0)
	case "k", "kelvin":
		celsius = val - 273.15
	}

	// Convert Celsius to toU
	var res float64
	var targetLabel string
	switch toU {
	case "c", "celsius":
		res = celsius
		targetLabel = "°C"
	case "f", "fahrenheit":
		res = (celsius * 9.0 / 5.0) + 32.0
		targetLabel = "°F"
	case "k", "kelvin":
		res = celsius + 273.15
		targetLabel = "K"
	}

	formatted := fmt.Sprintf("%.2f %s = %.2f %s", val, strings.ToUpper(fromU), res, targetLabel)
	resStr := fmt.Sprintf("%.2f %s", res, targetLabel)

	return ConvertResult{
		Success:   true,
		Query:     fmt.Sprintf("%g %s to %s", val, fromU, toU),
		FromValue: fmt.Sprintf("%g", val),
		FromUnit:  fromU,
		ToValue:   fmt.Sprintf("%.2f", res),
		ToUnit:    targetLabel,
		Result:    resStr,
		Category:  "temperature",
		Formatted: formatted,
	}
}

func convertCSSUnits(val float64, fromU, toU string) ConvertResult {
	var res float64
	if fromU == "px" && toU == "rem" {
		res = val / 16.0
	} else {
		res = val * 16.0
	}

	resStr := fmt.Sprintf("%.4g%s", res, toU)
	formatted := fmt.Sprintf("%g%s = %s", val, fromU, resStr)

	return ConvertResult{
		Success:   true,
		Query:     fmt.Sprintf("%g%s to %s", val, fromU, toU),
		FromValue: fmt.Sprintf("%g", val),
		FromUnit:  fromU,
		ToValue:   resStr,
		ToUnit:    toU,
		Result:    resStr,
		Category:  "css",
		Formatted: formatted,
	}
}

func formatFloat(val float64) string {
	if math.Abs(val) >= 1e9 || (math.Abs(val) > 0 && math.Abs(val) < 1e-4) {
		return fmt.Sprintf("%.4g", val)
	}
	if val == math.Trunc(val) {
		return fmt.Sprintf("%.0f", val)
	}
	s := fmt.Sprintf("%.4f", val)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func formatStandardResult(fromVal float64, fromU string, toVal float64, toU, cat, mainCat string) ConvertResult {
	formattedTo := formatFloat(toVal)
	formattedFrom := formatFloat(fromVal)
	resStr := fmt.Sprintf("%s %s", formattedTo, toU)
	formatted := fmt.Sprintf("%s %s = %s", formattedFrom, fromU, resStr)

	return ConvertResult{
		Success:   true,
		Query:     fmt.Sprintf("%s %s to %s", formattedFrom, fromU, toU),
		FromValue: formattedFrom,
		FromUnit:  fromU,
		ToValue:   formattedTo,
		ToUnit:    toU,
		Result:    resStr,
		Category:  cat,
		Formatted: formatted,
	}
}

func formatCurrencyResult(fromVal float64, fromU string, toVal float64, toU string) ConvertResult {
	resStr := fmt.Sprintf("%.2f %s", toVal, toU)
	formatted := fmt.Sprintf("%.2f %s = %s", fromVal, fromU, resStr)

	return ConvertResult{
		Success:   true,
		Query:     fmt.Sprintf("%.2f %s to %s", fromVal, fromU, toU),
		FromValue: fmt.Sprintf("%.2f", fromVal),
		FromUnit:  fromU,
		ToValue:   fmt.Sprintf("%.2f", toVal),
		ToUnit:    toU,
		Result:    resStr,
		Category:  "currency",
		Formatted: formatted,
	}
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func convertHexColor(query string) (ConvertResult, bool) {
	clean := strings.TrimPrefix(query, "#")
	if len(clean) == 3 {
		clean = string([]byte{clean[0], clean[0], clean[1], clean[1], clean[2], clean[2]})
	}
	if len(clean) != 6 || !isHex(clean) {
		return ConvertResult{}, false
	}

	r, err1 := strconv.ParseInt(clean[0:2], 16, 64)
	g, err2 := strconv.ParseInt(clean[2:4], 16, 64)
	b, err3 := strconv.ParseInt(clean[4:6], 16, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return ConvertResult{}, false
	}

	rgbStr := fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
	h, s, l := rgbToHSL(int(r), int(g), int(b))
	hslStr := fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l)
	formatted := fmt.Sprintf("#%s = %s • %s", strings.ToUpper(clean), rgbStr, hslStr)

	return ConvertResult{
		Success:   true,
		Query:     query,
		FromValue: "#" + strings.ToUpper(clean),
		FromUnit:  "hex",
		ToValue:   rgbStr,
		ToUnit:    "rgb",
		Result:    rgbStr,
		Category:  "color",
		Formatted: formatted,
	}, true
}

func convertRGBColor(query string) (ConvertResult, bool) {
	matches := rgbRegex.FindStringSubmatch(query)
	if len(matches) != 4 {
		return ConvertResult{}, false
	}

	r, _ := strconv.Atoi(matches[1])
	g, _ := strconv.Atoi(matches[2])
	b, _ := strconv.Atoi(matches[3])
	if r > 255 || g > 255 || b > 255 {
		return ConvertResult{}, false
	}

	hexStr := fmt.Sprintf("#%02X%02X%02X", r, g, b)
	h, s, l := rgbToHSL(r, g, b)
	hslStr := fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l)
	formatted := fmt.Sprintf("rgb(%d, %d, %d) = %s • %s", r, g, b, hexStr, hslStr)

	return ConvertResult{
		Success:   true,
		Query:     query,
		FromValue: fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		FromUnit:  "rgb",
		ToValue:   hexStr,
		ToUnit:    "hex",
		Result:    hexStr,
		Category:  "color",
		Formatted: formatted,
	}, true
}

func rgbToHSL(r, g, b int) (int, int, int) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	l := (max + min) / 2.0
	var h, s float64

	if delta == 0 {
		h = 0
		s = 0
	} else {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2.0 - max - min)
		}

		if max == rf {
			h = (gf - bf) / delta
			if gf < bf {
				h += 6.0
			}
		} else if max == gf {
			h = (bf-rf)/delta + 2.0
		} else {
			h = (rf-gf)/delta + 4.0
		}
		h *= 60.0
	}

	return int(math.Round(h)), int(math.Round(s * 100.0)), int(math.Round(l * 100.0))
}

func convertEpochTime(query string) (ConvertResult, bool) {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "now to epoch" || q == "now to unix" || q == "epoch now" || q == "unix now" {
		now := time.Now()
		epoch := now.Unix()
		formatted := fmt.Sprintf("Now = %d (epoch seconds)", epoch)
		return ConvertResult{
			Success:   true,
			Query:     query,
			FromValue: "now",
			FromUnit:  "time",
			ToValue:   strconv.FormatInt(epoch, 10),
			ToUnit:    "epoch",
			Result:    strconv.FormatInt(epoch, 10),
			Category:  "time",
			Formatted: formatted,
		}, true
	}

	// Match "epoch <num>" or "<num> to date" or "unix <num>"
	var numStr string
	if strings.HasPrefix(q, "epoch ") {
		numStr = strings.TrimSpace(strings.TrimPrefix(q, "epoch "))
	} else if strings.HasPrefix(q, "unix ") {
		numStr = strings.TrimSpace(strings.TrimPrefix(q, "unix "))
	} else if strings.HasSuffix(q, " to date") {
		numStr = strings.TrimSpace(strings.TrimSuffix(q, " to date"))
	}

	if numStr != "" {
		sec, err := strconv.ParseInt(numStr, 10, 64)
		if err == nil && sec > 0 && sec < 4102444800 { // between 1970 and 2100
			t := time.Unix(sec, 0)
			utcStr := t.UTC().Format("2006-01-02 15:04:05 UTC")
			localStr := t.Local().Format("2006-01-02 15:04:05 MST")
			formatted := fmt.Sprintf("Epoch %d = %s (%s)", sec, utcStr, localStr)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: numStr,
				FromUnit:  "epoch",
				ToValue:   utcStr,
				ToUnit:    "date",
				Result:    utcStr,
				Category:  "time",
				Formatted: formatted,
			}, true
		}
	}

	return ConvertResult{}, false
}

var monthNames = map[string]time.Month{
	"january": time.January, "jan": time.January,
	"february": time.February, "feb": time.February,
	"march": time.March, "mar": time.March,
	"april": time.April, "apr": time.April,
	"may":  time.May,
	"june": time.June, "jun": time.June,
	"july": time.July, "jul": time.July,
	"august": time.August, "aug": time.August,
	"september": time.September, "sep": time.September, "sept": time.September,
	"october": time.October, "oct": time.October,
	"november": time.November, "nov": time.November,
	"december": time.December, "dec": time.December,
}

func convertCalendarOrDate(query string) (ConvertResult, bool) {
	q := strings.ToLower(strings.TrimSpace(query))
	q = strings.Trim(q, "?!.,")
	q = strings.TrimSpace(q)
	if q == "" {
		return ConvertResult{}, false
	}

	now := time.Now()

	// 1. Weekday queries: "weekday", "what is the weekday", "what day is today", "what day is it", etc.
	if q == "weekday" || q == "weekday today" || q == "what is the weekday" || q == "what is the weekday today" ||
		q == "what day is today" || q == "what day is it" || q == "what day is it today" ||
		q == "day of the week" || q == "day of week" || q == "today weekday" || q == "current weekday" {
		weekday := now.Weekday().String()
		formatted := fmt.Sprintf("Today is %s (%s)", weekday, now.Format("January 2, 2006"))
		return ConvertResult{
			Success:   true,
			Query:     query,
			FromValue: "today",
			FromUnit:  "weekday",
			ToValue:   weekday,
			ToUnit:    "weekday",
			Result:    weekday,
			Category:  "calendar",
			Formatted: formatted,
		}, true
	}

	// Weekday of specific date: "weekday of <date>", "weekday on <date>"
	var dateTarget string
	if strings.HasPrefix(q, "weekday of ") {
		dateTarget = strings.TrimPrefix(q, "weekday of ")
	} else if strings.HasPrefix(q, "weekday on ") {
		dateTarget = strings.TrimPrefix(q, "weekday on ")
	}
	if dateTarget != "" {
		for _, layout := range []string{"2006-01-02", "01/02/2006", "Jan 2 2006", "January 2 2006", "2 Jan 2006", "2 January 2006"} {
			if parsedT, err := time.Parse(layout, dateTarget); err == nil {
				weekday := parsedT.Weekday().String()
				formatted := fmt.Sprintf("%s is %s", parsedT.Format("January 2, 2006"), weekday)
				return ConvertResult{
					Success:   true,
					Query:     query,
					FromValue: dateTarget,
					FromUnit:  "date",
					ToValue:   weekday,
					ToUnit:    "weekday",
					Result:    weekday,
					Category:  "calendar",
					Formatted: formatted,
				}, true
			}
		}
	}

	// 2. Today / current date
	if q == "today" || q == "date today" || q == "today date" || q == "current date" ||
		q == "what is today" || q == "what is the date" || q == "what is today's date" || q == "what is todays date" {
		dateStr := now.Format("Monday, January 2, 2006")
		return ConvertResult{
			Success:   true,
			Query:     query,
			FromValue: "today",
			FromUnit:  "date",
			ToValue:   now.Format("2006-01-02"),
			ToUnit:    "date",
			Result:    dateStr,
			Category:  "calendar",
			Formatted: dateStr,
		}, true
	}

	// 3. Days in month / year
	normalized := strings.TrimPrefix(q, "how many ")
	normalized = strings.TrimPrefix(normalized, "number of ")
	if strings.HasPrefix(normalized, "days in ") {
		rest := strings.TrimSpace(strings.TrimPrefix(normalized, "days in "))

		// "days in current month", "days in this month"
		if rest == "current month" || rest == "this month" || rest == "month" {
			m := now.Month()
			y := now.Year()
			days := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
			resStr := fmt.Sprintf("%d days", days)
			formatted := fmt.Sprintf("%s %d has %d days", m.String(), y, days)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: fmt.Sprintf("%s %d", m.String(), y),
				FromUnit:  "month",
				ToValue:   strconv.Itoa(days),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		// "days in next month"
		if rest == "next month" {
			nextM := now.AddDate(0, 1, 0)
			m := nextM.Month()
			y := nextM.Year()
			days := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
			resStr := fmt.Sprintf("%d days", days)
			formatted := fmt.Sprintf("%s %d has %d days", m.String(), y, days)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: fmt.Sprintf("%s %d", m.String(), y),
				FromUnit:  "month",
				ToValue:   strconv.Itoa(days),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		// "days in last month", "days in previous month"
		if rest == "last month" || rest == "previous month" {
			prevM := now.AddDate(0, -1, 0)
			m := prevM.Month()
			y := prevM.Year()
			days := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
			resStr := fmt.Sprintf("%d days", days)
			formatted := fmt.Sprintf("%s %d has %d days", m.String(), y, days)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: fmt.Sprintf("%s %d", m.String(), y),
				FromUnit:  "month",
				ToValue:   strconv.Itoa(days),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		// "days in current year", "days in this year", "days in year"
		if rest == "current year" || rest == "this year" || rest == "year" {
			y := now.Year()
			isLeap := (y%4 == 0 && y%100 != 0) || (y%400 == 0)
			days := 365
			leapStr := ""
			if isLeap {
				days = 366
				leapStr = " (leap year)"
			}
			resStr := fmt.Sprintf("%d days", days)
			formatted := fmt.Sprintf("Year %d has %d days%s", y, days, leapStr)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: strconv.Itoa(y),
				FromUnit:  "year",
				ToValue:   strconv.Itoa(days),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		// Check if rest is a 4-digit year like "2024"
		if y, err := strconv.Atoi(rest); err == nil && y > 0 && y < 10000 {
			isLeap := (y%4 == 0 && y%100 != 0) || (y%400 == 0)
			days := 365
			leapStr := ""
			if isLeap {
				days = 366
				leapStr = " (leap year)"
			}
			resStr := fmt.Sprintf("%d days", days)
			formatted := fmt.Sprintf("Year %d has %d days%s", y, days, leapStr)
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: strconv.Itoa(y),
				FromUnit:  "year",
				ToValue:   strconv.Itoa(days),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		// Check if rest is a month name, optionally followed by a year (e.g. "march", "february 2024", "mar")
		fields := strings.Fields(rest)
		if len(fields) >= 1 {
			monthStr := fields[0]
			if m, ok := monthNames[monthStr]; ok {
				y := now.Year()
				yearSpecified := false
				if len(fields) >= 2 {
					if parsedY, err := strconv.Atoi(fields[1]); err == nil && parsedY > 0 && parsedY < 10000 {
						y = parsedY
						yearSpecified = true
					}
				}
				days := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
				resStr := fmt.Sprintf("%d days", days)
				var formatted string
				if yearSpecified {
					isLeap := (y%4 == 0 && y%100 != 0) || (y%400 == 0)
					leapNote := ""
					if m == time.February && isLeap {
						leapNote = " (leap year)"
					}
					formatted = fmt.Sprintf("%s %d has %d days%s", m.String(), y, days, leapNote)
				} else {
					formatted = fmt.Sprintf("%s has %d days", m.String(), days)
				}

				return ConvertResult{
					Success:   true,
					Query:     query,
					FromValue: m.String(),
					FromUnit:  "month",
					ToValue:   strconv.Itoa(days),
					ToUnit:    "days",
					Result:    resStr,
					Category:  "calendar",
					Formatted: formatted,
				}, true
			}
		}
	}

	// 4. Days remaining: "days left in month", "days left in year"
	if strings.HasPrefix(normalized, "days left in ") || strings.HasPrefix(normalized, "days remaining in ") {
		rest := strings.TrimSpace(strings.TrimPrefix(normalized, "days left in "))
		rest = strings.TrimSpace(strings.TrimPrefix(rest, "days remaining in "))

		if rest == "month" || rest == "this month" || rest == "current month" {
			daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			remaining := daysInMonth - now.Day()
			resStr := fmt.Sprintf("%d days", remaining)
			formatted := fmt.Sprintf("%d days remaining in %s %d", remaining, now.Month().String(), now.Year())
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: fmt.Sprintf("%s %d", now.Month().String(), now.Year()),
				FromUnit:  "month",
				ToValue:   strconv.Itoa(remaining),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}

		if rest == "year" || rest == "this year" || rest == "current year" {
			dayOfYear := now.YearDay()
			isLeap := (now.Year()%4 == 0 && now.Year()%100 != 0) || (now.Year()%400 == 0)
			total := 365
			if isLeap {
				total = 366
			}
			remaining := total - dayOfYear
			resStr := fmt.Sprintf("%d days", remaining)
			formatted := fmt.Sprintf("%d days remaining in %d", remaining, now.Year())
			return ConvertResult{
				Success:   true,
				Query:     query,
				FromValue: strconv.Itoa(now.Year()),
				FromUnit:  "year",
				ToValue:   strconv.Itoa(remaining),
				ToUnit:    "days",
				Result:    resStr,
				Category:  "calendar",
				Formatted: formatted,
			}, true
		}
	}

	// 5. Week number: "week number", "current week"
	if q == "week number" || q == "current week" || q == "week of year" || q == "what week is it" {
		year, week := now.ISOWeek()
		result := fmt.Sprintf("Week %d", week)
		formatted := fmt.Sprintf("Week %d of %d", week, year)
		return ConvertResult{
			Success:   true,
			Query:     query,
			FromValue: "now",
			FromUnit:  "week",
			ToValue:   strconv.Itoa(week),
			ToUnit:    "week",
			Result:    result,
			Category:  "calendar",
			Formatted: formatted,
		}, true
	}

	// 6. Day of year: "day of year", "day of the year"
	if q == "day of year" || q == "day of the year" || q == "current day of year" {
		dayOfYear := now.YearDay()
		isLeap := (now.Year()%4 == 0 && now.Year()%100 != 0) || (now.Year()%400 == 0)
		total := 365
		if isLeap {
			total = 366
		}
		result := fmt.Sprintf("Day %d", dayOfYear)
		pct := float64(dayOfYear) / float64(total) * 100.0
		formatted := fmt.Sprintf("Day %d of %d (%.1f%% of %d)", dayOfYear, total, pct, now.Year())
		return ConvertResult{
			Success:   true,
			Query:     query,
			FromValue: "today",
			FromUnit:  "day",
			ToValue:   strconv.Itoa(dayOfYear),
			ToUnit:    "day",
			Result:    result,
			Category:  "calendar",
			Formatted: formatted,
		}, true
	}

	return ConvertResult{}, false
}
