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
	expectedData              map[string]interface{}
}

func createCollectionDefinition() data.CollectionDefinition {
	return data.CollectionDefinition{
		Name: "test",
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
		{
			Name: "test",
			Fields: map[string]string{
				"name":      "string",
				"lastname":  "string",
				"age":       "float",
				"is_active": "bool",
			},
		},
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

		var parsedItem map[string]interface{}
		json.Unmarshal(data, &parsedItem)
		if !reflect.DeepEqual(parsedItem, c.expectedData) {
			t.Logf(fmt.Sprintf("expected data %v, got %v", c.expectedData, parsedItem))
			t.Fail()
		}
	}
}
