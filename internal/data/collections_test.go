package data

import "testing"

func TestCollectionDefinition_IsDataValid(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"name": "string",
			"lastname": "string",
			"age": "float",
			"is_active": "bool",
		},
	}

	if !collection.IsDataValid(map[string]interface{}{
		"name": "Bob",
		"lastname": "Howards",
		"age": 20.0,
		"is_active": true,
	}) {
		t.Fatalf("unexpected invalid data")
	}
}

func TestCollectionDefinition_IsDataValid_partial(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"name": "string",
			"lastname": "string",
			"age": "float",
			"is_active": "bool",
		},
	}

	if !collection.IsDataValid(map[string]interface{}{
		"name": "Bob",
		"lastname": "Howards",
	}) {
		t.Fatalf("unexpected invalid data")
	}
}

func TestCollectionDefinition_IsDataValid_failsMismatchingFieldNames(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"name": "string",
			"lastname": "string",
			"age": "float",
			"is_active": "bool",
		},
	}

	if collection.IsDataValid(map[string]interface{}{
		"failing": "Bob",
	}) {
		t.Fatalf("unexpected success result")
	}
}

func TestCollectionDefinition_IsDataValid_failsMismatchingFieldTypeBool(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"expected": "bool",
		},
	}

	if collection.IsDataValid(map[string]interface{}{
		"expected": "not a boolean",
	}) {
		t.Fatalf("unexpected success result")
	}
}

func TestCollectionDefinition_IsDataValid_failsMismatchingFieldTypeString(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"expected": "string",
		},
	}

	if collection.IsDataValid(map[string]interface{}{
		"expected": false,
	}) {
		t.Fatalf("unexpected success result")
	}
}

func TestCollectionDefinition_IsDataValid_failsMismatchingFieldTypeFloat(t *testing.T) {
	collection := CollectionDefinition{
		Name: "test",
		Fields: map[string]string{
			"expected": "float",
		},
	}

	if collection.IsDataValid(map[string]interface{}{
		"expected": "not a float",
	}) {
		t.Fatalf("unexpected success result")
	}
}
