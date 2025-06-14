package models

import "time"

type StudyEntityType string

const (
	Teacher StudyEntityType = "teacher"
	Group   StudyEntityType = "group"
)

type SchoolWeek struct {
	Owner     string
	StartDate time.Time
	EndDate   time.Time
	Days      []SchoolDay
}

type SchoolDay struct {
	Weekday int
	Lessons []Lesson
}

type Lesson struct {
	Index     int
	Type      string
	StartTime string
	EndTime   string
	Name      string
	RelatedTo *string
	Cabinet   *string
}
