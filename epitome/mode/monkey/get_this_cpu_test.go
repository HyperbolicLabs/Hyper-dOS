package monkey

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_convertCPUModelToCPUName(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *cpuLabels
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				line: "model name\t: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz",
			},
			want:    &cpuLabels{name: "Intel_Core_i7-10710U_CPU", clock: "1.10GHz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertCPUModelToCPUName(tt.args.line)
			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCPUModelToCPUName() gotString = %v, want %v", got, tt.want)
			}
		})
	}
}
