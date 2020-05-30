package server

import (
	log "github.com/sirupsen/logrus"
	"monkiato/apio/internal/storage"
)

var (
	//Storage main and unique storage instance used across the api
	Storage storage.Storage
)

const (
	//StorageTypeMemory identifier for storage.memory
	StorageTypeMemory  = "memory"
	//StorageTypeMongoDB identifier for storage.mongodb
	StorageTypeMongoDB = "mongodb"
)

// InitStorage is an encapsulated function for the storage initialization process
func InitStorage(apiManifest string, storageType string) {
	log.Debug("initializing storage...")
	switch storageType {
	case StorageTypeMongoDB:
		Storage = storage.NewMongoStorage()
		break
	case StorageTypeMemory:
		Storage = storage.NewMemoryStorage()
		break
	default:
		log.Fatalf("unexoected storage type initialization: " + storageType)
		break
	}
	Storage.Initialize(apiManifest)
	log.Debugf("storage ready. type: %T", Storage)
}
