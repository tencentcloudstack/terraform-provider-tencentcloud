package igtm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewIgtmService(client *connectivity.TencentCloudClient) IgtmService {
	return IgtmService{client: client}
}

type IgtmService struct {
	client *connectivity.TencentCloudClient
}

func (me *IgtmService) DescribeIgtmAddressPoolById(ctx context.Context, poolId string) (ret *igtmv20231024.DescribeAddressPoolDetailResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := igtmv20231024.NewDescribeAddressPoolDetailRequest()
	response := igtmv20231024.NewDescribeAddressPoolDetailResponse()
	request.PoolId = helper.StrToInt64Point(poolId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseIgtmV20231024Client().DescribeAddressPoolDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe address pool detail failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *IgtmService) DescribeIgtmAddressPoolListByFilter(ctx context.Context, param map[string]interface{}) (ret []*igtmv20231024.AddressPool, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = igtmv20231024.NewDescribeAddressPoolListRequest()
		response = igtmv20231024.NewDescribeAddressPoolListResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*igtmv20231024.ResourceFilter)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseIgtmV20231024Client().DescribeAddressPoolList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.AddressPoolSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe address pool list failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.AddressPoolSet...)
		if len(response.Response.AddressPoolSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
