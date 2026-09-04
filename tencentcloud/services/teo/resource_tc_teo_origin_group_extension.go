package teo

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strings"
)

func teoOriginGroupRecordsHash(v interface{}) int {
	m, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}

	record := ""
	if rv, ok := m["record"].(string); ok {
		record = rv
	}

	recordType := ""
	if tv, ok := m["type"].(string); ok {
		recordType = tv
	}

	weight := 0
	if wv, ok := m["weight"].(int); ok {
		weight = wv
	}

	private := false
	if pv, ok := m["private"].(bool); ok {
		private = pv
	}

	privateParamsSig := ""
	if raw, ok := m["private_parameters"]; ok {
		if list, ok := raw.([]interface{}); ok {
			pairs := make([]string, 0, len(list))
			for _, item := range list {
				pm, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := pm["name"].(string)
				value, _ := pm["value"].(string)
				pairs = append(pairs, name+"="+value)
			}
			sort.Strings(pairs)
			privateParamsSig = strings.Join(pairs, ",")
		}
	}

	s := fmt.Sprintf("record=%s|type=%s|weight=%d|private=%t|private_parameters=%s", record, recordType, weight, private, privateParamsSig)
	return int(crc32.ChecksumIEEE([]byte(s)))
}
