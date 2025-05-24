package ir

import (
	"fmt"
	"time"

	"github.com/madeinheaven91/black-turtle-go/pkg/models"
)

// An interface for specific commands (e.g. пары, фио, звонки).
type Query interface {
	query()
	String() string
}

type LessonsQuery struct {
	StudyEntityName  string
	StudyEntityApiId int
	StudyEntityType  models.StudyEntityType
	Date             time.Time
}

func (l LessonsQuery) query() {}
func (l LessonsQuery) String() string {
	return fmt.Sprintf("пары %s (id: %d) на %s", l.StudyEntityName, l.StudyEntityApiId, l.Date.Format("02.01.2006"))
}
