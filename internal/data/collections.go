package data

import (
	"fmt"
	"strings"
)

// CollectionDefinition contains the main name structure for a rest API collection
type CollectionDefinition struct {
	Name string `json:"name"`
	Fields map[string]string  `json:"fields"`
}

// IsDataValid check if the specified item map contains valid structure and field types based on the collection definition
func (cd CollectionDefinition) IsDataValid(item map[string]interface{}) bool {
	for itemKey := range item {
		if !cd.isFieldNameValid(itemKey) ||
			!cd.isFieldTypeValid(itemKey, item[itemKey]) {
			return false
		}
	}
	return true
}

func (cd CollectionDefinition) isFieldNameValid(name string) bool {
	_, exists := cd.Fields[name]
	return exists
}

func (cd CollectionDefinition) isFieldTypeValid(name string, value interface{}) bool {
	definitionType := cd.Fields[name]
	valueType := fmt.Sprintf("%T", value)
	return strings.Contains(valueType, definitionType)
}

