package bh

import (
	"context"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DasbService struct {
	client *connectivity.TencentCloudClient
}

func (me *DasbService) DescribeDasbAclById(ctx context.Context, aclId string) (acl *dasb.Acl, errRet error) {
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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

func (me *DasbService) DescribeDasbDeviceAccountById(ctx context.Context, deviceAccountId string) (DeviceAccount *dasb.DeviceAccount, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeDeviceAccountsRequest()
	deviceAccountIdInt, _ := strconv.ParseUint(deviceAccountId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceAccountIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeDeviceAccounts(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DeviceAccountSet) != 1 {
		return
	}

	DeviceAccount = response.Response.DeviceAccountSet[0]
	return
}

func (me *DasbService) DeleteDasbDeviceAccountById(ctx context.Context, deviceAccountId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDeleteDeviceAccountsRequest()
	deviceAccountIdInt, _ := strconv.ParseUint(deviceAccountId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceAccountIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteDeviceAccounts(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbDeviceGroupMembersById(ctx context.Context, deviceGroupId string) (DeviceGroupMembers []uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeDeviceGroupMembersRequest()
	deviceGroupIdInt, _ := strconv.ParseUint(deviceGroupId, 10, 64)
	request.Id = &deviceGroupIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeDeviceGroupMembers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DeviceSet) < 1 {
		return
	}

	for _, item := range response.Response.DeviceSet {
		if item.Id != nil {
			DeviceGroupMembers = append(DeviceGroupMembers, *item.Id)
		}
	}

	return
}

func (me *DasbService) DeleteDasbDeviceGroupMembersById(ctx context.Context, deviceGroupId, memberIdSetStr string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDeleteDeviceGroupMembersRequest()
	deviceGroupIdInt, _ := strconv.ParseUint(deviceGroupId, 10, 64)
	request.Id = &deviceGroupIdInt
	memberIdSet := strings.Split(memberIdSetStr, tccommon.COMMA_SP)
	tmpList := make([]*uint64, 0)
	for _, item := range memberIdSet {
		itemInt, _ := strconv.ParseUint(item, 10, 64)
		tmpList = append(tmpList, &itemInt)
	}

	request.MemberIdSet = tmpList

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteDeviceGroupMembers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbUserGroupMembersById(ctx context.Context, userGroupId string) (UserGroupMembers []uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeUserGroupMembersRequest()
	userGroupIdInt, _ := strconv.ParseUint(userGroupId, 10, 64)
	request.Id = &userGroupIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeUserGroupMembers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.UserSet) < 1 {
		return
	}

	for _, item := range response.Response.UserSet {
		if item.Id != nil {
			UserGroupMembers = append(UserGroupMembers, *item.Id)
		}
	}

	return
}

func (me *DasbService) DeleteDasbUserGroupMembersById(ctx context.Context, userGroupId, memberIdSetStr string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDeleteUserGroupMembersRequest()
	userGroupIdInt, _ := strconv.ParseUint(userGroupId, 10, 64)
	request.Id = &userGroupIdInt
	memberIdSet := strings.Split(memberIdSetStr, tccommon.COMMA_SP)
	tmpList := make([]*uint64, 0)
	for _, item := range memberIdSet {
		itemInt, _ := strconv.ParseUint(item, 10, 64)
		tmpList = append(tmpList, &itemInt)
	}

	request.MemberIdSet = tmpList

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteUserGroupMembers(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbResourceById(ctx context.Context, resourceId string) (Resource *dasb.Resource, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeResourcesRequest()
	request.ResourceIds = common.StringPtrs([]string{resourceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ResourceSet) != 1 {
		return
	}

	Resource = response.Response.ResourceSet[0]
	return
}

func (me *DasbService) DescribeDasbDeviceById(ctx context.Context, deviceId string) (device *dasb.Device, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeDevicesRequest()
	deviceIdInt, _ := strconv.ParseUint(deviceId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeDevices(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DeviceSet) != 1 {
		return
	}

	device = response.Response.DeviceSet[0]
	return
}

func (me *DasbService) DeleteDasbDeviceById(ctx context.Context, deviceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDeleteDevicesRequest()
	deviceIdInt, _ := strconv.ParseUint(deviceId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteDevices(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DescribeDasbUserGroupById(ctx context.Context, userGroupId string) (UserGroup *dasb.Group, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDescribeUserGroupsRequest()
	userGroupIdInt, _ := strconv.ParseUint(userGroupId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{userGroupIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DescribeUserGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.GroupSet) != 1 {
		return
	}

	UserGroup = response.Response.GroupSet[0]
	return
}

func (me *DasbService) DeleteDasbUserGroupById(ctx context.Context, userGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewDeleteUserGroupsRequest()
	userGroupIdInt, _ := strconv.ParseUint(userGroupId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{userGroupIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().DeleteUserGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DeleteDasbBindDeviceAccountPrivateKeyById(ctx context.Context, deviceAccountId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewResetDeviceAccountPrivateKeyRequest()
	deviceAccountIdInt, _ := strconv.ParseUint(deviceAccountId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceAccountIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().ResetDeviceAccountPrivateKey(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DasbService) DeleteDasbBindDeviceAccountPasswordById(ctx context.Context, deviceAccountId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dasb.NewResetDeviceAccountPasswordRequest()
	deviceAccountIdInt, _ := strconv.ParseUint(deviceAccountId, 10, 64)
	request.IdSet = common.Uint64Ptrs([]uint64{deviceAccountIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDasbClient().ResetDeviceAccountPassword(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
