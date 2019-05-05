package tencentcloud

import (
	"context"
	"fmt"
	"log"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type RedisService struct {
	client *connectivity.TencentCloudClient
}

type TencentCloudRedisDetail struct {
	RedisId    string
	Name       string
	Zone       string
	ProjectId  int64
	Type       string
	MemSize    int64
	Status     string
	VpcId      string
	SubnetId   string
	Ip         string
	Port       int64
	CreateTime string
}

func (me *RedisService) DescribeRedisZoneConfig(ctx context.Context) (sellConfigures []*redis.RegionConf, errRet error) {

	logId := GetLogId(ctx)
	request := redis.NewDescribeProductInfoRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseRedisClient().DescribeProductInfo(request)
	if err != nil {
		errRet = err
		return
	}
	sellConfigures = response.Response.RegionSet
	return
}

func (me *RedisService) DescribeInstances(ctx context.Context, zoneName, searchKey string,
	projectId, needLimit int64) (instances []TencentCloudRedisDetail, errRet error) {

	logId := GetLogId(ctx)

	var zoneId int64 = -1

	if zoneName != "" {
		for id, name := range REDIS_ZONE_ID2NAME {
			if name == zoneName {
				zoneId = id
				break
			}
		}
		if zoneId == -1 {
			errRet = fmt.Errorf("redis instances not support this zone search.")
			return
		}
	}

	listInitSize := map[bool]int64{true: 500, false: needLimit}[needLimit > 500 || needLimit < 1]
	instances = make([]TencentCloudRedisDetail, 0, listInitSize)

	request := redis.NewDescribeInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var limit, offset uint64 = 20, 0
	var leftNumber int64 = 0
	var leftNumberInit bool = false

	request.Limit = &limit
	request.Offset = &offset

needMoreItems:
	if searchKey != "" {
		request.SearchKey = &searchKey
	}
	if projectId >= 0 {
		request.ProjectIds = []*int64{&projectId}
	}
	response, err := me.client.UseRedisClient().DescribeInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if !leftNumberInit {
		leftNumber = *response.Response.TotalCount
	}
	leftNumber = leftNumber - int64(limit)
	offset = offset + limit
	for _, item := range response.Response.InstanceSet {
		if zoneId != -1 && *item.ZoneId != zoneId {
			continue
		}
		var instance TencentCloudRedisDetail
		if REDIS_NAMES[*item.Type] == "" {
			instance.Type = "unknown"
		} else {
			instance.Type = REDIS_NAMES[*item.Type]
		}
		if REDIS_STATUS[*item.Status] == "" {
			instance.Status = "unknown"
		} else {
			instance.Status = REDIS_STATUS[*item.Status]
		}

		if REDIS_ZONE_ID2NAME[*item.ZoneId] == "" {
			instance.Zone = "unknown"
		} else {
			instance.Zone = REDIS_ZONE_ID2NAME[*item.ZoneId]
		}

		instance.CreateTime = *item.Createtime
		instance.Ip = *item.WanIp
		instance.MemSize = int64(*item.Size)
		instance.Name = *item.InstanceName
		instance.Port = *item.Port
		instance.ProjectId = *item.ProjectId
		instance.RedisId = *item.InstanceId
		instance.SubnetId = *item.UniqSubnetId
		instance.VpcId = *item.UniqVpcId
		instances = append(instances, instance)

		if needLimit > 0 && int64(len(instances)) >= needLimit {
			return
		}
	}

	if leftNumber < 0 {
		return
	} else {
		goto needMoreItems
	}
	return
}
