package hyperutils

import "encoding/json"

func CovertStructToMap(v any) map[string]any {
	var m map[string]any
	data, _ := json.Marshal(v)
	json.Unmarshal(data, &m)
	return m
}
