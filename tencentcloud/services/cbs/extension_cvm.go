package cbs

// Only client error can cvm retry, others will directly returns
var CVM_RETRYABLE_ERROR = []string{
	// client
	//"ClientError.NetworkError",
	"ClientError.HttpStatusCodeError",
}
