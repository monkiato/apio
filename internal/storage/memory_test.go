package storage

import "testing"

func createCollection() collectionData {
	return map[string]interface{}{
		"1": createItem(),
	}
}

func createItem() map[string]interface{} {
	return map[string]interface{} {
		"name": "Bob",
		"lastname": "Howards",
		"age": 20.0,
		"is_active": true,
	}
}

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()
	if storage == nil {
		t.Fatalf("unexpected nil instance")
	}
}

func TestMemoryCollectionHandler_AddItem(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: map[string]interface{}{},
	}
	id, err := handler.AddItem(createItem())
	if err != nil {
		t.Errorf("unexpected error: " + err.Error())
	}
	if id != "1" {
		t.Errorf("unexpected id: " + id)
	}
}

func TestMemoryCollectionHandler_DeleteItem(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	if err := handler.DeleteItem("1"); err != nil {
		t.Fatalf("unexpected error: " + err.Error())
	}
}

func TestMemoryCollectionHandler_DeleteItem_notFound(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	if err := handler.DeleteItem("2"); err == nil {
		t.Fatalf("unexpected success result")
	}
}

func TestMemoryCollectionHandler_GetItem(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	data, found := handler.GetItem("1")
	if !found {
		t.Errorf("item not found")
	}
	if data == nil {
		t.Fatalf("unexpected nil data")
	}
	if len(data.(map[string]interface{})) != 4 {
		t.Fatalf("unexpected properties in data")
	}
}

func TestMemoryCollectionHandler_GetItem_notFound(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	data, found := handler.GetItem("2")
	if found {
		t.Errorf("unexpected item found")
	}
	if data != nil {
		t.Fatalf("unexpected valid data")
	}
}

func TestMemoryCollectionHandler_Query(t *testing.T) {
	// not implemented here
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	list, err := handler.Query(QueryParams{})
	if list != nil || err != nil {
		t.Fatalf("unexpected values for not implemented method")
	}
}

func TestMemoryCollectionHandler_UpdateItem(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	if err := handler.UpdateItem("1", map[string]interface{} {
		"name": "Bob updated",
		"lastname": "Howards updated",
		"age": 10.0,
		"is_active": true,
	}); err != nil {
		t.Fatalf("unexpected error: " + err.Error())
	}
}

func TestMemoryCollectionHandler_UpdateItem_wrongId(t *testing.T) {
	handler := &MemoryCollectionHandler{
		collection: createCollection(),
	}
	if err := handler.UpdateItem("2", map[string]interface{} {
		"name": "Bob updated",
		"lastname": "Howards updated",
		"age": 10.0,
		"is_active": true,
	}); err == nil {
		t.Fatalf("expected error for unexisting id 2")
	}
}