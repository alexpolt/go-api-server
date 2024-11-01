package common

import (
	"errors"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func Default_string(value string, default_value string) string {
	if value == "" {
		value = default_value
	}
	return value
}

func Log_Time(start time.Time) {
	Log("Time (us):", time.Since(start).Microseconds())
}

func Join_Errors(err ...error) error {
	out := ""
	for _, e := range err {
		if e != nil {
			out = out + "; " + e.Error()
		}
	}
	return errors.New(out)
}
