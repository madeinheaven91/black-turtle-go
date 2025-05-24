package models

import "time"

type StudyEntityType string 
const (
	Teacher StudyEntityType = "teacher"
	Group StudyEntityType = "group"
)

type SchoolWeek struct {
	EntityName string 
	StartDate time.Time 
	EndDate time.Time 
	Days []SchoolDay
}

type SchoolDay struct {
	Weekday int
	Date string 
	Lessons []Lesson
}

type Lesson struct {
	Index int 
	StartTime string 
	EndTime string 
	Name string 
	Teacher *string 
	Cabinet *string
}
