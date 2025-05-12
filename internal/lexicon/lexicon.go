package lexicon 

import "sync"

var (
	lexiconMap = map[string]string {
		"errorMsg": "Что-то пошло не так",
	}

	mutex = &sync.RWMutex{}
)

func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, exists := lexiconMap[key]
	return val, exists
}
