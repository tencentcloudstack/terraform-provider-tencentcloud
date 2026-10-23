package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func (me *FwmService) DescribeFwmOrganMemberById(ctx context.Context, appId, memberId string) (organMember *fwm.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwm.NewDescribeOrganMembersRequest()
	response := fwm.NewDescribeOrganMembersResponse()

	operatorType := int64(1)
	request.Filters = []*fwm.CommonFilter{
		{
			Name:         common.StringPtr("AppId"),
			Values:       common.StringPtrs([]string{appId}),
			OperatorType: common.Int64Ptr(operatorType),
		},
	}
	request.Limit = common.Int64Ptr(100)
	request.Offset = common.Int64Ptr(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = memberId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient(iacExtInfo).DescribeOrganMembersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Members == nil {
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

	for _, member := range response.Response.Members {
		if member.MemberId != nil && *member.MemberId == memberId {
			organMember = member
			return
		}
	}

	return
}
