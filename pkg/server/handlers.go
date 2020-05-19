package server

import (
	"encoding/json"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"monkiato/apio/internal/data"
	"net/http"
)

// GetHandler used to handle GET requests, the collectionDefinition is provided based on the endpoint being called
func GetHandler(collectionDefinition data.CollectionDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle GET for collection
		item := context.Get(r, "item")
		data, err := json.Marshal(item)
		if err != nil {
			log.Error(err.Error())
			addErrorResponse(w, http.StatusInternalServerError, "unable to parse item data")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// PutHandler used to handle PUT requests, the collectionDefinition is provided based on the endpoint being called
func PutHandler(collectionDefinition data.CollectionDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle PUT for collection
		item := context.Get(r, "parsedBody").(map[string]interface{})

		storageCollection, _ := Storage.GetCollection(collectionDefinition.Name)
		if !collectionDefinition.IsDataValid(item) {
			addErrorResponse(w, http.StatusBadRequest, "invalid item data, no matching collection definition")
			return
		}
		if id, err := storageCollection.AddItem(item); err != nil {
			log.Error(err.Error())
			addErrorResponse(w, http.StatusInternalServerError, "can't add new item")
		} else {
			addSuccessResponse(w, http.StatusCreated, map[string]interface{}{
				"id": id,
			})
		}
	}
}

// PostHandler used to handle POST requests, the collectionDefinition is provided based on the endpoint being called
func PostHandler(collectionDefinition data.CollectionDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle POST for collection
		id := context.Get(r, "id").(string)
		newItem := context.Get(r, "parsedBody").(map[string]interface{})

		storageCollection, _ := Storage.GetCollection(collectionDefinition.Name)
		if err := storageCollection.UpdateItem(id, newItem); err != nil {
			log.Error(err.Error())
			addErrorResponse(w, http.StatusInternalServerError, "can't update item")
			return
		}

		addSuccessResponse(w, http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

// DeleteHandler used to handle DELETE requests, the collectionDefinition is provided based on the endpoint being called
func DeleteHandler(collectionDefinition data.CollectionDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle DELEte for collection
		id := context.Get(r, "id").(string)
		storageCollection, _ := Storage.GetCollection(collectionDefinition.Name)
		if err := storageCollection.DeleteItem(id); err != nil {
			log.Error(err.Error())
			addErrorResponse(w, http.StatusInternalServerError, "can't delete item")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ListCollectionHandler used to get a list of items in the collection using pagination
func ListCollectionHandler(collectionDefinition data.CollectionDefinition) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		lastValue := queryParams.Get("last")
		storageCollection, _ := Storage.GetCollection(collectionDefinition.Name)
		items := storageCollection.List(lastValue)
		data, err := json.Marshal(items)
		if err != nil {
			log.Error(err.Error())
			addErrorResponse(w, http.StatusInternalServerError, "unable to parse items list data")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
