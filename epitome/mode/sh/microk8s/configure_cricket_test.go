package microk8s

import (
	"os"
	"testing"
)

func Test_upgradeServiceIPRange_EditsFile(t *testing.T) {
	// Setup temp file
	tmpFile, err := os.CreateTemp("", "kube-apiserver-test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content to test file
	testContent := `--service-cluster-ip-range=10.96.5.0/24
--other-flag=value`
	if err := os.WriteFile(tmpFile.Name(), []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}

	// Run command without sudo
	err = upgradeServiceIPRange(false, tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to upgrade service IP range: %v", err)
	}

	// Verify results
	modified, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read modified file: %v", err)
	}

	expected := `--service-cluster-ip-range=10.96.0.0/16
--other-flag=value`
	if string(modified) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(modified))
	}
}

func Test_upgradeNodePortRange(t *testing.T) {
	tests := []struct {
		name        string
		testContent string
		newRange    string
		want        string
		wantErr     bool
	}{
		{
			name: "noop",
			testContent: `--asdf=asdf
			--service-node-port-range=80,443,30000-40000`,
			newRange: `80,443,30000-40000`,
			want: `--asdf=asdf
--service-node-port-range=80,443,30000-40000
`,
			wantErr: false,
		},
		{
			name: "update",
			testContent: `--asdf=asdf
			--service-node-port-range=80,443`,
			newRange: `80,443,30000-40000`,
			want: `--asdf=asdf
--service-node-port-range=80,443,30000-40000
`,
			wantErr: false,
		},
		{
			name: "not present, should add to end",
			testContent: `--asdf=asdf
--other-flag=value`,
			newRange: `80,443,30000-40000`,
			want: `--asdf=asdf
--other-flag=value
--service-node-port-range=80,443,30000-40000
`,
		},
	}

	sudo := false
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "nodeportrange-test")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write test content to test file
			if err := os.WriteFile(tmpFile.Name(), []byte(tt.testContent), 0644); err != nil {
				t.Fatalf("failed to write test content: %v", err)
			}
			if err := upgradeNodePortRange(sudo, tmpFile.Name(), tt.newRange); (err != nil) != tt.wantErr {
				t.Errorf("upgradeNodePortRange() error = %v, wantErr %v", err, tt.wantErr)
			}

			// read the temp file and make sure it matches
			modified, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("failed to read modified file: %v", err)
			}
			if string(modified) != tt.want {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.want, string(modified))
			}
		})
	}
}
