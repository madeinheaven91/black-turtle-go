package msgbuilder

import (
	"fmt"
	"strings"
	"time"

	mymodels "github.com/madeinheaven91/black-turtle-go/internal/models"
)

// TODO:
func BuildDayMsg(day mymodels.SchoolDay, date time.Time, entityName string) string {
	var sb strings.Builder
	sb.WriteString("<b>")
	sb.WriteString(entityName)
	sb.WriteRune('\n')
	weekday := ""
	switch date.Weekday() - 1 {
	case 0:
		weekday = "Понедельник"
	case 1:
		weekday = "Вторник"
	case 2:
		weekday = "Вторник"
	case 3:
		weekday = "Вторник"
	case 4:
		weekday = "Вторник"
	case 5:
		weekday = "Вторник"
	case 6:
		weekday = "Вторник"
	default:
		weekday = "???"
	}

	if len(day.Lessons) == 0 {
		sb.WriteString(weekday)
		sb.WriteRune('\n')
		sb.WriteString(date.Format("02.01.06"))
		sb.WriteString("</b>\n\n")
		sb.WriteString("——————| Нет уроков |——————\n\n")
		sb.WriteString("<b>Гуляем!</b>")
	} else {
		sb.WriteString(fmt.Sprintf("%s, %d уроков\n%s</b>\n\n", weekday, len(day.Lessons), date.Format("02.01.06")))
		for _, lesson := range day.Lessons {
			sb.WriteString(fmt.Sprintf("——————| %d урок |——————\n\n", lesson.Index+1))
			sb.WriteString(fmt.Sprintf("⏳%s — %s\n", lesson.StartTime, lesson.EndTime))
			sb.WriteString(fmt.Sprintf("📖 <b>%s</b>\n", lesson.Name))
			sb.WriteString(fmt.Sprintf("🎓 %s\n", *lesson.Teacher))
			sb.WriteString(fmt.Sprintf("🔑 %s\n\n", *lesson.Cabinet))
		}
	}
	return sb.String()
}
