package sheets

import (
	"testing"

	"github.com/dathan/go-github-license-check/pkg/license"
)

func TestService_Save(t *testing.T) {

	rows1 := append(license.LicenseCheckResults{}, license.LicenseCheckResult{"test1", "test2", "test3", "test4"})
	rows2 := append(license.LicenseCheckResults{}, license.LicenseCheckResult{"test5", "test6", "test7", "test8"})
	type args struct {
		input license.LicenseCheckResults
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{"sheet", args{rows1}, false},
		{"sameSheet", args{rows2}, false},
	}
	s := NewService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Save(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Service.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
