package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
)

func FetchWeek(payload string) (*models.APIResponseGroup, error) {
	// logging.Debug("Requesting %s with api id %d for %s\n", entityType, entityId, queryDate.Format("02.01.2006"))

	b := bytes.NewBuffer([]byte(payload))

	resp, err := http.Post(
		"https://schedule.mstimetables.ru/api/publications/group/lessons",
		"application/json",
		b,
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
