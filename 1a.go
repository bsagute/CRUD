package main

import (
	"strconv"
	"testing"
)

// Assuming the model package and AppComponentLayout type are defined as below
package model

type AppComponentLayout struct {
	SortValue int
}

func sortIndex(containers []model.AppComponentLayout, sort string) int {
	sortValue, err := strconv.Atoi(sort)
	if err != nil {
		return len(containers) // in case sort is not present
	}

	for i, graphType := range containers {
		existingSortValue := graphType.SortValue
		if int64(sortValue) < existingSortValue {
			return i
		}
	}
	return len(containers)
}

func TestSortIndex(t *testing.T) {
	tests := []struct {
		containers []model.AppComponentLayout
		sort       string
		expected   int
	}{
		{
			containers: []model.AppComponentLayout{{SortValue: 2}, {SortValue: 4}, {SortValue: 6}},
			sort:       "5",
			expected:   2,
		},
		{
			containers: []model.AppComponentLayout{{SortValue: 2}, {SortValue: 4}, {SortValue: 6}},
			sort:       "1",
			expected:   0,
		},
		{
			containers: []model.AppComponentLayout{{SortValue: 2}, {SortValue: 4}, {SortValue: 6}},
			sort:       "7",
			expected:   3,
		},
		{
			containers: []model.AppComponentLayout{{SortValue: 2}, {SortValue: 4}, {SortValue: 6}},
			sort:       "abc",
			expected:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.sort, func(t *testing.T) {
			result := sortIndex(tt.containers, tt.sort)
			if result != tt.expected {
				t.Errorf("sortIndex(%v, %s) = %d; expected %d", tt.containers, tt.sort, result, tt.expected)
			}
		})
	}
}
