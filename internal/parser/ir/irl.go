package ir

import (
	"fmt"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/parser/token"
)

type Query struct {
	CommandToken token.Token
	Command      Command
}

type Command interface {
	command()
}

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
	return fmt.Sprintf("пары %s %s %s %s", name, day, week, date)
}
