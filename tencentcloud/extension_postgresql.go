package tencentcloud

const (
	POSTGRESQL_PAYTYPE_PREPAID  = "prepaid"
	POSTGRESQL_PAYTYPE_POSTPAID = "postpaid"
)

const (
	COMMON_PAYTYPE_PREPAID  = "PREPAID"
	COMMON_PAYTYPE_POSTPAID = "POSTPAID_BY_HOUR"
)

var POSTGRESQL_PAYTYPE = []string{COMMON_PAYTYPE_POSTPAID}

const (
	POSTGRESQL_DB_VERSION_9_3_5 = "9.3.5"
	POSTGRESQL_DB_VERSION_9_5_4 = "9.5.4"
	POSTGRESQL_DB_VERSION_10_4  = "10.4"
)

var POSTSQL_DB_VERSION = []string{POSTGRESQL_DB_VERSION_9_3_5, POSTGRESQL_DB_VERSION_9_5_4, POSTGRESQL_DB_VERSION_10_4}

const (
	POSTGRESQL_DB_CHARSET_UTF8   = "UTF8"
	POSTGRESQL_DB_CHARSET_LATIN1 = "LATIN1"
)

var POSTGRESQL_DB_CHARSET = []string{POSTGRESQL_DB_CHARSET_UTF8, POSTGRESQL_DB_CHARSET_LATIN1}

const (
	POSTGRESQL_STAUTS_RUNNING = "running"
)

var POSTGRESQL_RETRYABLE_STATUS = []string{
	"initing",
	"expanding",
	"switching",
	// deployment changing not exposed at response struct but actually exists
	"deployment changing",
}

const (
	PostgresqlResourceNotFound = "ResourceNotFound"
)
