package jsonbuilder

import (
	"fmt"
	"strconv"

	"github.com/madeinheaven91/black-turtle-go/internal/db"
	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"

	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
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

func Payload(query ir.Query, chatID int64) (string, error) {
	switch query.Command.(type) {
	case *ir.LessonsQuery:
		logging.Trace("Building payload for lesson query")
		input := query.Command.(*ir.LessonsQuery)
		return lessonsPayload(input, chatID)
	default:
		logging.Trace("Unrecognized command type, no json payload")
		err := errors.From(fmt.Errorf("unrecognized command type %T", query.Command), "jsonbuilder error", "unknownCommand", map[string]any{})
		return "", err
	}
}

func lessonsPayload(query *ir.LessonsQuery, chatID int64) (string, error) {
	publicationId := config.PublicationID()

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
