package emotedownloader

import (
	"encoding/json"
	"fmt"
	"math"
)

func humanSize(s uintptr) string {
	unitMap := map[int]string{
		0: "B",
		1: "Ki",
		2: "Mi",
		3: "Gi",
		4: "Ti",
		5: "Pi",
		6: "Ei",
	}

	size := float64(s)
	unitIdx := 0
	for i := range len(unitMap) {
		if d, _ := math.Modf(size / 1024.0); d == 0 {
			unitIdx = i
			break
		}
		unitIdx++
		size /= 1024.0
	}
	return fmt.Sprintf("%.2f %s", size, unitMap[unitIdx])
}

func prettyPrint(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
