package monkey

import (
	"fmt"
	"os"
	"strings"
)

func (agent *agent) getThisCPU() (*cpuLabels, error) {
	// Read CPU info from Linux proc filesystem
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to read cpuinfo: %w", err)
	}

	// Find the first model name (all cores will be the same)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "model name") {
			return convertCPUModelToCPUName(line)
		}
	}

	return nil, fmt.Errorf("model name field not found in /proc/cpuinfo")
}

type cpuLabels struct {
	name  string
	clock string
}

func convertCPUModelToCPUName(line string) (*cpuLabels, error) {
	combo := ""

	// start: "model name\t: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz"
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		combo = strings.TrimSpace(parts[1])
	} else {
		return nil, fmt.Errorf("invalid line format: %s", line)
	}
	// now we have: "Intel (R) Core(TM) i7-10710U CPU @ 1.10GHz"

	// remove all instances of (R) and (TM) from output
	combo = strings.ReplaceAll(combo, "(R)", "")
	combo = strings.ReplaceAll(combo, "(TM)", "")

	name := ""
	clock := ""

	// now we have: "Intel Core i7-10710U CPU @ 1.10GHz"
	// @ is an invalid char in k8s labels, so we will have to break it up
	parts = strings.SplitN(combo, "@", 2)
	if len(parts) == 2 {
		name = strings.TrimSpace(parts[0])
		clock = strings.TrimSpace(parts[1])
	} else {
		return nil, fmt.Errorf("could not split invalid line format into name and clock: %s", line)
	}

	// replace all spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	return &cpuLabels{
		name:  name,
		clock: clock,
	}, nil
}
