package helper

import "fmt"

func PortString(port int) string {
	return fmt.Sprintf(":%d", port)
}
