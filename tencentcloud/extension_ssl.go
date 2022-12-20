package tencentcloud

const (
	SSL_MODULE_TYPE = "ssl"
	SSL_WITH_CERT   = "1"

	SSL_STATUS_AVAILABLE = 1
)

const (
	SSL_CERT_TYPE_SERVER = "SVR"
	SSL_CERT_TYPE_CA     = "CA"
)

const (
	CsrTypeOnline = "online"
	CsrTypeParse  = "parse"
)

const (
	DNSAuto = "DNS_AUTO"
	DNS     = "DNS"
	File    = "FILE"
)

const (
	InvalidParam          = "FailedOperation.InvalidParam"
	CertificateNotFound   = "FailedOperation.CertificateNotFound"
	InvalidParameter      = "InvalidParameter"
	InvalidParameterValue = "InvalidParameterValue"
	CertificateInvalid    = "FailedOperation.CertificateInvalid"
)

var CsrTypeArr = []string{
	CsrTypeOnline,
	CsrTypeParse,
}

var VerifyType = []string{
	DNSAuto,
	DNS,
	File,
}

var SSL_CERT_TYPE = []string{
	SSL_CERT_TYPE_SERVER,
	SSL_CERT_TYPE_CA,
}
var DNSPOD_OV_EV_TYPE = []int64{51, 52, 53}
var GEOTRUST_OV_EV_TYPE = []int64{8, 9, 10}
var SECURESITE_OV_EV_TYPE = []int64{3, 4, 5, 6, 7}
var TRUSTASIA_OV_EV_TYPE = []int64{13, 14, 15, 16, 17}
var GLOBALSIGN_OV_EV_TYPE = []int64{18, 19, 20, 21, 22, 23, 24}

func IsContainProductId(productId int64, lists ...[]int64) bool {
	for _, list := range lists {
		for _, item := range list {
			if item == productId {
				return true
			}
		}
	}
	return false
}
