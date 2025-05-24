package models

import (
	"sort"
	"time"
)

type APIResponse interface {
	IntoDay() SchoolWeek
}

type APIResponseGroup struct {
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
	Group APIGroup `json:"group"`
	Lessons []APILesson `json:"lessons"`
}

func (r *APIResponseGroup) IntoWeek() SchoolWeek {
	lessons := make(map[int][]Lesson, 0)
	for _, l := range r.Lessons {
		var teacher *string
		if len(l.Teachers) == 0 {
			teacher = nil
		} else {
			teacher = l.Teachers[0].Fio
		}
		lessons[l.Weekday - 1] = append(lessons[l.Weekday - 1], Lesson {
			l.Lesson - 1,
			l.StartTime,
			l.EndTime,
			l.Subject.Name,
			teacher,
			l.Cabinet.Name,
		})
	}

	days := make([]SchoolDay, 6)
	for k, v := range lessons {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Index < v[j].Index 
		})
		days[k] = SchoolDay {
			k,
			"a",
			v,
		}
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Weekday< days[j].Weekday
	})
	return SchoolWeek {
		r.Group.Name,
r.StartDate,
		r.EndDate,
		days,
	}
}

type APIGroup struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type APILesson struct {
	Cabinet APICabinet `json:"cabinet"`
	Teachers []APITeacher `json:"teachers"`
	Subject APISubject `json:"subject"`
	Lesson int `json:"lesson"`
	Weekday int `json:"weekday"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
}

type APICabinet struct {
	Name *string `json:"name"`
}

type APITeacher struct {
	Fio *string `json:"fio"`
}

type APISubject struct {
	Name string `json:"name"`
}
