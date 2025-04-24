package monkey

import "testing"

func Test_labelsAreGood(t *testing.T) {
	type args struct {
		existingLabels map[string]string
		newLabels      map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				existingLabels: map[string]string{
					"hyperbolic.xyz/cpu-name": "Intel_Core_i7-10710U_CPU",
				},
				newLabels: map[string]string{
					"hyperbolic.xyz/cpu-name": "Intel_Core_i7-10710U_CPU",
				},
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				existingLabels: map[string]string{
					"hyperbolic.xyz/cpu-name": "Intel_Core_i7-10710U_CPU",
				},
				newLabels: map[string]string{
					"hyperbolic.xyz/cpu-name":  "Intel_Core i7-10710U_CPU",
					"hyperbolic.xyz/cpu-clock": "3.6GHz",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := labelsAreGood(tt.args.existingLabels, tt.args.newLabels); got != tt.want {
				t.Errorf("labelsAreGood() = %v, want %v", got, tt.want)
			}
		})
	}
}
