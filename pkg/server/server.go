package server

import (
	log "github.com/sirupsen/logrus"
	"monkiato/apio/internal/storage"
)

var (
	Storage storage.Storage
)

// InitStorage is an encapsulated function for the storage initialization process
func InitStorage(apiManifest string) {
	log.Debug("initializing storage...")
	Storage = storage.NewMongoStorage()
	Storage.Initialize(apiManifest)
	log.Debugf("storage ready. type: %T", Storage)
}