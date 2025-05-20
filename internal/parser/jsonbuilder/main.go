package jsonbuilder

import (
	"os"
	"strconv"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/internal/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/models"
	"github.com/madeinheaven91/black-turtle-go/internal/parser/ir"
)

func validateStudyEntity(name string) (*models.DBStudyEntity, error) {
	conn := db.GetConnection()
	defer db.CloseConn(conn)

	entity, err := db.GetStudyEntity(conn, name)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func BuildPayload(query ir.Query, chatID int64) (string, error) {
	switch query.Command.(type) {
	case *ir.LessonsQuery:
		logging.Trace("Building payload for lesson query")
		input := query.Command.(*ir.LessonsQuery)
		return buildLessonsPayload(input, chatID)
	}
	logging.Trace("Unrecognized command type, no json payload")
	return "", nil
}

func buildLessonsPayload(query *ir.LessonsQuery, chatID int64) (string, error) {
	publicationId := os.Getenv("PUBLICATION_ID")

	var entity *models.DBStudyEntity
	var err error
	if query.StudyEntityName == nil {
		conn := db.GetConnection()
		defer db.CloseConn(conn)
		entity, err = db.GetStudyEntityByChat(conn, chatID)
		if err != nil {
			return "", err
		}
	} else {
		entity, err = validateStudyEntity(*query.StudyEntityName)
		if err != nil {
			return "", err
		}
	}
	date := query.Date()

	var idKey string
	switch entity.Kind {
	case models.Group:
		idKey = "groupId"
	case models.Teacher:
		idKey = "teacherId"
	}

	payloadText := `{"` + idKey + `":"` + strconv.Itoa(entity.Api_id) + `","date":"` + date.Format("2006-01-02") + `","publicationId":"` + publicationId + `"}`

	return payloadText, nil
}
