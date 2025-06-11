package monitor

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

// notify_way
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

// receiver_type
const monitorReceiverTypeUser = "user"
const monitorReceiverTypeGroup = "group"

var monitorReceiverTypes = []string{
	monitorReceiverTypeUser,
	monitorReceiverTypeGroup,
}

// receive_language
const monitorReceiveLanguageCN = "zh-CN"
const monitorReceiveLanguageUS = "en-US"

var monitorReceiveLanguages = []string{
	monitorReceiveLanguageCN,
	monitorReceiveLanguageUS,
}

/*regions in monitor*/
// https://tapd.woa.com/qcloud_api/markdown_wikis/show/#1210161711000430909
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
	"ap-qingyuan":        "qy",
	"ap-xibei-ec":        "xbec",
	"ap-hefei-ec":        "hfeec",
	"ap-jakarta":         "jkt",
	"sa-saopaulo":        "sao",
}

var MonitorRegionMapName = map[string]string{
	"-":      "ap-guangzhou",
	"gz":     "ap-guangzhou",
	"szjr":   "ap-shenzhen-fsi",
	"gzopen": "ap-guangzhou-open",
	"szx":    "ap-shenzhen",
	"sh":     "ap-shanghai",
	"shjr":   "ap-shanghai-fsi",
	"nj":     "ap-nanjing",
	"jnec":   "ap-jinan-ec",
	"hzec":   "ap-hangzhou-ec",
	"fzec":   "ap-fuzhou-ec",
	"bj":     "ap-beijing",
	"tsn":    "ap-tianjin",
	"sjwec":  "ap-shijiazhuang-ec",
	"bjjr":   "ap-beijing-fsi",
	"whec":   "ap-wuhan-ec",
	"csec":   "ap-changsha-ec",
	"cd":     "ap-chengdu",
	"cq":     "ap-chongqing",
	"tpe":    "ap-taipei",
	"hk":     "ap-hongkong",
	"sg":     "ap-singapore",
	"th":     "ap-bangkok",
	"in":     "ap-mumbai",
	"kr":     "ap-seoul",
	"jp":     "ap-tokyo",
	"usw":    "na-siliconvalley",
	"use":    "na-ashburn",
	"ca":     "na-toronto",
	"de":     "eu-frankfurt",
	"ru":     "eu-moscow",
	"qy":     "ap-qingyuan",
	"xbec":   "ap-xibei-ec",
	"hfeec":  "ap-hefei-ec",
	"jkt":    "ap-jakarta",
	"sao":    "sa-saopaulo",
}
