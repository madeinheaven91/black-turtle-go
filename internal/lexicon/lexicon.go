package lexicon

import "sync"

var (
	lexiconMap = map[string]string{
		"err": "Что-то пошло не так",
		"greeting": ` ____  _     ____  ____  _  __   _____  _     ____  _____  _     _____
/  _ \/ \   /  _ \/   _\/ |/ /  /__ __\/ \ /\/  __\/__ __\/ \   /  __/
| | //| |   | / \||  /  |   /     / \  | | |||  \/|  / \  | |   |  \  
| |_\\| |_/\| |-|||  \_ |   \     | |  | \_/||    /  | |  | |_/\|  /_ 
\____/\____/\_/ \|\____/\_|\_\    \_/  \____/\_/\_\  \_/  \____/\____\`,
	}

	mutex = &sync.RWMutex{}
)

func Get(key string) string {
	mutex.RLock()
	defer mutex.RUnlock()

	val, exists := lexiconMap[key]
	if !exists {
		return "null"
	}
	return val
}
