package models

const (
	CacheCommandRemove CacheCommand = "REMOVE"
	CacheCommandPurge  CacheCommand = "Purge"
)

type CacheCommand string

type CacheMsg struct {
	Command CacheCommand `json:"command"`
	Key     interface{}  `json:"key"`
}
