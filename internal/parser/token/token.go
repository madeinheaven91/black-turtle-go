package token

import "time"

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	LESSONS = "LESSONS"
	NAME    = "NAME"
	DAY     = "DAY"
	WEEK    = "WEEK"
	DATE    = "DATE"

	HELP     = "HELP"
	REGISTER = "REGISTER"
	FIO      = "FIO"
	BELLS    = "BELLS"
)

var keywords = map[string]TokenType{
	string([]byte{0}): EOF,
	"пары":            LESSONS,
	"помощь":          HELP,
	"регистрация":     REGISTER,
	"фио":             FIO,
	"звонки":          BELLS,
	"завтра":          DAY,
	"вчера":           DAY,
	"позавчера":       DAY,
	"сегодня":         DAY,
	"послезавтра":     DAY,
	"пн":              DAY,
	"понедельник":     DAY,
	"вт":              DAY,
	"вторник":         DAY,
	"ср":              DAY,
	"среда":           DAY,
	"чт":              DAY,
	"четверг":         DAY,
	"пт":              DAY,
	"пятница":         DAY,
	"сб":              DAY,
	"суббота":         DAY,
	"вс":              DAY,
	"воскресенье":     DAY,
	"эта":             WEEK,
	"этот":            WEEK,
	"след":            WEEK,
	"следующий":       WEEK,
	"следующая":       WEEK,
	"пред":            WEEK,
	"предыдущий":      WEEK,
	"предудыщуя":      WEEK,
	"прош":            WEEK,
	"прошлый":         WEEK,
	"прошлая":         WEEK,
	"неделя":          WEEK,
}

func Lookup(key string) TokenType {
	if tok, ok := keywords[key]; ok {
		return tok
	}
	if parseDate(key) {
		return DATE
	}
	return NAME
}

func parseDate(input string) bool {
	// форматы:
	// 02.01.2006
	// 02.01.06
	// 02.01

	_, err := time.Parse("02.01", input)
	if err == nil {
		return true
	}
	_, err = time.Parse("02.01.06", input)
	if err == nil {
		return true
	}
	_, err = time.Parse("02.01.2006", input)
	return err == nil
}

func New(literal string, typ TokenType) Token {
	return Token{literal, typ}
}
