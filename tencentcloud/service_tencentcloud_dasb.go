package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DasbService struct {
	client *connectivity.TencentCloudClient
}

func (me *DasbService) DescribeDasbAclById(ctx context.Context, aclId string) (acl *dasb.Acl, errRet error) {
	logId := getLogId(ctx)
	request := dasb.NewDescribeAclsRequest()
	aclIdInt, _ := strconv.ParseUint(aclId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{aclIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeAcls(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.AclSet) != 1 {
		return
	}

	acl = response.Response.AclSet[0]
	return
}

func (me *DasbService) DeleteDasbAclById(ctx context.Context, aclId string) (errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDeleteAclsRequest()
	aclIdInt, _ := strconv.ParseUint(aclId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{aclIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteAcls(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbCmdTemplateById(ctx context.Context, templateId string) (CmdTemplate *dasb.CmdTemplate, errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDescribeCmdTemplatesRequest()
	templateIdInt, _ := strconv.ParseUint(templateId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{templateIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeCmdTemplates(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.CmdTemplateSet) != 1 {
		return
	}

	CmdTemplate = response.Response.CmdTemplateSet[0]
	return
}

func (me *DasbService) DeleteDasbCmdTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDeleteCmdTemplatesRequest()
	templateIdInt, _ := strconv.ParseUint(templateId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{templateIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteCmdTemplates(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbDeviceGroupById(ctx context.Context, deviceGroupId string) (DeviceGroup *dasb.Group, errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDescribeDeviceGroupsRequest()
	deviceGroupIdInt, _ := strconv.ParseUint(deviceGroupId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceGroupIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeDeviceGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.GroupSet) != 1 {
		return
	}

	DeviceGroup = response.Response.GroupSet[0]
	return
}

func (me *DasbService) DeleteDasbDeviceGroupById(ctx context.Context, deviceGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDeleteDeviceGroupsRequest()
	deviceGroupIdInt, _ := strconv.ParseUint(deviceGroupId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceGroupIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteDeviceGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbUserById(ctx context.Context, userId string) (user *dasb.User, errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDescribeUsersRequest()
	userIdInt, _ := strconv.ParseUint(userId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{userIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeUsers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.UserSet) != 1 {
		return
	}

	user = response.Response.UserSet[0]
	return
}

func (me *DasbService) DeleteDasbUserById(ctx context.Context, userId string) (errRet error) {
	logId := getLogId(ctx)

	request := dasb.NewDeleteUsersRequest()
	userIdInt, _ := strconv.ParseUint(userId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{userIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteUsers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
