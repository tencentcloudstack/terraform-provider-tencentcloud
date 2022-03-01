package tencentcloud

import "encoding/json"

const (
	DDOS_EIP_BIND_STATUS_BINDING   = "BINDING"
	DDOS_EIP_BIND_STATUS_BIND      = "BIND"
	DDOS_EIP_BIND_STATUS_UNBINDING = "UNBINDING"
	DDOS_EIP_BIND_STATUS_UNBIND    = "UNBIND"
)

var DDOS_EIP_BIND_STATUS = []string{DDOS_EIP_BIND_STATUS_BINDING, DDOS_EIP_BIND_STATUS_BIND, DDOS_EIP_BIND_STATUS_UNBINDING, DDOS_EIP_BIND_STATUS_UNBIND}

const (
	DDOS_EIP_BIND_RESOURCE_TYPE_CVM = "cvm"
	DDOS_EIP_BIND_RESOURCE_TYPE_CLB = "clb"
	s
)

var DDOS_EIP_BIND_RESOURCE_TYPE = []string{DDOS_EIP_BIND_RESOURCE_TYPE_CVM, DDOS_EIP_BIND_RESOURCE_TYPE_CLB}

const (
	DDOS_BLACK_WHITE_IP_TYPE_BLACK = "black"
	DDOS_BLACK_WHITE_IP_TYPE_WHITE = "white"
)

func DeltaList(oldInstanceList []interface{}, newInstanceList []interface{}) (increment []string, decrement []string) {
	oldInstanceMaps := make(map[string]int)
	newInstanceMaps := make(map[string]int)
	for _, oldInstance := range oldInstanceList {
		buf, _ := json.Marshal(oldInstance)
		oldInstanceMaps[string(buf)] = 1
	}
	for _, newInstance := range newInstanceList {
		buf, _ := json.Marshal(newInstance)
		newInstanceMaps[string(buf)] = 1
	}

	for _, oldInstance := range oldInstanceList {
		buf, _ := json.Marshal(oldInstance)
		key := string(buf)
		if newInstanceMaps[key] == 0 {
			decrement = append(decrement, key)
		}
	}
	for _, newInstance := range newInstanceList {
		buf, _ := json.Marshal(newInstance)
		key := string(buf)
		if oldInstanceMaps[key] == 0 {
			increment = append(increment, key)
		}
	}
	return
}
