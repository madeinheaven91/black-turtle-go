package ir

import (
	"fmt"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

type QueryRaw interface {
	queryRaw()
	String() string
}

// A struct for 'пары' command.
//
// Is produced by Parser. StudyEntityName could be nil. Either no, Date, Day or Day + Week fields for TimeFrame are not nil, otherwise parser will throw an error.
type LessonsQueryRaw struct {
	StudyEntityName *string
	TimeFrame       struct {
		Date *time.Time
		Day  *string
		Week *string
	}
}

func (l LessonsQueryRaw) queryRaw() {}
func (l *LessonsQueryRaw) String() string {
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
func (l *LessonsQueryRaw) Date() *time.Time {
	if l.TimeFrame.Date != nil {
		return l.TimeFrame.Date
	} else if l.TimeFrame.Day != nil {
		var res time.Time
		var inc int
		currWeekday := shared.NormalizeWeekday(time.Now().Weekday())
		switch *l.TimeFrame.Day {
		case "сегодня":
			inc = 0
		case "завтра":
			inc = 1
		case "послезавтра":
			inc = 2
		case "вчера":
			inc = -1
		case "позавчера":
			inc = -2
		case "пн", "понедельник":
			// FIXME: these if statements should be avoidable, but i cant figure out how rn
			inc = 0 - currWeekday
			if currWeekday > 0 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "вт", "вторник":
			inc = 1 - currWeekday
			if currWeekday > 1 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "ср", "среда":
			inc = 2 - currWeekday
			if currWeekday > 2 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "чт", "четверг":
			inc = 3 - currWeekday
			if currWeekday > 3 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "пт", "пятница":
			inc = 4 - currWeekday
			if currWeekday > 4 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "сб", "суббота":
			inc = 5 - currWeekday
			if currWeekday > 5 && l.TimeFrame.Week == nil {
				inc += 7
			}
		case "вс", "воскресенье":
			inc = 6 - currWeekday
			if currWeekday > 6 && l.TimeFrame.Week == nil {
				inc += 7
			}
		}

		if l.TimeFrame.Week != nil {
			switch *l.TimeFrame.Week {
			case "след", "следующая", "следующий":
				inc += 7
			case "пред", "предыдущая", "предыдущий", "прош", "прошлый", "прошлая":
				inc -= 7
			}
		}

		res = time.Now().AddDate(0, 0, inc)
		return &res
	} else {
		res := time.Now()
		return &res
	}
}

func (l LessonsQueryRaw) Validate(chatID int64) (*LessonsQuery, error) {
	var res LessonsQuery
	res.Date = *l.Date()

	var entity *models.DBStudyEntity
	var err error
	if l.StudyEntityName == nil {
		entity, err = db.GetStudyEntityByChat(chatID)
		if err != nil {
			return nil, err
		}
	} else {
		entity, err = validateStudyEntity(*l.StudyEntityName)
		if err != nil {
			return nil, err
		}
	}

	res.StudyEntityName = entity.Name
	res.StudyEntityType = entity.Kind
	res.StudyEntityApiId = entity.ApiID

	return &res, nil
}

func validateStudyEntity(name string) (*models.DBStudyEntity, error) {
	entities, err := db.GetStudyEntities(name)
	if err != nil || len(entities) == 0 {
		return nil, err
	}
	// FIXME: maybe shouldnt be that way
	return &entities[0], nil
}

type HelpQueryRaw struct {
	Command models.Command
}

func (h HelpQueryRaw) queryRaw() {}
func (h *HelpQueryRaw) String() string {
	if h == nil {
		return ""
	}
	return fmt.Sprintf("помощь %s", h.Command)
}

func (h HelpQueryRaw) Validate() *HelpQuery {
	res := HelpQuery(h)
	return &res
}
