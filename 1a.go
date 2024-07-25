package handlerapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/amex-eng/go-paved-road/pkg/api"
	"github.com/amex-eng/go-paved-road/pkg/model"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// Ensure MockAppInfoApi implements the api.AppInfoApi interface
var _ api.AppInfoApi = (*MockAppInfoApi)(nil)

type MockAppInfoApi struct{}

func (m *MockAppInfoApi) GetAllAppInfo(queryValues url.Values) ([]model.AppInfo, error) {
	if queryValues.Get("app_id") == "error" {
		return nil, errors.New("mock error")
	}
	return []model.AppInfo{{App_id: 1, App_name: "Test App", App_type: "Test Type", App_description: "Test Description"}}, nil
}

func TestGetAllAppInfo(t *testing.T) {
	mockService := &MockAppInfoApi{}
	handler := &AppInfoHandler{
		Apiservice: mockService,
	}

	req := httptest.NewRequest("GET", "/appinfo?app_id=test&level=level&component_id=comp&metric_id=metric", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler.GetAllAppInfo(w, req, ps)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var applist []model.AppInfo
	err := json.NewDecoder(resp.Body).Decode(&applist)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(applist))
	assert.Equal(t, int32(1), applist[0].App_id)
	assert.Equal(t, "Test App", applist[0].App_name)
	assert.Equal(t, "Test Type", applist[0].App_type)
	assert.Equal(t, "Test Description", applist[0].App_description)
}

func TestGetAllAppInfoWithError(t *testing.T) {
	mockService := &MockAppInfoApi{}
	handler := &AppInfoHandler{
		Apiservice: mockService,
	}

	req := httptest.NewRequest("GET", "/appinfo?app_id=error&level=level&component_id=comp&metric_id=metric", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler.GetAllAppInfo(w, req, ps)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var applist []model.AppInfo
	err := json.NewDecoder(resp.Body).Decode(&applist)
	assert.NoError(t, err)

	assert.Nil(t, applist)
}