package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"io/ioutil"
	"monkiato/apio/internal/data"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	testItemId string = "abcd"
)

type TestCase struct {
	description               string
	methodType                string
	endpoint                  string
	id                        string
	item                      map[string]interface{}
	parsedBody                map[string]interface{}
	responseBodyInvalidFormat bool
	expectedStatus            int
	expectedData              interface{}
}

func createCollectionDefinition() data.CollectionDefinition {
	return data.CollectionDefinition{
		Name: "books",
		Fields: map[string]string{
			"name":      "string",
			"lastname":  "string",
			"age":       "float",
			"is_active": "bool",
		},
	}
}

func createCollectionItem() map[string]interface{} {
	return map[string]interface{}{
		"name":      "Bob",
		"lastname":  "Howards",
		"age":       20.0,
		"is_active": true,
	}
}

func createManifest(t *testing.T) string {
	definition := []data.CollectionDefinition{
		createCollectionDefinition(),
	}
	data, err := json.Marshal(definition)
	if err != nil {
		t.Fatalf("unexpected error preparing data for test")
	}
	return string(data)
}

func TestGetHandler(t *testing.T) {
	handler := GetHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}
}

func TestGetHandler_responses(t *testing.T) {
	handler := GetHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}

	cases := []TestCase{
		{
			description:               "should succeed to fetch item",
			methodType:                http.MethodGet,
			endpoint:                  "/api/books/" + testItemId,
			id:                        testItemId,
			item:                      createCollectionItem(),
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
			expectedData:              createCollectionItem(),
		},
	}

	runTestCases(t, handler, cases)
}

func TestPutHandler(t *testing.T) {
	handler := PutHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}

	InitStorage(createManifest(t), StorageTypeMemory)

	cases := []TestCase{
		{
			description:               "should succeed and create new item",
			methodType:                http.MethodPut,
			endpoint:                  "/api/books/",
			parsedBody:                createCollectionItem(),
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusCreated,
			expectedData: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "1",
				},
				"success": true,
			},
		},
		{
			description:               "should fail due to wrong item structure",
			methodType:                http.MethodPut,
			endpoint:                  "/api/books/",
			parsedBody:                map[string]interface{}{"bad": "structure"},
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "invalid item data, no matching collection definition",
				},
				"success": false,
			},
		},
		{
			description: "should fail due to wrong data type",
			methodType:  http.MethodPut,
			endpoint:    "/api/books/",
			parsedBody: map[string]interface{}{
				"name":      "Bob",
				"lastname":  "Howards",
				"age":       "invalid type",
				"is_active": true,
			},
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "invalid item data, no matching collection definition",
				},
				"success": false,
			},
		},
	}

	runTestCases(t, handler, cases)
}

func TestPostHandler(t *testing.T) {
	handler := PostHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}

	InitStorage(createManifest(t), StorageTypeMemory)

	collection, _ := Storage.GetCollection("books")
	itemId, _ := collection.AddItem(map[string]interface{}{
		"name":      "old name",
		"lastname":  "old lastname",
		"age":       5,
		"is_active": false,
	})

	cases := []TestCase{
		{
			description:               "should succeed and update item",
			methodType:                http.MethodPost,
			endpoint:                  "/api/books/" + itemId,
			id:                        itemId,
			parsedBody:                createCollectionItem(),
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
			expectedData: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "1",
				},
				"success": true,
			},
		},
		{
			description:               "should fail due to item not found in storage",
			methodType:                http.MethodPost,
			endpoint:                  "/api/books/100",
			id:                        "100",
			parsedBody:                createCollectionItem(),
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "item not found",
				},
				"success": false,
			},
		},
		{
			description:               "should fail due to wrong item structure",
			methodType:                http.MethodPost,
			endpoint:                  "/api/books/" + itemId,
			id:                        itemId,
			parsedBody:                map[string]interface{}{"bad": "structure"},
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "invalid item data, no matching collection definition",
				},
				"success": false,
			},
		},
		{
			description: "should fail due to wrong data type",
			methodType:  http.MethodPost,
			endpoint:    "/api/books/" + itemId,
			id:          itemId,
			parsedBody: map[string]interface{}{
				"name":      "Bob",
				"lastname":  "Howards",
				"age":       "invalid type",
				"is_active": true,
			},
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "invalid item data, no matching collection definition",
				},
				"success": false,
			},
		},
	}

	runTestCases(t, handler, cases)
}

func TestDeleteHandler(t *testing.T) {
	handler := DeleteHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}

	InitStorage(createManifest(t), StorageTypeMemory)

	collection, _ := Storage.GetCollection("books")
	itemId, _ := collection.AddItem(map[string]interface{}{
		"name":      "old name",
		"lastname":  "old lastname",
		"age":       5,
		"is_active": false,
	})

	cases := []TestCase{
		{
			description:               "should succeed and delete item",
			methodType:                http.MethodDelete,
			endpoint:                  "/api/books/" + itemId,
			id:                        itemId,
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusNoContent,
		},
		{
			description:               "should fail due to item not found in storage",
			methodType:                http.MethodDelete,
			endpoint:                  "/api/books/100",
			id:                        "100",
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusBadRequest,
			expectedData: map[string]interface{}{
				"error": map[string]interface{}{
					"msg": "item not found",
				},
				"success": false,
			},
		},
	}

	runTestCases(t, handler, cases)
}

func TestListCollectionHandler(t *testing.T) {
	handler := ListCollectionHandler(createCollectionDefinition())
	if handler == nil {
		t.Fatalf("unexpected null handler")
	}

	InitStorage(createManifest(t), StorageTypeMemory)

	collection, _ := Storage.GetCollection("books")
	collection.AddItem(map[string]interface{}{
		"name":      "name1",
		"lastname":  "lastname1",
		"age":       5,
		"is_active": true,
	})
	collection.AddItem(map[string]interface{}{
		"name":      "name2",
		"lastname":  "lastname2",
		"age":       10,
		"is_active": true,
	})

	cases := []TestCase{
		{
			description:               "should succeed and get collection items",
			methodType:                http.MethodGet,
			endpoint:                  "/api/books/",
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
			expectedData: []interface{}{
				map[string]interface{}{
					"name":      "name1",
					"lastname":  "lastname1",
					"age":       float64(5),
					"is_active": true,
				},
				map[string]interface{}{
					"name":      "name2",
					"lastname":  "lastname2",
					"age":       float64(10),
					"is_active": true,
				},
			},
		},
		{
			description:               "should skip 1",
			methodType:                http.MethodGet,
			endpoint:                  "/api/books/?skip=1",
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
			expectedData: []interface{}{
				map[string]interface{}{
					"name":      "name2",
					"lastname":  "lastname2",
					"age":       float64(10),
					"is_active": true,
				},
			},
		},
		{
			description:               "should skip 2",
			methodType:                http.MethodGet,
			endpoint:                  "/api/books/?skip=2",
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
		},
		{
			description:               "should limit 1",
			methodType:                http.MethodGet,
			endpoint:                  "/api/books/?limit=1",
			responseBodyInvalidFormat: false,
			expectedStatus:            http.StatusOK,
			expectedData: []interface{}{
				map[string]interface{}{
					"name":      "name1",
					"lastname":  "lastname1",
					"age":       float64(5),
					"is_active": true,
				},
			},
		},
	}

	runTestCases(t, handler, cases)
}

func runTestCases(t *testing.T, handler func(w http.ResponseWriter, r *http.Request), cases []TestCase) {
	for _, c := range cases {
		t.Logf("running test case: [%s]%s", c.methodType, c.description)
		req := httptest.NewRequest(c.methodType, c.endpoint, nil)
		w := httptest.NewRecorder()
		context.Set(req, "item", c.item)
		context.Set(req, "id", c.id)
		context.Set(req, "parsedBody", c.parsedBody)
		handler(w, req)

		if w.Code != c.expectedStatus {
			t.Logf(fmt.Sprintf("expected status %d got %d", c.expectedStatus, w.Code))
			t.Fail()
		}

		data, err := ioutil.ReadAll(w.Body)
		if (err != nil) != c.responseBodyInvalidFormat {
			t.Logf(fmt.Sprintf("unexpected error result. error expected %t, got %v", c.responseBodyInvalidFormat, err))
			t.Fail()
		}

		var parsedItem interface{}
		json.Unmarshal(data, &parsedItem)
		if !reflect.DeepEqual(parsedItem, c.expectedData) {
			t.Logf(fmt.Sprintf("expected data %v, got %v", c.expectedData, parsedItem))
			t.Fail()
		}
	}
}
