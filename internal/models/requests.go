package models

import "time"

type LessonRequest struct {
	Date time.Time
	StudyEntityType StudyEntityType
	StudyEntityId int
}

