package main

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_writeOnYamlFile(t *testing.T) {
	type args struct {
		fileName string
		data     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "work",
			args: args{
				fileName: "sample1",
				data:     testdata,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeOnYamlFile(tt.args.fileName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("writeOnYamlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
