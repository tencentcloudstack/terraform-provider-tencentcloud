package connectivity

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

var AllSupportedRegions = []string{regions.Beijing,
	regions.Chengdu,
	regions.Chongqing,
	regions.Guangzhou,
	regions.GuangzhouOpen,
	regions.HongKong,
	regions.Mumbai,
	regions.Seoul,
	regions.Shanghai,
	regions.ShanghaiFSI,
	regions.ShenzhenFSI,
	regions.Singapore,
	regions.Frankfurt,
	regions.Moscow,
	regions.Ashburn,
	regions.SiliconValley,
	regions.Toronto}

var MysqlSupportedRegions = AllSupportedRegions
