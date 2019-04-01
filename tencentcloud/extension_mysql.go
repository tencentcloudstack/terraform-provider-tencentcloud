package tencentcloud

const (
	ZONE_SELL_STATUS_ONLINE = 1
	ZONE_SELL_STATUS_NEW    = 2
)

var MYSQL_ALLOW_BACKUP_TIME = []string{"02:00-06:00", "06:00-10:00", "10:00-14:00", "14:00-18:00", "18:00-22:00", "22:00-02:00"}

var MYSQL_ALLOW_BACKUP_MODEL = []string{"logical", "physical"}

//Async  task  status,  from  https://cloud.tencent.com/document/api/236/20410
const (
	MYSQL_TASK_STATUS_INITIAL = "INITIAL"
	MYSQL_TASK_STATUS_RUNNING = "RUNNING"
	MYSQL_TASK_STATUS_SUCCESS = "SUCCESS"
	MYSQL_TASK_STATUS_FAILED  = "FAILED"
	MYSQL_TASK_STATUS_REMOVED = "REMOVED"
	MYSQL_TASK_STATUS_PAUSED  = "PAUSED "
)

//default to all host
var MYSQL_DEFAULT_ACCOUNT_HOST = "%"

var MYSQL_DATABASE_PRIVILEGE = []string{"SELECT", "INSERT", "UPDATE", "DELETE",
	"CREATE", "DROP", "REFERENCES", "INDEX",
	"ALTER", "CREATE TEMPORARY TABLES", "LOCK TABLES",
	"EXECUTE", "CREATE VIEW", "SHOW VIEW",
	"CREATE ROUTINE", "ALTER ROUTINE", "EVENT", "TRIGGER"}

var MYSQL_DATABASE_MUST_PRIVILEGE = "SHOW VIEW"

var MYSQL_ROLE_MAP = map[int64]string{
	1: "master",
	2: "ro",
	3: "dr",
}

//mysql available period value
var MYSQL_AVAILABLE_PERIOD = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}

var MYSQL_SUPPORTS_ENGINE = []string{"5.5", "5.6", "5.7"}
