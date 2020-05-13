package tencentcloud

const monitorEventTypeStatusChange = "status_change"
const monitorEventTypeAbnormal = "abnormal"

var monitorEventTypes = []string{
	monitorEventTypeStatusChange,
	monitorEventTypeAbnormal,
}

const monitorEventStatusRecover = "recover"
const monitorEventStatusAlarm = "alarm"
const monitorEventStatusNothing = "-"

var monitorEventStatus = []string{
	monitorEventStatusRecover,
	monitorEventStatusAlarm,
	monitorEventStatusNothing,
}

//notify_way
const monitorNotifyWaySMS = "SMS"
const monitorNotifyWaySITE = "SITE"
const monitorNotifyWayEMAIL = "EMAIL"
const monitorNotifyWayCALL = "CALL"
const monitorNotifyWayWECHAT = "WECHAT"

var monitorNotifyWays = []string{
	monitorNotifyWaySMS,
	monitorNotifyWaySITE,
	monitorNotifyWayEMAIL,
	monitorNotifyWayCALL,
	monitorNotifyWayWECHAT,
}

//receiver_type
const monitorReceiverTypeUser = "user"
const monitorReceiverTypeGroup = "group"

var monitorReceiverTypes = []string{
	monitorReceiverTypeUser,
	monitorReceiverTypeGroup,
}

//receive_language
const monitorReceiveLanguageCN = "zh-CN"
const monitorReceiveLanguageUS = "en-US"

var monitorReceiveLanguages = []string{
	monitorReceiveLanguageCN,
	monitorReceiveLanguageUS,
}

/*regions in monitor*/
var MonitorRegionMap = map[string]string{
	"ap-guangzhou":       "gz",
	"ap-shenzhen-fsi":    "szjr",
	"ap-guangzhou-open":  "gzopen",
	"ap-shenzhen":        "szx",
	"ap-shanghai":        "sh",
	"ap-shanghai-fsi":    "shjr",
	"ap-nanjing":         "nj",
	"ap-jinan-ec":        "jnec",
	"ap-hangzhou-ec":     "hzec",
	"ap-fuzhou-ec":       "fzec",
	"ap-beijing":         "bj",
	"ap-tianjin":         "tsn",
	"ap-shijiazhuang-ec": "sjwec",
	"ap-beijing-fsi":     "bjjr",
	"ap-wuhan-ec":        "whec",
	"ap-changsha-ec":     "csec",
	"ap-chengdu":         "cd",
	"ap-chongqing":       "cq",
	"ap-taipei":          "tpe",
	"ap-hongkong":        "hk",
	"ap-singapore":       "sg",
	"ap-bangkok":         "th",
	"ap-mumbai":          "in",
	"ap-seoul":           "kr",
	"ap-tokyo":           "jp",
	"na-siliconvalley":   "usw",
	"na-ashburn":         "use",
	"na-toronto":         "ca",
	"eu-frankfurt":       "de",
	"eu-moscow":          "ru",
}
