package cynosdb

const (
	ZONE_SELL_STATUS_ONLINE = 0
)

var MYSQL_ALLOW_BACKUP_TIME = []string{"02:00-06:00", "06:00-10:00", "10:00-14:00", "14:00-18:00", "18:00-22:00", "22:00-02:00"}

var MYSQL_ALLOW_BACKUP_MODEL = []string{"logical", "physical"}

// mysql Status	https://cloud.tencent.com/document/api/236/15872
const (
	MYSQL_STATUS_DELIVING  = 0
	MYSQL_STATUS_RUNNING   = 1
	MYSQL_STATUS_ISOLATING = 4
	MYSQL_STATUS_ISOLATED  = 5
	//https://cloud.tencent.com/document/api/236/36197
	//Internal business state , not public
	MYSQL_STATUS_ISOLATED_1 = 6
	MYSQL_STATUS_ISOLATED_2 = 7
)

// Async  task  status,  from  https://cloud.tencent.com/document/api/236/20410
const (
	MYSQL_TASK_STATUS_INITIAL = "INITIAL"
	MYSQL_TASK_STATUS_RUNNING = "RUNNING"
	MYSQL_TASK_STATUS_SUCCESS = "SUCCESS"
	MYSQL_TASK_STATUS_FAILED  = "FAILED"
	MYSQL_TASK_STATUS_REMOVED = "REMOVED"
	MYSQL_TASK_STATUS_PAUSED  = "PAUSED "
)

// default to all host
var MYSQL_DEFAULT_ACCOUNT_HOST = "%"

var MYSQL_GlOBAL_PRIVILEGE = []string{
	"ALTER", "ALTER ROUTINE", "CREATE", "CREATE ROUTINE", "CREATE TEMPORARY TABLES",
	"CREATE USER", "CREATE VIEW", "DELETE", "DROP", "EVENT", "EXECUTE", "INDEX", "INSERT",
	"LOCK TABLES", "PROCESS", "REFERENCES", "RELOAD", "REPLICATION CLIENT",
	"REPLICATION SLAVE", "SELECT", "SHOW DATABASES", "SHOW VIEW", "TRIGGER", "UPDATE",
}
var MYSQL_DATABASE_PRIVILEGE = []string{"SELECT", "INSERT", "UPDATE", "DELETE",
	"CREATE", "DROP", "REFERENCES", "INDEX",
	"ALTER", "CREATE TEMPORARY TABLES", "LOCK TABLES",
	"EXECUTE", "CREATE VIEW", "SHOW VIEW",
	"CREATE ROUTINE", "ALTER ROUTINE", "EVENT", "TRIGGER"}

var MYSQL_TABLE_PRIVILEGE = []string{
	"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "REFERENCES", "INDEX",
	"ALTER", "CREATE VIEW", "SHOW VIEW", "TRIGGER",
}
var MYSQL_COLUMN_PRIVILEGE = []string{
	"SELECT", "INSERT", "UPDATE", "REFERENCES",
}

var MYSQL_DATABASE_MUST_PRIVILEGE = "SHOW VIEW"

var MYSQL_ROLE_MAP = map[int64]string{
	1: "master",
	2: "ro",
	3: "dr",
}

var MysqlDelStates = map[int64]bool{
	MYSQL_STATUS_ISOLATING:  true,
	MYSQL_STATUS_ISOLATED:   true,
	MYSQL_STATUS_ISOLATED_1: true,
	MYSQL_STATUS_ISOLATED_2: true,
}

// mysql available period value
var MYSQL_AVAILABLE_PERIOD = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}

var MYSQL_SUPPORTS_ENGINE = []string{"5.5", "5.6", "5.7", "8.0"}

// automatic renewal status code
const (
	MYSQL_RENEW_NOUSE = 0
	MYSQL_RENEW_OPEN  = 1
	MYSQL_RENEW_CLOSE = 2
)

// type of pay
var (
	MysqlPayByMonth = 0
	MysqlPayByUse   = 1
)

const (
	MYSQL_CHARGE_TYPE_PREPAID  = "PREPAID"
	MYSQL_CHARGE_TYPE_POSTPAID = "POSTPAID"
)

var MYSQL_CHARGE_TYPE = map[int]string{
	MysqlPayByMonth: MYSQL_CHARGE_TYPE_PREPAID,
	MysqlPayByUse:   MYSQL_CHARGE_TYPE_POSTPAID,
}

const (
	MysqlInstanceIdNotFound  = "InvalidParameter.InstanceNotFound"
	MysqlInstanceIdNotFound2 = "InvalidParameter"
	MysqlInstanceIdNotFound3 = "InternalError.DatabaseAccessError"
)

var MYSQL_TASK_STATUS = map[string]int64{
	"UNDEFINED": -1,
	"INITIAL":   0,
	"RUNNING":   1,
	"SUCCEED":   2,
	"FAILED":    3,
	"KILLED":    4,
	"REMOVED":   5,
	"PAUSED":    6,
}

var MYSQL_TASK_TYPES = map[string]int64{
	"ROLLBACK":            1,
	"SQL OPERATION":       2,
	"IMPORT DATA":         3,
	"MODIFY PARAM":        5,
	"INITIAL":             6,
	"REBOOT":              7,
	"OPEN GTID":           8,
	"UPGRADE RO":          9,
	"BATCH ROLLBACK":      10,
	"UPGRADE MASTER":      11,
	"DROP TABLES":         12,
	"SWITCH DR TO MASTER": 13,
}
