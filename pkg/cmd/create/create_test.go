// ABOUTME: Tests for create command functionality
// ABOUTME: Validates service list computation logic
package create

import (
	"reflect"
	"testing"
)

func TestComputeServiceList(t *testing.T) {
	tests := []struct {
		name        string
		services    []string
		noServices  []string
		addServices []string
		expected    []string
	}{
		{
			name:     "explicit services override",
			services: []string{"apiserver", "bot"},
			expected: []string{"apiserver", "bot"},
		},
		{
			name:     "services none",
			services: []string{"none"},
			expected: []string{},
		},
		{
			name:     "no flags uses defaults",
			expected: []string{"apiserver"},
		},
		{
			name:       "exclude service",
			noServices: []string{"apiserver"},
			expected:   []string{},
		},
		{
			name:        "add services",
			addServices: []string{"bot", "scheduler"},
			expected:    []string{"apiserver", "bot", "scheduler"},
		},
		{
			name:        "combined exclude and add",
			noServices:  []string{"apiserver"},
			addServices: []string{"admserver"},
			expected:    []string{"admserver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreateOptions{
				Services:    tt.services,
				NoServices:  tt.noServices,
				AddServices: tt.addServices,
			}
			result := o.computeServiceList()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("computeServiceList() = %v, want %v", result, tt.expected)
			}
		})
	}
}
