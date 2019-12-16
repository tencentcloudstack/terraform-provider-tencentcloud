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

var SSL_CERT_TYPE = []string{
	SSL_CERT_TYPE_SERVER,
	SSL_CERT_TYPE_CA,
}
