package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	botmodels "github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"

	"github.com/madeinheaven91/black-turtle-go/pkg/config"
	"github.com/madeinheaven91/black-turtle-go/pkg/errors"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func GetConnection() *pgx.Conn {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.PgUser(),
		config.PgPassword(),
		config.PgHost(),
		config.PgPort(),
		config.PgName(),
	)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		logging.Critical("%s\n", err)
		os.Exit(1)
	}
	return conn
}

// NOTE: so that importing context isnt necessary
func CloseConn(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func GetStudyEntity(conn *pgx.Conn, input string) (*models.DBStudyEntity, error) {
	rows, err := conn.Query(context.Background(), "select * from study_entity")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var api_id int
		var kind string
		var name string
		err := rows.Scan(&id, &api_id, &kind, &name)
		if err != nil {
			return nil, err
		}
		if strings.Contains(strings.ToLower(name), strings.ToLower(input)) {
			res := models.DBStudyEntity{
				Id:     id,
				Api_id: api_id,
				Kind:   models.StudyEntityType(kind),
				Name:   name,
			}
			return &res, nil
		}
	}

	err = errors.From(fmt.Errorf("study entity not found"), "db error", "studyEntityNotFound", map[string]any{
		"studyEntityName": input,
	})
	return nil, err
}

func GetStudyEntityByChat(conn *pgx.Conn, chatId int64) (*models.DBStudyEntity, error) {
	row := conn.QueryRow(context.Background(), "select study_entity.* from study_entity join chat on chat.study_entity_id=study_entity.id where chat.id=$1", chatId)
	var id int
	var api_id int
	var kind string
	var name string
	err := row.Scan(&id, &api_id, &kind, &name)
	if err != nil {
		return nil, err
	}
	res := models.DBStudyEntity{
		Id:     id,
		Api_id: api_id,
		Kind:   models.StudyEntityType(kind),
		Name:   name,
	}
	return &res, nil
}

func AddChat(conn *pgx.Conn, update *botmodels.Update) error {
	{
		var existing string
		err := conn.QueryRow(context.Background(), "select name from chat where id=$1", update.Message.Chat.ID).Scan(&existing)
		if err == nil {
			return errors.From(fmt.Errorf("chat already exists"), "db error", "", map[string]any{
				"id": update.Message.Chat.ID,
			})
		}
	}

	name := shared.GetChatName(update)
	username := update.Message.Chat.Username
	_, err := conn.Exec(context.Background(), "insert into chat(id, kind, name, username, is_banned) values ($1, $2, $3, $4, false)", update.Message.Chat.ID, update.Message.Chat.Type, name, username)
	return err
}
