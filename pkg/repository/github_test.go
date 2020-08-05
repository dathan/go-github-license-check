package repository

import (
	"testing"
)

func TestNewRepository(t *testing.T) {

	rep := NewRepository()

	tests := []struct {
		name string
	}{
		{"NEWREPO"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := rep.GetLicenses("WeConnect", "spaceman")
			if err != nil {
				t.Fatalf("DEAD: %s\n", err)
			}
			if err := rep.SaveLicenses(res); err != nil {
				t.Fatalf("Save Failed: %s", err.Error())
			}
		})
	}
}
