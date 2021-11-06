package files

import (
	"reflect"
	"testing"
)

func TestAddPath(t *testing.T) {
	type args struct {
		localPath string
		filenames []string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
	}{
		{
			"AddPath",
			args{
				"static",
				[]string{"test.md", "../../../../../../etc/passwd"},
			},
			[]string{"static/test.md", "static/passwd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPaths := addPath(tt.args.localPath, tt.args.filenames); !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("AddPath() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}
}
