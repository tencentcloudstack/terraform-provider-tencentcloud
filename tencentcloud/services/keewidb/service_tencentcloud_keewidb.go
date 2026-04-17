package keewidb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	keewidbv20220308 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/keewidb/v20220308"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type KeewidbService struct {
	client *connectivity.TencentCloudClient
}

func (me *KeewidbService) DescribeKeewidbInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*keewidbv20220308.InstanceInfo, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = keewidbv20220308.NewDescribeInstancesRequest()
		response = keewidbv20220308.NewDescribeInstancesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "InstanceName" {
			request.InstanceName = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
		if k == "UniqVpcIds" {
			request.UniqVpcIds = v.([]*string)
		}
		if k == "UniqSubnetIds" {
			request.UniqSubnetIds = v.([]*string)
		}
		if k == "ProjectIds" {
			request.ProjectIds = v.([]*int64)
		}
		if k == "Status" {
			request.Status = v.([]*int64)
		}
		if k == "BillingMode" {
			request.BillingMode = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*int64)
		}
		if k == "AutoRenew" {
			request.AutoRenew = v.([]*int64)
		}
		if k == "VpcIds" {
			request.VpcIds = v.([]*string)
		}
		if k == "SubnetIds" {
			request.SubnetIds = v.([]*string)
		}
		if k == "SearchKeys" {
			request.SearchKeys = v.([]*string)
		}
		if k == "TagKeys" {
			request.TagKeys = v.([]*string)
		}
		if k == "TagList" {
			request.TagList = v.([]*keewidbv20220308.InstanceTagInfo)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 1000
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseKeewidbV20220308Client().DescribeInstancesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeInstances failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
