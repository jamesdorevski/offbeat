package client

import (
	"time"
)

type Worklog struct {
	Issue struct {
		Self string `json:"self"`
		ID   int    `json:"id"`
	} `json:"issue"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
}

type GetWorklogsRequest struct {
	Start  string
	End     string
}

type GetWorklogsResponse struct {
	Results []Worklog `json:"results"`
}

type CreateWorklogRequest struct {
	IssueId   string
	StartDate string
	StartTime string
	TimeSpentSeconds int
}

func (l *Worklog) TimeStarted() (time.Time, error) {
	str := l.StartDate + "T" + l.StartTime
	t, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		return time.Time{}, err
	}

	// Assumes that l.StartTime is in UTC time, so we need to convert the result to UTC
	return t.UTC(), nil
}

func (l *Worklog) TimeFinished() (time.Time, error) {
	t, err := l.TimeStarted()
	if err != nil {
		return time.Time{}, err
	}

	seconds := t.Unix()
	seconds += int64(l.TimeSpentSeconds)

	// Assumes that l.StartTime is in UTC time, so we need to convert the result to UTC
	t = time.Unix(seconds, 0).UTC()
	return t, nil
}
