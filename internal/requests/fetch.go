package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/madeinheaven91/black-turtle-go/internal/models"
)

func FetchWeek(entityType models.StudyEntity, entityId int, queryDate time.Time) (*models.APIResponseGroup, error) {
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
		fmt.Println(err)
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
