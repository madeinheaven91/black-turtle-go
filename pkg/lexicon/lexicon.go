package lexicon

import (
	"fmt"
	"strings"
	"sync"
)

var (
	lexicon = map[string]string{
		"greeting": ` ____  _     ____  ____  _  __   _____  _     ____  _____  _     _____
/  _ \/ \   /  _ \/   _\/ |/ /  /__ __\/ \ /\/  __\/__ __\/ \   /  __/
| | //| |   | / \||  /  |   /     / \  | | |||  \/|  / \  | |   |  \  
| |_\\| |_/\| |-|||  \_ |   \     | |  | \_/||    /  | |  | |_/\|  /_ 
\____/\____/\_/ \|\____/\_|\_\    \_/  \____/\_/\_\  \_/  \____/\____\`,
	}

	errorLexicon = map[string]string{
		"":                    "Что-то пошло не так...",
		"general":             "Что-то пошло не так...",
		"parserError":         "Неправильно написана команда.",
		"studyEntityNotFound": "Нет такой группы/преподавателя.",
		"unknownCommand":      "Неизвестная команда.",
	}
	mutex = &sync.RWMutex{}
)

func Get(key string) string {
	mutex.RLock()
	defer mutex.RUnlock()

	val, exists := lexicon[key]
	if !exists {
		return "null"
	}
	return val
}

// Get display error message by key and provide metadata
func Error(key string, metadata ...any) string {
	mutex.RLock()
	defer mutex.RUnlock()

	val, exists := errorLexicon[key]
	if !exists {
		val = errorLexicon[""]
	}
	data := make([]string, 0)
	for _, item := range metadata {
		switch v := item.(type) {
		case string:
			data = append(data, v)
		case fmt.Stringer:
			data = append(data, v.String())
		case error:
			data = append(data, v.Error())
		default:
			data = append(data, fmt.Sprintf("%v", v))
		}
	}
	if len(data) != 0 {
		return "🚫 Ошибка!\n\n" + val + " " + strings.Join(data, ", ") + "\n\nПропишите <b>помощь</b> для вывода справки или обратитесь в техподдержку"
	} else {
		return "🚫 Ошибка!\n\n" + val + "\n\nПропишите <b>помощь</b> для вывода справки или обратитесь в техподдержку"
	}
}
