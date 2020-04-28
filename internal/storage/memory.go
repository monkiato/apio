package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"rodrigocollavo/apio/internal/data"
	"strconv"
)

type collectionData map[string]interface{}

type MemoryStorage struct {
	collectionsDefinitions []data.CollectionDefinition
	dataCollections        map[string]collectionData
	collectionHandlers     map[string]CollectionHandler
}

type MemoryCollectionHandler struct {
	collection collectionData
	lastId     int64
}

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

func (msc *MemoryCollectionHandler) GetItem(itemId string) (interface{}, bool) {
	if data, ok := msc.collection[itemId]; ok {
		return data, true
	}
	return nil, false
}

func (msc *MemoryCollectionHandler) AddItem(item map[string]interface{}) (string, error) {
	msc.lastId++
	id := strconv.FormatInt(msc.lastId, 16)
	msc.collection[id] = item
	return id, nil
}

func (msc *MemoryCollectionHandler) UpdateItem(itemId string, newItem map[string]interface{}) error {
	_, found := msc.GetItem(itemId)
	if !found {
		return fmt.Errorf("item '%s' not found", itemId)
	}
	msc.collection[itemId] = newItem
	return nil
}

func (msc *MemoryCollectionHandler) DeleteItem(itemId string) error {
	_, found := msc.GetItem(itemId)
	if !found {
		return fmt.Errorf("item '%s' not found", itemId)
	}
	delete(msc.collection, itemId)
	return nil
}

func (ms *MemoryStorage) Initialize(manifest string) {
	ms.initializeCollectionDefinitions(manifest)
	ms.initializeCollections()
}

func (ms *MemoryStorage) GetCollectionDefinitions() []data.CollectionDefinition {
	return ms.collectionsDefinitions
}

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