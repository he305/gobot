package database

type Database interface {
	GetEntryByName(collectionName string, key string, name string) (map[string]interface{}, error)
	AddEntry(collectionName string, entry map[string]interface{}) error
	AddEntries(collectionName string, entries []map[string]interface{}) error
}
