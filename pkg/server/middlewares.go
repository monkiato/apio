package server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io/ioutil"
	"monkiato/apio/internal/data"
	"net/http"
)

// ParseBody middleware used to apply a JSON parse for the request body. Data will be stored in Gorilla Context
// it can be obtained from subsequence handlers through context.Get(r, "parseBody")
func ParseBody(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			addErrorResponse(w, http.StatusBadRequest, "can't read body")
			return
		}
		var parsedBody map[string]interface{}
		jsonErr := json.Unmarshal(body, &parsedBody)
		if jsonErr != nil {
			addErrorResponse(w, http.StatusBadRequest, "can't parse body")
			return
		}
		context.Set(r, "parsedBody", parsedBody)
		handler.ServeHTTP(w, r)
	}
}

// ValidateID middleware used to detect an item ID in the request, if exists it means the endpoint is trying to operate
// over an existing item, and the middleware will try to find and get the item, otherwise an error is returned if the
// item was not found. The item will be stored in Gorilla Context, it can be obtained from subsequence handlers through
// context.Get(r, "item")
func ValidateID(collection data.CollectionDefinition) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			// collection validation
			collection, error := Storage.GetCollection(collection.Name)
			if error != nil {
				addErrorResponse(w, http.StatusBadRequest, "can't fetch collection. error: "+error.Error())
				return
			}

			if id, ok := vars["id"]; ok {
				// id exists, validate if the item exists in the specified collections
				item, found := collection.GetItem(id)
				if !found {
					//not found
					addErrorResponse(w, http.StatusNotFound, "item not found")
					return
				}

				context.Set(r, "id", id)
				context.Set(r, "item", item)
			}
			next.ServeHTTP(w, r)
		})
	}
}
