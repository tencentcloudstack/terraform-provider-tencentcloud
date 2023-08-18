package tencentcloud

import (
	"context"
	"log"

	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DlcService struct {
	client *connectivity.TencentCloudClient
}

func (me *DlcService) DescribeDlcWorkGroupById(ctx context.Context, workGroupId string) (workGroup *dlc.WorkGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDescribeWorkGroupsRequest()
	request.WorkGroupId = helper.Int64(helper.StrToInt64(workGroupId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeWorkGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WorkGroupSet) < 1 {
		return
	}

	workGroup = response.Response.WorkGroupSet[0]
	return
}

func (me *DlcService) DeleteDlcWorkGroupById(ctx context.Context, workGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDeleteWorkGroupRequest()
	request.WorkGroupIds = []*int64{helper.Int64(helper.StrToInt64(workGroupId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DeleteWorkGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DlcService) DescribeDlcUserById(ctx context.Context, userId string) (user *dlc.UserInfo, errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDescribeUsersRequest()
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeUsers(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.UserSet) < 1 {
		return
	}

	user = response.Response.UserSet[0]
	return
}

func (me *DlcService) DeleteDlcUserById(ctx context.Context, userId string) (errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDeleteUserRequest()
	request.UserIds = []*string{&userId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DeleteUser(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
