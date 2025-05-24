package json

import (
	"fmt"
	"strconv"

	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"

	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func Payload(query ir.Query) (string, error) {
	// NOTE: there are currently no other scenarios where payload is needed
	// but let it be here for scalability (i think???)
	switch query := query.(type) {
	case ir.LessonsQuery:
		return lessonsPayload(&query)
	default:
		logging.Trace("Unrecognized command type, no json payload")
		err := errors.From(fmt.Errorf("unrecognized command type %T", query), "json error", "unknownCommand", map[string]any{})
		return "", err
	}
}

func lessonsPayload(query *ir.LessonsQuery) (string, error) {
	logging.Trace("making lessons payload for %s %s (id: %d) for %s", query.StudyEntityType, query.StudyEntityName, query.StudyEntityApiId, query.Date.Format("02.01.2006"))
	publicationId := config.PublicationID()

	date := shared.GetMonday(query.Date)

	var idKey string
	switch query.StudyEntityType {
	case models.Group:
		idKey = "groupId"
	case models.Teacher:
		idKey = "teacherId"
	}

	payloadText := `{"` + idKey + `":"` + strconv.Itoa(query.StudyEntityApiId) + `","date":"` + date.Format("2006-01-02") + `","publicationId":"` + publicationId + `"}`

	return payloadText, nil
}
