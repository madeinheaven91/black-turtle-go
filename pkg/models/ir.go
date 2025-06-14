package models

type Command string

const (
	Nil     Command = ""
	Lessons Command = "пары"
	Bells   Command = "звонки"
	Fio     Command = "фио"
)
