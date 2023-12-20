package dcg

// https://cloud.tencent.com/document/product/215/19192
const (
	DCG_NETWORK_TYPE_VPC = "VPC"
	DCG_NETWORK_TYPE_CCN = "CCN"
)
const (
	DCG_GATEWAY_TYPE_NORMAL = "NORMAL"
	DCG_GATEWAY_TYPE_NAT    = "NAT"
)

// https://cloud.tencent.com/document/product/215/30643
const (
	DCG_CCN_ROUTE_TYPE_BGP    = "BGP"
	DCG_CCN_ROUTE_TYPE_STATIC = "STATIC"
)

var DCG_NETWORK_TYPES = []string{DCG_NETWORK_TYPE_VPC, DCG_NETWORK_TYPE_CCN}
var DCG_GATEWAY_TYPES = []string{DCG_GATEWAY_TYPE_NORMAL, DCG_GATEWAY_TYPE_NAT}
