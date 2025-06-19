package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/madeinheaven91/black-turtle-go/internal/query/ir"
	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func FetchWeek(query *ir.LessonsQuery) (*models.APIResponse, error) {
	payload := lessonsPayload(query)
	b := bytes.NewBuffer([]byte(payload))

	resp, err := http.Post(
		fmt.Sprintf("https://schedule.mstimetables.ru/api/publications/%s/lessons", query.StudyEntityType),
		"application/json",
		b,
	)
	if err != nil {
		logging.Error("%s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	var obj models.APIResponse
	if err := json.Unmarshal(res, &obj); err != nil {
		panic(err)
	}

	return &obj, err
}

func lessonsPayload(query *ir.LessonsQuery) string {
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

	return payloadText
}
