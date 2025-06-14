package messages

import (
	"fmt"
	"strings"
	"time"

	mymodels "github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

// TODO:
func BuildDayMsg(day mymodels.SchoolDay, date time.Time, entityName string) string {
	var sb strings.Builder
	sb.WriteString("<b>")
	sb.WriteString(entityName)
	sb.WriteRune('\n')
	weekday := ""
	switch shared.NormalizeWeekday(date.Weekday()) {
	case 0:
		weekday = "Понедельник"
	case 1:
		weekday = "Вторник"
	case 2:
		weekday = "Среда"
	case 3:
		weekday = "Четверг"
	case 4:
		weekday = "Пятница"
	case 5:
		weekday = "Суббота"
	case 6:
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
			if lesson.Type == "" {
				sb.WriteString(fmt.Sprintf("📖 <b>%s</b>\n", lesson.Name))
			} else {
				sb.WriteString(fmt.Sprintf("📖 <b>%s (%s)</b>\n", lesson.Name, lesson.Type))
			}
			teacher := "<i>Преподаватель не указан</i>"
			if lesson.RelatedTo != nil {
				teacher = *lesson.RelatedTo
			}
			sb.WriteString(fmt.Sprintf("🎓 %s\n", teacher))
			cabinet := "<i>Кабинет не указан</i>"
			if lesson.Cabinet != nil {
				cabinet = *lesson.Cabinet
			}
			sb.WriteString(fmt.Sprintf("🔑 %s\n\n", cabinet))
		}
	}
	return sb.String()
}
