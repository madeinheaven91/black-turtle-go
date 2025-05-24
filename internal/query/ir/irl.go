package ir

import (
	"fmt"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/query/token"
)

// An intermediate representation of user query
type Query struct {
	CommandToken token.Token
	Command      Command
}

// An interface for specific commands (e.g. пары, фио, звонки).
// Types that implement this interface are produced by Parser.
//
// String() is mainly for debugging purposes, command() is needed to associate types with interface
type Command interface {
	command()
	String() string
}

// A struct for 'пары' command.
// 
// Is produced by Parser. StudyEntityName could be nil. Either no, Date, Day or Day + Week fields for TimeFrame are not nil, otherwise parser will throw an error.
type LessonsQuery struct {
	StudyEntityName *string
	TimeFrame       struct {
		Date *time.Time
		Day  *string
		Week *string
	}
}

func (l LessonsQuery) command() {}
func (l *LessonsQuery) String() string {
	if l == nil {
		return ""
	}
	var name, date, day, week string
	if l.StudyEntityName != nil {
		name = *l.StudyEntityName
	} else {
		name = "nil"
	}
	if l.TimeFrame.Date != nil {
		date = l.TimeFrame.Date.Format("02.01.2006")
	} else {
		date = "nil"
	}
	if l.TimeFrame.Day != nil {
		day = *l.TimeFrame.Day
	} else {
		day = "nil"
	}
	if l.TimeFrame.Week != nil {
		week = *l.TimeFrame.Week
	} else {
		week = "nil"
	}
	return fmt.Sprintf("пары,%s,%s,%s,%s", name, day, week, date)
}

// Returns *time.Time based on lessons query timeframe
func (l *LessonsQuery) Date() *time.Time {
	if l.TimeFrame.Date != nil {
		return l.TimeFrame.Date
	} else if l.TimeFrame.Day != nil {
		var res time.Time
		currWeekday := int(time.Now().Weekday())
		switch *l.TimeFrame.Day {
		case "сегодня":
			res = time.Now()
		case "завтра":
			res = time.Now().AddDate(0, 0, 1)
		case "послезавтра":
			res = time.Now().AddDate(0, 0, 2)
		case "вчера":
			res = time.Now().AddDate(0, 0, -1)
		case "позавчера":
			res = time.Now().AddDate(0, 0, -2)
		case "пн", "понедельник":
			res = time.Now().AddDate(0, 0, 1-currWeekday)
		case "вт", "вторник":
			res = time.Now().AddDate(0, 0, 2-currWeekday)
		case "ср", "среда":
			res = time.Now().AddDate(0, 0, 3-currWeekday)
		case "чт", "четверг":
			res = time.Now().AddDate(0, 0, 4-currWeekday)
		case "пт", "пятница":
			res = time.Now().AddDate(0, 0, 5-currWeekday)
		case "сб", "суббота":
			res = time.Now().AddDate(0, 0, 6-currWeekday)
		case "вс", "воскресенье":
			res = time.Now().AddDate(0, 0, 7-currWeekday)
		}

		if l.TimeFrame.Week != nil {
			switch *l.TimeFrame.Week {
			case "след", "следующая", "следующий":
				res = res.AddDate(0, 0, 7)
			case "пред", "предыдущая", "предыдущий":
				res = res.AddDate(0, 0, -7)
			case "эта", "этот":
			default:
			}
		}
		return &res
	} else {
		res := time.Now()
		return &res
	}
}
