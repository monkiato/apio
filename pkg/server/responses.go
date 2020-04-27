package server

import (
	"encoding/json"
	"net/http"
)

func addErrorResponse(w http.ResponseWriter, status int, error string) {
	data := map[string]interface{}{
		"success": false,
		"error": map[string]string{
			"msg": error,
		},
	}
	bytes, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(bytes)
}

func addSuccessResponse(w http.ResponseWriter, status int, extraData map[string]interface{}) {
	data := map[string]interface{}{
		"success": true,
		"data":    extraData,
	}
	bytes, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(bytes)
}
