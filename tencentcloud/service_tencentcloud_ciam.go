package tencentcloud

import (
	"context"
	"log"

	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CiamService struct {
	client *connectivity.TencentCloudClient
}

func (me *CiamService) DescribeCiamUserGroupById(ctx context.Context, userStoreId string, userGroupId string) (userGroup *ciam.UserGroup, errRet error) {
	logId := getLogId(ctx)

	request := ciam.NewListUserGroupsRequest()
	request.UserStoreId = &userStoreId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiamClient().ListUserGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Content) < 1 {
		return
	}

	userGroup = response.Response.Content[0]
	return
}

func (me *CiamService) DeleteCiamUserGroupById(ctx context.Context, userStoreId string, userGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := ciam.NewDeleteUserGroupsRequest()
	request.UserStoreId = &userStoreId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiamClient().DeleteUserGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
