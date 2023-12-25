package ssl

const (
	SSL_MODULE_TYPE = "ssl"
	SSL_WITH_CERT   = "1"
)

const (
	SSL_STATUS_PENDING = iota
	SSL_STATUS_AVAILABLE
	SSL_STATUS_REJECTED
	SSL_STATUS_EXPIRED
	SSL_STATUS_DNS_ADDED
	SSL_STATUS_PENDING_SUB
	SSL_STATUS_CANCELING
	SSL_STATUS_CANCELED
	SSL_STATUS_DATA_PENDING
	SSL_STATUS_REVOKING
	SSL_STATUS_REVOKED
	SSL_STATUS_REISSUING
	SSL_STATUS_REVOCATION_PENDING
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
const SSL_ERR_CANCELING = `
	The update of the information field is still ongoing. Please retry the "terraform apply" later and then check whether the update process is complete. 
	For more information, please refer to the documentation: 
	https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/resources/ssl_pay_certificate.`

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

var SslCanCancelStatus = map[uint64]struct{}{
	SSL_STATUS_PENDING:      {},
	SSL_STATUS_DNS_ADDED:    {},
	SSL_STATUS_DATA_PENDING: {},
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
