package messages

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
	switch date.Weekday() {
	case 1:
		weekday = "Понедельник"
	case 2:
		weekday = "Вторник"
	case 3:
		weekday = "Среда"
	case 4:
		weekday = "Четверг"
	case 5:
		weekday = "Пятница"
	case 6:
		weekday = "Суббота"
	case 7:
		weekday = "Воскресенье"
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
			var teacher = "<i>Преподаватель не указан</i>"
			if lesson.Cabinet != nil {
				teacher = *lesson.Teacher
			}
			sb.WriteString(fmt.Sprintf("🎓 %s\n", teacher))
			var cabinet = "<i>Кабинет не указан</i>"
			if lesson.Cabinet != nil {
				cabinet = *lesson.Cabinet
			}
			sb.WriteString(fmt.Sprintf("🔑 %s\n\n", cabinet))
		}
	}
	return sb.String()
}
