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

func FetchLessons(entityType models.StudyEntity, entityId int, queryDate time.Time) (*models.APIResponseGroup, error) {
	publicationId := os.Getenv("PUBLICATION_ID")
	var idKey string
	switch entityType {
	case models.Group:
		idKey = "groupId"
	case models.Teacher: idKey = "teacherId"
	}

	payloadText := `{"` + idKey + `":"` + strconv.Itoa(entityId) + `","date":"` + queryDate.Format("2006-01-02") + `","publicationId":"` + publicationId + `"}`
	fmt.Println(payloadText)
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

	var obj1 map[string]any
	if err := json.Unmarshal(res, &obj1); err != nil {
		panic(err)
	}
	// WRITE TO RAW FILE
	raw, _ := json.MarshalIndent(obj1, "", " ")
	f, _ := os.OpenFile("raw.json", os.O_RDWR|os.O_CREATE, 0644)
	f.Write(raw)
	if err := f.Close(); err != nil {
		panic(err)
	}

	var obj models.APIResponseGroup
	if err := json.Unmarshal(res, &obj); err != nil {
		panic(err)
	}
	// WRITE TO DESER FILE
	pretty, _ := json.MarshalIndent(obj, "", " ")
	f, _ = os.OpenFile("deser.json", os.O_RDWR|os.O_CREATE, 0644)
	f.Write(pretty)
	if err := f.Close(); err != nil {
		panic(err)
	}
	fmt.Println(string(pretty))

	return &obj, err
}

