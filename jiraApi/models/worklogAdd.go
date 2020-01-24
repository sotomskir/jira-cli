package models

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// WorklogAdd represents request to JIRA API to add worklog
type WorklogAdd struct {
	Comment          string `json:"comment"`
	Started          string `json:"started,omitempty"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
}

// Prepare worklog from user data
func InitilizeWorklogAdd(com string, workedSec uint64, d string, t string) (WorklogAdd, error) {
	w := WorklogAdd{Comment: com, TimeSpentSeconds: workedSec * 60}
	if d != "" || t != "" {
		n := time.Now().Format("-0700")
		ut, e := time.Parse("2006-01-02 15:04 -0700", d+" "+t+" "+n)
		if e != nil {
			msg := fmt.Sprintln("If provided the date and time must adhere to formats: [YYYY-MM-DD] and [HH:ss]. You provided: date=[", d, "] and time=[", t, "]")
			return WorklogAdd{}, errors.New(msg)
		}
		//fraction of seconds is really important https://jira.atlassian.com/browse/JRASERVER-61378
		w.Started = ut.Format("2006-01-02T15:04:05.000-0700")
	}
	return w, nil
}
