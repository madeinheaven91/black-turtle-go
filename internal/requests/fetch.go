package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/logging"
	"github.com/madeinheaven91/black-turtle-go/internal/models"
)

func FetchWeek(entityType models.StudyEntityType, entityId int, queryDate time.Time) (*models.APIResponseGroup, error) {
	logging.Debug("Requesting %s with api id %d for %s\n", entityType, entityId, queryDate.Format("02.01.2006"))
	publicationId := os.Getenv("PUBLICATION_ID")

	var idKey string
	switch entityType {
	case models.Group:
		idKey = "groupId"
	case models.Teacher:
		idKey = "teacherId"
	}

	payloadText := `{"` + idKey + `":"` + strconv.Itoa(entityId) + `","date":"` + queryDate.Format("2006-01-02") + `","publicationId":"` + publicationId + `"}`
	payload := bytes.NewBuffer([]byte(payloadText))

	resp, err := http.Post(
		"https://schedule.mstimetables.ru/api/publications/group/lessons",
		"application/json",
		payload,
	)
	if err != nil {
		logging.Error("%s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	res, _ := io.ReadAll(resp.Body)

	var obj models.APIResponseGroup
	if err := json.Unmarshal(res, &obj); err != nil {
		panic(err)
	}

	return &obj, err
}
