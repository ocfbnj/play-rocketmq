package utils

import (
	"encoding/json"
	"fmt"
)

func MarshalJson(obj any) string {
	res, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("<%s>", err.Error())
	}

	return string(res)
}
