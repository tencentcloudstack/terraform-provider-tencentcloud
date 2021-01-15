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
