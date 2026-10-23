package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func (me *FwmService) DescribeFwmMemberGroup(ctx context.Context, appId, memberId string) (member *fwm.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwm.NewDescribeOrganMembersRequest()
	filterName := "AppId"
	filterValue := appId
	operatorType := int64(1)
	request.Filters = []*fwm.CommonFilter{
		{
			Name:         &filterName,
			Values:       []*string{&filterValue},
			OperatorType: &operatorType,
		},
	}
	limit := int64(100)
	request.Limit = &limit

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	request.Offset = &offset

	for {
		response, err := me.client.UseFwmClient().DescribeOrganMembersWithContext(ctx, request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		for _, item := range response.Response.Members {
			if item.MemberId != nil && *item.MemberId == memberId {
				member = item
				return
			}
		}

		total := int64(0)
		if response.Response.TotalCount != nil {
			total = *response.Response.TotalCount
		}
		offset += limit
		request.Offset = &offset
		if offset >= total {
			return
		}
	}
}

func (me *FwmService) ModifyFwmMemberGroup(ctx context.Context, userMemberId, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwm.NewModifyMemberGroupRequest()
	request.UserMemberId = helper.String(userMemberId)
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmClient().ModifyMemberGroupWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *FwmService) FwmMemberGroupIdToAppIdMemberId(id string) (appId, memberId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid fwm_member_group id: %s", id)
		return
	}
	appId = parts[0]
	memberId = parts[1]
	return
}
