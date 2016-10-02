package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SmartPrint formats messages as text or JSON with a given severity.
func SmartPrint(severity, m string, jsonOut bool) {
	if jsonOut {
		if severity == "" {
			fmt.Printf(m)
			return
		}

		outMap := make(map[string]string)
		outMap[severity] = m

		b, _ := json.Marshal(outMap)
		fmt.Printf(fmt.Sprintf("%s\n", b))

		return
	}

	if severity == "" {
		fmt.Printf("%s", m)
		return
	}

	fmt.Printf("[%s] %s", strings.ToUpper(severity), m)

	return
}
