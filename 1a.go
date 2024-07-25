package handlerapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/amex-eng/go-paved-road/pkg/api"
	"github.com/amex-eng/go-paved-road/pkg/model"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// Ensure MockAppComponentApi implements the api.AppComponentApi interface
var _ api.AppComponentApi = (*MockAppComponentApi)(nil)

type MockAppComponentApi struct{}

func (m *MockAppComponentApi) GetAllAppComponent(queryValues url.Values) (interface{}, error) {
	if queryValues.Get("app_id") == "error" {
		return nil, errors.New("mock error")
	}
	return "mock data", nil
}

func TestGetAllAppComp(t *testing.T) {
	mockService := &MockAppComponentApi{}
	handler := &AppComponentHandler{
		AppCompservice: mockService,
	}

	req := httptest.NewRequest("GET", "/appcomponent?app_id=test&time_range=range&drill_level=level&component_id=comp", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler.GetAllAppComp(w, req, ps)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData model.AppResponse
	err := json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)

	assert.Equal(t, 200, responseData.Status)
	assert.Equal(t, "success", responseData.Message)
	assert.Equal(t, "mock data", responseData.Data)
}

func TestGetAllAppCompWithError(t *testing.T) {
	mockService := &MockAppComponentApi{}
	handler := &AppComponentHandler{
		AppCompservice: mockService,
	}

	req := httptest.NewRequest("GET", "/appcomponent?app_id=error&time_range=range&drill_level=level&component_id=comp", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler.GetAllAppComp(w, req, ps)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData model.AppResponse
	err := json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)

	assert.Equal(t, 500, responseData.Status)
	assert.Equal(t, "mock error", responseData.Message)
	assert.Nil(t, responseData.Data)
}