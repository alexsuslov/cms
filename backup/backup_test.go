package backup

import (
	"archive/zip"
	"os"
	"testing"
)

func Test_addFile(t *testing.T) {

	f, err := os.Create("test.zip")
	if err!= nil{
		panic(err)
	}
	defer f.Close()

	type args struct {
		w         *zip.Writer
		file      *os.File
		baseInZip string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"create file",
			args{
				w: zip.NewWriter(f),
				file: TrustOpenFile("test.txt"),
				baseInZip:"",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.args.w.Close()
			if err := addFile(tt.args.w, tt.args.file, tt.args.baseInZip); (err != nil) != tt.wantErr {
				t.Errorf("addFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TrustOpenFile(filename string)*os.File{
	f, _ := os.Open(filename)
	return f
}