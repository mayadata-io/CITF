// +build linux darwin

package system

import (
	"os"
	"strings"
	"testing"
)

func TestBinPathFromPathEnv(t *testing.T) {
	type args struct {
		BinName string
	}

	tests := []struct {
		Name string
		Args args
	}{
		{
			Name: "check path of minikube",
			Args: args{
				BinName: "minikube",
			},
		},
		{
			Name: "check path of kubectl",
			Args: args{
				BinName: "kubectl",
			},
		},
		{
			Name: "check path of some-radom-binary",
			Args: args{
				BinName: "some-radom-binary",
			},
		},
	}

	os.Setenv("PATH", "~/front:"+os.Getenv("PATH")+":~/back")
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Log("os.Getenv(\"PATH\")=", os.Getenv("PATH"))
			// Assumption: `which` is present in the system
			want, _ := ExecCommand("which " + tt.Args.BinName)
			want = strings.TrimSpace(want)

			get, err := BinPathFromPathEnv(tt.Args.BinName)
			if err != nil {
				t.Errorf("error occurred while getting path for %q: %+v", tt.Args.BinName, err)
			}
			get = strings.TrimSpace(get)

			if get != want {
				t.Errorf("BinPathFromPathEnv(%q)=%q; want=%q", tt.Args.BinName, get, want)
			}
		})
	}
}
