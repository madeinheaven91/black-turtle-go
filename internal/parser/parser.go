package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/models"
)

func parseDate(input string) (*time.Time, error) {
	// форматы:
	// 02.01.2006
	// 02.01.06
	// 02.01
	// пн
	// завтра

	currWeekday := int(time.Now().Weekday())
	switch input {
	case "сегодня":
		res := time.Now()
		return &res, nil
	case "завтра":
		res := time.Now().AddDate(0, 0, 1)
		return &res, nil
	case "послезавтра":
		res := time.Now().AddDate(0, 0, 2)
		return &res, nil
	case "вчера":
		res := time.Now().AddDate(0, 0, -1)
		return &res, nil
	case "позавчера":
		res := time.Now().AddDate(0, 0, -2)
		return &res, nil
	case "пн", "понедельник":
		res := time.Now().AddDate(0, 0, 1 - currWeekday)
		return &res, nil
	case "вт", "вторник":
		res := time.Now().AddDate(0, 0, 2 - currWeekday)
		return &res, nil
	case "ср", "среда":
		res := time.Now().AddDate(0, 0, 3 - currWeekday)
		return &res, nil
	case "чт", "четверг":
		res := time.Now().AddDate(0, 0, 4 - currWeekday)
		return &res, nil
	case "пт", "пятница":
		res := time.Now().AddDate(0, 0, 5 - currWeekday)
		return &res, nil
	case "сб", "суббота":
		res := time.Now().AddDate(0, 0, 6 - currWeekday)
		return &res, nil
	case "вс", "воскресенье":
		res := time.Now().AddDate(0, 0, 7 - currWeekday)
		return &res, nil
	}

	{
		res, err := time.Parse("02.01", input)
		if err == nil {
			res = res.AddDate(time.Now().Year(), 0, 0)
			return &res, nil
		}
	}
	{
		res, err := time.Parse("02.01.06", input)
		if err == nil {
			return &res, nil
		}
	}
	{
		res, err := time.Parse("02.01.2006", input)
		if err == nil {
			return &res, nil
		}
	}

	e := fmt.Errorf("couldn't parse date")
	return nil, e
}

func ParseToRequest(input string) (*models.LessonRequest, error) {
	// ...
	// пн
	// пн след
	// 314
	// 314 пн
	// 314 пн след
	// димитриев
	// димитриев пн
	// димитриев пн след
	// александр оелгович
	// александр оелгович пн
	// александр оелгович пн след
	// <имя?> <дата?> | <день?> <неделя?>

	tokens := strings.Split(input, " ")
	tokens = tokens[1:]
	fmt.Println(tokens)
	if len(tokens) == 0 {
		return &models.LessonRequest{
			Date:            time.Now(),
			StudyEntityType: models.Group, // fetch from db
			StudyEntityId:   2,            // fetch from db
		}, nil
	}

	nameToken := ""
	var dayTokenIndex, weekTokenIndex *int = nil, nil
	var date *time.Time = nil

	for i, token := range tokens {
		if d, err := parseDate(token); err == nil {
			date = d
			dayTokenIndex = &i
		}
	}

	if dayTokenIndex != nil {
		tmp := strings.Join(tokens[:*dayTokenIndex], " ")
		nameToken = tmp
		if len(tokens) > *dayTokenIndex + 1 {
			weekTmp := (*dayTokenIndex + 1)
			weekTokenIndex = &weekTmp
		}
	}else {
		tmp := strings.Join(tokens, " ")
		nameToken = tmp
	}


	namePrint, dayPrint, weekPrint := "", "", ""
	if nameToken != "" {
		namePrint = nameToken
	}else {
		namePrint = "null"
	}

	if dayTokenIndex != nil {
		dayPrint = tokens[*dayTokenIndex]
	}else {
		tmp := time.Now()
		date = &tmp
		dayPrint = "null"
	}

	fmt.Println(date)
	if weekTokenIndex != nil {
		weekPrint = tokens[*weekTokenIndex]
		date, _ = addWeek(date, weekPrint)
		fmt.Println(date)
	}else {
		weekPrint = "null"
	}

	fmt.Println(namePrint, dayPrint, weekPrint)


	// TODO: fetch from db to validate and stuff

	e := fmt.Errorf("todo")
	return nil, e
}

func addWeek(date *time.Time, modifier string) (*time.Time, error) {
	switch modifier {
	case "след", "следующая", "следующий":
		tmp := date.AddDate(0, 0, 7)
		date = &tmp
	case "пред", "предыдущая", "предыдущий":
		tmp := date.AddDate(0, 0, -7)
		date = &tmp
	case "эта", "этот":
		//
	default:
		e := fmt.Errorf("invalid week modifier")
		return nil, e
	}
	return date, nil
}
