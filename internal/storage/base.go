package storage

import "monkiato/apio/internal/data"

// Storage handles data for multiple collections, it's the main entry points to initialize and manage all API collections
type Storage interface {
	// Initialize must be called before any other method to initialize main collection structures based on the specified json formatted manifest
	Initialize(manifest string)
	// GetCollectionDefinitions get a map with all collection definitions
	GetCollectionDefinitions() []data.CollectionDefinition
	// GetCollection get a single collection handler for the specified collection name
	GetCollection(collectionName string) (CollectionHandler, error)
}

// CollectionHandler used to operate over a single collection
type CollectionHandler interface {
	// GetItem get a collection item for the specified item ID
	GetItem(itemID string) (interface{}, bool)
	// AddItem insert new item. (itemID, error) is returned
	AddItem(item map[string]interface{}) (string, error)
	// UpdateItem used to update an existing item, it must exists previously, otherwise an error will be returned
	UpdateItem(itemID string, item map[string]interface{}) error
	// DeleteItem remove the specified itemID
	DeleteItem(itemID string) error
	// Query returns a list of items from a collection filtered by some criteria declared in QueryParams
	Query(query QueryParams) ([]interface{}, error)
}

// QueryParams used to filter data on a query
type QueryParams struct {
	Skip   int64
	Limit  int64
	SortBy string
}
