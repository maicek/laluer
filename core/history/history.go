package history

import "os"

type EntryTypes string

var CACHE_HOME = os.Getenv("XDG_CACHE_HOME")
var HISTORY_PATH = CACHE_HOME + "/laluer/history"

const (
	ENTRY_TYPE_APP EntryTypes = "app"
)

type HistoryEntry struct {
	Type      EntryTypes
	EntryName string // Name of the entry, e.g. app name
}

func init() {

}

func PushHistory(entry HistoryEntry) {

}
