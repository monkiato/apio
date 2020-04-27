package server

import "rodrigocollavo/apio/internal/storage"

var (
	Storage storage.Storage
)

// InitStorage is an encapsulated function for the storage initialization process
func InitStorage(apiManifest string) {
	Storage = storage.NewMemoryStorage()
	Storage.Initialize(apiManifest)
}