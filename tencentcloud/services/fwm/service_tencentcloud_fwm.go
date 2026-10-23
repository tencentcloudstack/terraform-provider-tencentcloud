package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func NewFwmService(client *connectivity.TencentCloudClient) FwmService {
	return FwmService{client: client}
}

func (me *FwmService) DescribeFwmCfwManagedByAppId(ctx context.Context, appId string, memberId string) (member *fwm.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwm.NewDescribeOrganMembersRequest()
	response := fwm.NewDescribeOrganMembersResponse()
	request.Filters = []*fwm.CommonFilter{
		{
			Name:         helper.String("AppId"),
			Values:       helper.Strings([]string{appId}),
			OperatorType: helper.Int64(1),
		},
	}
	request.Limit = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient().DescribeOrganMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm organ members failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Members) < 1 {
		return
	}

	for _, item := range response.Response.Members {
		if item.MemberId != nil && *item.MemberId == memberId {
			member = item
			return
		}
	}

	return
}

func (me *FwmService) ModifyFwmCfwManaged(ctx context.Context, memberId string, cfwManaged int64) error {
	logId := tccommon.GetLogId(ctx)

	request := fwm.NewModifyCfwManagedRequest()
	request.MemberId = helper.String(memberId)
	request.CfwManaged = helper.Int64(cfwManaged)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient().ModifyCfwManaged(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update fwm cfw managed failed, reason:%+v", logId, err)
	}

	return err
}
