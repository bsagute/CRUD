package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"your_module_path/model" // Make sure to replace with the actual module path
)

func TestGetMetricGraphMultiline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockMetricGraphMultilineRepository(ctrl)
	mockResult := &model.MetricGraphMultiline{}
	mockRepo.EXPECT().GetMetricGraphMultiline().Return(mockResult, nil)

	service := NewMetricGraphMultilineApi(mockRepo)

	result, err := service.GetMetricGraphMultiline()
	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}

func TestGetMetricGraphMultilineJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockMetricGraphMultilineRepository(ctrl)
	mockRequest := &model.GraphRequest{}
	mockResult := &model.MetricGraphMultiline{}
	mockRepo.EXPECT().GetMetricGraphMultilineJson(mockRequest).Return(mockResult, nil)

	service := NewMetricGraphMultilineApi(mockRepo)

	result, err := service.GetMetricGraphMultilineJson(mockRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}

func TestGetMetricGraphMultilineMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockMetricGraphMultilineRepository(ctrl)
	mockRequest := &model.GraphRequest{}
	mockResult := &model.MetricGraphMultiline{}
	mockRepo.EXPECT().GetMetricGraphMultilineMetric(mockRequest).Return(mockResult, nil)

	service := NewMetricGraphMultilineApi(mockRepo)

	result, err := service.GetMetricGraphMultilineMetric(mockRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}

func TestGetMetricGraphMultilineDynamic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := model.NewMockMetricGraphMultilineRepository(ctrl)
	mockRequest := &model.GraphRequest{}
	mockResult := &model.PodMetricGraphMultiline{}
	mockRepo.EXPECT().GetMetricGraphMultilineDynamic(mockRequest).Return(mockResult, nil)

	service := NewMetricGraphMultilineApi(mockRepo)

	result, err := service.GetMetricGraphMultilineDynamic(mockRequest)
	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}