package monkey

import (
	"fmt"
	"os"
	"regexp"
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
			return convertCPUModelLineToCPULabels(line)
		}
	}

	return nil, fmt.Errorf("model name field not found in /proc/cpuinfo")
}

type cpuLabels struct {
	name string
}

// convertCPUModelLineToCPULabels takes a line from /proc/cpuinfo and converts it
// into kubernetes label format. If the line format is invalid, an error is returned.
// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set
// for example:
// model name	: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz
// returns: &cpuLabels{name: "Intel_Core_i7-10710U_CPU_1.10GHz"}
func convertCPUModelLineToCPULabels(line string) (*cpuLabels, error) {
	// Valid label value:
	//     must be 63 characters or less (can be empty),
	//     unless empty, must begin and end with an alphanumeric character ([a-z0-9A-Z]),
	//     could contain dashes (-), underscores (_), dots (.), and alphanumerics between.
	name := ""

	// start: "model name\t: Intel(R) Core(TM) i7-10710U CPU @ 1.10GHz"
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid line format: %s", line)
	}
	name = strings.TrimSpace(parts[1])
	// now we have: "Intel (R) Core(TM) i7-10710U CPU @ 1.10GHz"

	// remove all instances of (R), (TM), "w/"
	// for example: "model name\t: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics"
	// should return "AMD_Ryzen_7_8845HS_w_Radeon_780M_Graphics"
	name = strings.ReplaceAll(name, "(R)", "")
	name = strings.ReplaceAll(name, "(TM)", "")
	name = strings.ReplaceAll(name, "w/", "w")
	name = strings.ReplaceAll(name, "@", "")

	// remove all special characters other than letters, numbers, dashes, dots, spaces and underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9\.\-_\-_\s]`)
	name = re.ReplaceAllString(name, "")

	// replace multiple spaces with a single space
	// https://stackoverflow.com/a/55437544
	name = strings.Join(strings.Fields(strings.TrimSpace(name)), " ")

	// replace all spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// if longer than 63 characters, truncate
	// and replace the last three chars with ---
	if len(name) > 63 {
		name = name[:63]
		// replace the last three chars with "..."
		name = name[:len(name)-3] + "..."
	}

	return &cpuLabels{
		name: name,
	}, nil
}
