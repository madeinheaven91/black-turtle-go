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
	"github.com/madeinheaven91/black-turtle-go/pkg/lexicon"
	"github.com/madeinheaven91/black-turtle-go/pkg/logging"
	"github.com/madeinheaven91/black-turtle-go/pkg/models"
	"github.com/madeinheaven91/black-turtle-go/pkg/shared"
)

func connection() *pgx.Conn {
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
func close(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func StudyEntities(input string) ([]models.DBStudyEntity, error) {
	conn := connection()
	defer close(conn)
	rows, err := conn.Query(context.Background(), "select * from study_entity")
	if err != nil {
		return nil, err
	}
	res := make([]models.DBStudyEntity, 0, 1)
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
			new := models.DBStudyEntity{
				ID:    id,
				ApiID: api_id,
				Kind:  models.StudyEntityType(kind),
				Name:  name,
			}
			res = append(res, new)
		}
	}

	return res, nil
}

func StudyEntityByChat(chatId int64) (*models.DBStudyEntity, error) {
	conn := connection()
	defer close(conn)
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
		ID:    id,
		ApiID: api_id,
		Kind:  models.StudyEntityType(kind),
		Name:  name,
	}
	return &res, nil
}

func StudyEntityByID(id int) (*models.DBStudyEntity, error) {
	conn := connection()
	defer close(conn)
	row := conn.QueryRow(context.Background(), "select api_id, kind, name from study_entity where id = $1;", id)
	var api_id int
	var kind string
	var name string
	err := row.Scan(&api_id, &kind, &name)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	res := models.DBStudyEntity{
		ID:    id,
		ApiID: api_id,
		Kind:  models.StudyEntityType(kind),
		Name:  name,
	}
	return &res, nil
}

func AddChat(update *botmodels.Update) error {
	conn := connection()
	defer close(conn)
	{
		var existing string
		err := conn.QueryRow(context.Background(), "select name from chat where id=$1", update.Message.Chat.ID).Scan(&existing)
		if err == nil {
			return errors.From(fmt.Errorf("chat already exists"), "db error", lexicon.EGeneral, map[string]any{
				"id": update.Message.Chat.ID,
			})
		}
	}

	name := shared.GetChatName(update)
	username := update.Message.Chat.Username
	_, err := conn.Exec(context.Background(), "insert into chat(id, kind, name, username, is_banned) values ($1, $2, $3, $4, false)", update.Message.Chat.ID, update.Message.Chat.Type, name, username)
	return err
}

func Chat(chatId int64) *models.DBChat {
	conn := connection()
	defer close(conn)
	row := conn.QueryRow(context.Background(), "select * from chat where id=$1", chatId)
	var id int64
	var kind string
	var name string
	var username string
	var study_entity_id int
	var is_banned bool
	err := row.Scan(&id, &kind, &name, &username, &study_entity_id, &is_banned)
	if err != nil {
		return nil
	}
	res := models.DBChat{
		ID:            id,
		Kind:          kind,
		Name:          name,
		Username:      &username,
		StudyEntityID: &study_entity_id,
		IsBanned:      is_banned,
	}
	return &res
}

func AssignStudyEntity(update *botmodels.Update, studyEntity *models.DBStudyEntity) error {
	conn := connection()
	defer close(conn)
	chatID := shared.GetChatID(update)
	_, err := conn.Exec(context.Background(), "update chat set study_entity_id=$1 where id=$2", studyEntity.ID, chatID)
	return err
}

func CheckAdmin(chatID int64) bool {
	conn := connection()
	defer close(conn)
	var id int
	err := conn.QueryRow(context.Background(), "select id from admin where chat_id=$1", chatID).Scan(&id)
	return err == nil
}

func Chats() ([]models.DBChat, error) {
	conn := connection()
	defer close(conn)
	rows, err := conn.Query(context.Background(), "select * from chat")
	if err != nil {
		return nil, err
	}
	res := make([]models.DBChat, 0, 1)
	for rows.Next() {
		var id int64
		var kind string
		var name string
		var username *string
		var studyEntityID *int
		var isBanned bool
		err := rows.Scan(&id, &kind, &name, &username, &studyEntityID, &isBanned)
		if err == nil {
			res = append(res, models.DBChat{
				ID:            id,
				Kind:          kind,
				Name:          name,
				Username:      username,
				StudyEntityID: studyEntityID,
				IsBanned:      isBanned,
			})
		}
	}

	return res, nil
}

func ChatCount() (int, int, error) {
	conn := connection()
	defer close(conn)
	private, groups := 0, 0
	row := conn.QueryRow(context.Background(), `select 
		count(case when kind='private' then 1 end),
		count(case when kind!='private' then 1 end) 
	from chat`)
	err := row.Scan(&private, &groups)
	return private, groups, err
}

func SyncName(update *botmodels.Update) error {
	type pair struct {
		fst bool
		snd bool
	}
	conn := connection()
	defer close(conn)
	row := conn.QueryRow(context.Background(), "select * from chat where id=$1", update.Message.Chat.ID)
	var name string
	var username string
	err := row.Scan(nil, nil, &name, &username, nil, nil); if err != nil {
		return err
	}

	newName := shared.GetChatName(update)
	newUsername := update.Message.Chat.Username
	matches := pair{
		name == newName,
		username == newUsername,
	}
	switch matches {
	case pair{false, true}:
		_, err = conn.Exec(context.Background(), "update chat set name=$1 where id=$2", newName, update.Message.Chat.ID)
	case pair{true, false}:
		_, err = conn.Exec(context.Background(), "update chat set username=$1 where id=$2", newUsername, update.Message.Chat.ID)
	case pair{false, false}:
		_, err = conn.Exec(context.Background(), "update chat set name=$1, username=$2 where id=$3", newName, newUsername, update.Message.Chat.ID)
	}
	return err
}
