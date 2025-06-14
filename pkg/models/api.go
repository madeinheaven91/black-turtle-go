package models

import (
	"sort"
	"time"
)

type APIResponse struct {
	StartDate time.Time   `json:"startDate"`
	EndDate   time.Time   `json:"endDate"`
	Group     *APIGroup   `json:"group"`
	Teacher   *APITeacher `json:"teacher"`
	Lessons   []APILesson `json:"lessons"`
}

// FIXME: вроде работает, но я не понимаю что тут происходит
func (r *APIResponse) IntoWeek() SchoolWeek {
	if r.Group != nil {
		// deser group
		lessons := make(map[int][]Lesson)
		for _, l := range r.Lessons {
			var teacher *string
			var lessonType string
			if len(l.Teachers) == 0 {
				teacher = nil
			} else {
				teacher = l.Teachers[0].Fio
			}
			if l.Type == nil {
				lessonType = ""
			} else {
				lessonType = l.Type.Name
			}
			lessons[l.Weekday-1] = append(lessons[l.Weekday-1], Lesson{
				l.Lesson - 1,
				lessonType,
				l.StartTime,
				l.EndTime,
				l.Subject.Name,
				teacher,
				l.Cabinet.Name,
			})
		}

		days := make([]SchoolDay, 7)
		for k, v := range lessons {
			sort.Slice(v, func(i, j int) bool {
				return v[i].Index < v[j].Index
			})
			days[k] = SchoolDay{
				k,
				v,
			}
		}
		sort.Slice(days, func(i, j int) bool {
			return days[i].Weekday < days[j].Weekday
		})
		for len(days) < 7 {
			days = append(days, SchoolDay{})
		}
		res := SchoolWeek{
			r.Group.Name,
			r.StartDate,
			r.EndDate,
			days,
		}
		return res
	} else if r.Teacher != nil {
		// deser teacher
		lessons := make(map[int][]Lesson)
		for _, l := range r.Lessons {
			var group *string
			var lessonType string
			if len(l.UnionGroups) == 0 {
				group = nil
			} else {
				group = &l.UnionGroups[0].Group.Name
			}
			if l.Type == nil {
				lessonType = ""
			} else {
				lessonType = l.Type.Name
			}
			lessons[l.Weekday-1] = append(lessons[l.Weekday-1], Lesson{
				l.Lesson - 1,
				lessonType,
				l.StartTime,
				l.EndTime,
				l.Subject.Name,
				group,
				l.Cabinet.Name,
			})
		}

		days := make([]SchoolDay, 7)
		for k, v := range lessons {
			sort.Slice(v, func(i, j int) bool {
				return v[i].Index < v[j].Index
			})
			days[k] = SchoolDay{
				k,
				v,
			}
		}
		sort.Slice(days, func(i, j int) bool {
			return days[i].Weekday < days[j].Weekday
		})
		for len(days) < 7 {
			days = append(days, SchoolDay{})
		}
		res := SchoolWeek{
			*r.Teacher.Fio,
			r.StartDate,
			r.EndDate,
			days,
		}
		return res
	} else {
		panic("watafuk")
	}
}

type APICabinet struct {
	Name *string `json:"name"`
}

type APITeacher struct {
	Fio *string `json:"fio"`
}

type APIUnionGroup struct {
	Group APIGroup `json:"group"`
}

type APIGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APISubject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APITypeLesson struct {
	Name string `json:"name"`
}

type APILesson struct {
	Cabinet     APICabinet      `json:"cabinet"`
	Teachers    []APITeacher    `json:"teachers"`
	UnionGroups []APIUnionGroup `json:"unionGroups"`
	Subject     APISubject      `json:"subject"`
	Type        *APITypeLesson  `json:"typeLesson"`
	Lesson      int             `json:"lesson"`
	Weekday     int             `json:"weekday"`
	StartTime   string          `json:"startTime"`
	EndTime     string          `json:"endTime"`
}
