package lexicon

import (
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

	errorLexicon = map[string]string {
		"": "Что-то пошло не так...",
		"general": "Что-то пошло не так...",
		"studyEntityNotFound": "Нет такой группы/преподавателя.",
		"unknownCommand": "Неизвестная команда.",
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

func Error(key string) string {
	mutex.RLock()
	defer mutex.RUnlock()

	val, exists := errorLexicon[key]
	if !exists {
		val = errorLexicon[""]
	}
	return "🚫 Ошибка!\n\n" + val + "\n\nПропишите <b>помощь</b> для вывода справки или обратитесь в техподдержку"
}
