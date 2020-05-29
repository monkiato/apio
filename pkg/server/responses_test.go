package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_addErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	addErrorResponse(recorder, http.StatusInternalServerError, "testing error")
	if recorder.Code != 500 {
		t.Error("unexpected status code")
	}
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("unexpected error reading body: " + err.Error())
	}
	if string(data) != "{\"error\":{\"msg\":\"testing error\"},\"success\":false}" {
		t.Error("unexpected body: " + string(data))
	}
}

func Test_addSuccessResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	addSuccessResponse(recorder, http.StatusOK, nil)
	if recorder.Code != 200 {
		t.Error("unexpected status code")
	}
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("unexpected error reading body: " + err.Error())
	}
	if string(data) != "{\"data\":null,\"success\":true}" {
		t.Error("unexpected body: " + string(data))
	}
}

func Test_addSuccessResponse_withExtraData(t *testing.T) {
	recorder := httptest.NewRecorder()
	addSuccessResponse(recorder, http.StatusOK, map[string]interface{}{
		"more": "data",
	})
	if recorder.Code != 200 {
		t.Error("unexpected status code")
	}
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("unexpected error reading body: " + err.Error())
	}
	if string(data) != "{\"data\":{\"more\":\"data\"},\"success\":true}" {
		t.Error("unexpected body: " + string(data))
	}
}
