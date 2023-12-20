package ciam

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CiamService struct {
	client *connectivity.TencentCloudClient
}

func (me *CiamService) DescribeCiamUserGroupById(ctx context.Context, userStoreId string, userGroupId string) (userGroup *ciam.UserGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ciam.NewListUserGroupsRequest()
	request.UserStoreId = &userStoreId

	filter := &ciam.Filter{
		Key:    helper.String("UserGroupId"),
		Values: []*string{helper.String(userGroupId)},
	}

	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 1
		limit  int64 = 20
	)
	instances := make([]*ciam.UserGroup, 0)
	for {
		page := ciam.Pageable{
			PageSize:   &limit,
			PageNumber: &offset,
		}
		request.Pageable = &page

		response, err := me.client.UseCiamClient().ListUserGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Content) < 1 {
			break
		}
		instances = append(instances, response.Response.Content...)
		if len(response.Response.Content) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	userGroup = instances[0]
	return
}

func (me *CiamService) DeleteCiamUserGroupById(ctx context.Context, userStoreId string, userGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ciam.NewDeleteUserGroupsRequest()
	request.UserStoreId = &userStoreId
	request.UserGroupIds = []*string{&userGroupId}

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

func (me *CiamService) DescribeCiamUserStoreById(ctx context.Context, userPoolId string) (userStore *ciam.UserStore, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ciam.NewListUserStoreRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiamClient().ListUserStore(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.UserStoreSet) < 1 {
		return
	}

	for _, ins := range response.Response.UserStoreSet {
		if *ins.UserStoreId == userPoolId {
			userStore = ins
			break
		}
	}
	return
}

func (me *CiamService) DeleteCiamUserStoreById(ctx context.Context, userPoolId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ciam.NewDeleteUserStoreRequest()
	request.UserPoolId = &userPoolId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiamClient().DeleteUserStore(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
