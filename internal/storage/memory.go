package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"monkiato/apio/internal/data"
	"strconv"
)

type collectionData map[string]interface{}

//MemoryStorage structure for the storage using in-memory data (ideal for testing, not for production)
type MemoryStorage struct {
	collectionsDefinitions []data.CollectionDefinition
	dataCollections        map[string]collectionData
	collectionHandlers     map[string]CollectionHandler
}

//MemoryCollectionHandler data handler used for a specific collection
type MemoryCollectionHandler struct {
	collection collectionData
	lastID     int64
}

//NewMemoryStorage create a new MemoryStarage instance
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		dataCollections:    map[string]collectionData{},
		collectionHandlers: map[string]CollectionHandler{},
	}
}

func newMemoryStorageCollectionHandler(collection collectionData) CollectionHandler {
	return &MemoryCollectionHandler{
		collection: collection,
	}
}

//GetItem implements storage.CollectionHandler.GetItem
func (msc *MemoryCollectionHandler) GetItem(itemID string) (interface{}, bool) {
	if data, ok := msc.collection[itemID]; ok {
		return data, true
	}
	return nil, false
}

//AddItem implements storage.CollectionHandler.AddItem
func (msc *MemoryCollectionHandler) AddItem(item map[string]interface{}) (string, error) {
	msc.lastID++
	id := strconv.FormatInt(msc.lastID, 16)
	msc.collection[id] = item
	return id, nil
}

//UpdateItem implements storage.CollectionHandler.UpdateItem
func (msc *MemoryCollectionHandler) UpdateItem(itemID string, newItem map[string]interface{}) error {
	_, found := msc.GetItem(itemID)
	if !found {
		return fmt.Errorf("item '%s' not found", itemID)
	}
	msc.collection[itemID] = newItem
	return nil
}

//DeleteItem implements storage.CollectionHandler.DeleteItem
func (msc *MemoryCollectionHandler) DeleteItem(itemID string) error {
	_, found := msc.GetItem(itemID)
	if !found {
		return fmt.Errorf("item '%s' not found", itemID)
	}
	delete(msc.collection, itemID)
	return nil
}

//Query implements storage.CollectionHandler.Query
func (msc *MemoryCollectionHandler) Query(query QueryParams) ([]interface{}, error) {
	return nil, nil
}

//Initialize implements storage.Storage.Initialize
func (ms *MemoryStorage) Initialize(manifest string) {
	ms.initializeCollectionDefinitions(manifest)
	ms.initializeCollections()
}

//GetCollectionDefinitions implements storage.Storage.GetCollectionDefinitions
func (ms *MemoryStorage) GetCollectionDefinitions() []data.CollectionDefinition {
	return ms.collectionsDefinitions
}

//GetCollection implements storage.Storage.GetCollection
func (ms *MemoryStorage) GetCollection(collectionName string) (CollectionHandler, error) {
	if collection, ok := ms.dataCollections[collectionName]; ok {
		storageCollection, exists := ms.collectionHandlers[collectionName]
		if !exists {
			storageCollection = newMemoryStorageCollectionHandler(collection)
			ms.collectionHandlers[collectionName] = storageCollection
		}
		return storageCollection, nil
	}
	return nil, fmt.Errorf("collection %s not found", collectionName)
}

func (ms *MemoryStorage) initializeCollectionDefinitions(manifest string) {
	err := json.Unmarshal([]byte(manifest), &ms.collectionsDefinitions)
	if err != nil {
		log.Fatal("Unable to parse manifest")
	}
}

func (ms *MemoryStorage) initializeCollections() {
	for _, collectionDefinition := range ms.collectionsDefinitions {
		ms.dataCollections[collectionDefinition.Name] = collectionData{}
	}
}
