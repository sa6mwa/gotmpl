package main

// Majority copied from github.com/hairyhenderson/gomplate/v4/funcs
// (funcs/time.go) Copyright (c) 2016-2023 Dave Henderson

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sa6mwa/gotmpl/conv"
)

type TimeFuncs struct {
	ANSIC       string
	UnixDate    string
	RubyDate    string
	RFC822      string
	RFC822Z     string
	RFC850      string
	RFC1123     string
	RFC1123Z    string
	RFC3339     string
	RFC3339Nano string
	Kitchen     string
	Stamp       string
	StampMilli  string
	StampMicro  string
	StampNano   string
	DateTime    string
	DateOnly    string
	TimeOnly    string
}

func (TimeFuncs) Parse(layout string, value interface{}) (time.Time, error) {
	return time.Parse(layout, conv.ToString(value))
}

func (f TimeFuncs) ParseLocal(layout string, value interface{}) (time.Time, error) {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Local"
	}
	return f.ParseInLocation(layout, tz, value)
}

func (TimeFuncs) ParseInLocation(layout, location string, value interface{}) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(layout, conv.ToString(value), loc)
}

func (TimeFuncs) Now() time.Time {
	return time.Now()
}

func (TimeFuncs) Unix(in interface{}) (time.Time, error) {
	sec, nsec, err := parseNum(in)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, nsec), nil
}

func (TimeFuncs) Nanosecond(n interface{}) time.Duration {
	return time.Nanosecond * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) Microsecond(n interface{}) time.Duration {
	return time.Microsecond * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) Millisecond(n interface{}) time.Duration {
	return time.Millisecond * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) Second(n interface{}) time.Duration {
	return time.Second * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) Minute(n interface{}) time.Duration {
	return time.Minute * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) Hour(n interface{}) time.Duration {
	return time.Hour * time.Duration(conv.ToInt64(n))
}

func (TimeFuncs) ParseDuration(n interface{}) (time.Duration, error) {
	return time.ParseDuration(conv.ToString(n))
}

func (TimeFuncs) Since(n time.Time) time.Duration {
	return time.Since(n)
}

func (TimeFuncs) Until(n time.Time) time.Duration {
	return time.Until(n)
}

// convert a number input to a pair of int64s, representing the integer portion and the decimal remainder
// this can handle a string as well as any integer or float type
// precision is at the "nano" level (i.e. 1e+9)
func parseNum(in interface{}) (integral int64, fractional int64, err error) {
	if s, ok := in.(string); ok {
		ss := strings.Split(s, ".")
		if len(ss) > 2 {
			return 0, 0, fmt.Errorf("can not parse '%s' as a number - too many decimal points", s)
		}
		if len(ss) == 1 {
			integral, err := strconv.ParseInt(s, 0, 64)
			return integral, 0, err
		}
		integral, err := strconv.ParseInt(ss[0], 0, 64)
		if err != nil {
			return integral, 0, err
		}
		fractional, err = strconv.ParseInt(padRight(ss[1], "0", 9), 0, 64)
		return integral, fractional, err
	}
	if s, ok := in.(fmt.Stringer); ok {
		return parseNum(s.String())
	}
	if i, ok := in.(int); ok {
		return int64(i), 0, nil
	}
	if u, ok := in.(uint64); ok {
		return int64(u), 0, nil
	}
	if f, ok := in.(float64); ok {
		return 0, 0, fmt.Errorf("can not parse floating point number (%f) - use a string instead", f)
	}
	if in == nil {
		return 0, 0, nil
	}
	return 0, 0, nil
}

// pads a number with zeroes
func padRight(in, pad string, length int) string {
	for {
		in += pad
		if len(in) > length {
			return in[0:length]
		}
	}
}
