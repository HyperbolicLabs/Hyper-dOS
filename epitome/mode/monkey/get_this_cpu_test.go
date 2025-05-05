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
			want:    &cpuLabels{name: "Intel_Core_i7-10710U_CPU_1.10GHz"},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				line: "model name\t: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics",
			},
			want:    &cpuLabels{name: "AMD_Ryzen_7_8845HS_w_Radeon_780M_Graphics"},
			wantErr: false,
		},
		{
			name: "longer than 63 characters",
			args: args{
				line: "model name\t: Intel(R) Core(TM) i7-10710U CPU thing thing thing thing thing thing @ 1.10GHz",
			},
			want:    &cpuLabels{name: "Intel_Core_i7-10710U_CPU_thing_thing_thing_thing_thing_thing..."},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertCPUModelLineToCPULabels(tt.args.line)
			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCPUModelToCPUName() gotString = %v, want %v", got, tt.want)
			}
		})
	}
}
