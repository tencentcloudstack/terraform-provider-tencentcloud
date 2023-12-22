package tcr

const (
	TCR_VPC_DNS_STATUS_ENABLED  = "ENABLED"
	TCR_VPC_DNS_STATUS_DISABLED = "DISABLED"
)

// FIXME: temp data, use api request instead if support
var RegionIdMap = map[string]string{
	"ap-guangzhou":       "1",
	"ap-shanghai":        "4",
	"ap-hongkong":        "5",
	"na-toronto":         "6",
	"ap-shanghai-fsi":    "7",
	"ap-beijing":         "8",
	"ap-singapore":       "9",
	"ap-shenzhen-fsi":    "11",
	"ap-guangzhou-open":  "12",
	"ap-shanghai-ysx":    "13",
	"na-siliconvalley":   "15",
	"ap-chengdu":         "16",
	"eu-frankfurt":       "17",
	"ap-seoul":           "18",
	"ap-chongqing":       "19",
	"ap-mumbai":          "21",
	"na-ashburn":         "22",
	"ap-bangkok":         "23",
	"eu-moscow":          "24",
	"ap-tokyo":           "25",
	"ap-jinan-ec":        "31",
	"ap-hangzhou-ec":     "32",
	"ap-nanjing":         "33",
	"ap-fuzhou-ec":       "34",
	"ap-wuhan-ec":        "35",
	"ap-tianjin":         "36",
	"ap-shenzhen":        "37",
	"ap-taipei":          "39",
	"ap-changsha-ec":     "45",
	"ap-beijing-fsi":     "46",
	"ap-shijiazhuang-ec": "53",
	"ap-qingyuan":        "54",
	"ap-hefei-ec":        "55",
	"ap-shenyang-ec":     "56",
	"ap-xian-ec":         "57",
	"ap-xibei-ec":        "58",
	"ap-zhengzhou-ec":    "71",
	"ap-jakarta":         "72",
	"ap-qingyuan-xinan":  "73",
	"sa-saopaulo":        "74",
	"ap-guiyang":         "76",
	"ap-shenzhen-sycft":  "77",
	"ap-shanghai-adc":    "78",
}

const (
	REGISTRY_CHARGE_TYPE_POSTPAID = 0
	REGISTRY_CHARGE_TYPE_PREPAID  = 1
)

const (
	TCR_NAME_PREFIX = "tcr$"
)
