package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	botmodels "github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"

	"github.com/madeinheaven91/black-turtle-go/internal/errors"
	"github.com/madeinheaven91/black-turtle-go/internal/models"
)

func GetConnection() *pgx.Conn {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
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

	err = errors.NewBotError("study entity not found")
	return nil, err
}

func GetStudyEntityByChat(conn *pgx.Conn, chatId int64) (*models.DBStudyEntity, error) {
	row := conn.QueryRow(context.Background(), "select study_entity.* from study_entity on join chat on chat.study_entity_id = study_entity.id where chat.id=%d", chatId)
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
			return errors.NewBotError("user already exists")
		}
	}

	var name string
	if update.Message.Chat.FirstName != "" && update.Message.Chat.LastName != "" {
		name = update.Message.Chat.FirstName + " " + update.Message.Chat.LastName
	} else if update.Message.Chat.FirstName != "" {
		name = update.Message.Chat.FirstName
	} else if update.Message.Chat.Title != "" {
		name = update.Message.Chat.Title
	} else {
		name = "<unknown>"
	}
	username := update.Message.Chat.Username
	_, err := conn.Exec(context.Background(), "insert into chat(id, kind, name, username, is_banned) values ($1, $2, $3, $4, false)", update.Message.Chat.ID, update.Message.Chat.Type, name, username)
	return err
}
