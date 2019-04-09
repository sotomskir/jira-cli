package models

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestInitilizeWorklogAdd(t *testing.T) {
	tz := time.Now().Format("-0700")
	exp := WorklogAdd{Comment: "Test", TimeSpentSeconds: 300, Started: "2019-04-02T20:21:00.000"+tz}
	payload, _ := InitilizeWorklogAdd("Test", 5, "2019-04-02", "20:21")
	assert.DeepEqual(t, payload, exp)
}

func TestInitilizeWorklogAddCrashOnValidation(t *testing.T) {
	_, e := InitilizeWorklogAdd("T", 1, "3", "dd")
	assert.Error(t, e, "If provided the date and time must adhere to formats: [YYYY-MM-DD] and [HH:ss]. You provided: date=[ 3 ] and time=[ dd ]\n")
}
