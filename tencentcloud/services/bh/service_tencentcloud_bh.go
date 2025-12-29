package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewBhService(client *connectivity.TencentCloudClient) BhService {
	return BhService{client: client}
}

type BhService struct {
	client *connectivity.TencentCloudClient
}

func (me *BhService) DescribeBhAccountGroupsByFilter(ctx context.Context, param map[string]interface{}) (ret []*bhv20230418.AccountGroup, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = bhv20230418.NewDescribeAccountGroupsRequest()
		response = bhv20230418.NewDescribeAccountGroupsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DeepIn" {
			request.DeepIn = v.(*int64)
		}
		if k == "ParentId" {
			request.ParentId = v.(*int64)
		}
		if k == "GroupName" {
			request.GroupName = v.(*string)
		}
		if k == "PageNum" {
			request.PageNum = v.(*int64)
		}
	}

	var (
		pageNum  int64 = 0
		pageSize int64 = 100
	)

	for {
		request.PageNum = &pageNum
		request.PageSize = &pageSize
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBhV20230418Client().DescribeAccountGroups(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.AccountGroupSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe account groups failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.AccountGroupSet) < 1 {
			break
		}

		ret = append(ret, response.Response.AccountGroupSet...)
		if len(response.Response.AccountGroupSet) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *BhService) DescribeBhSourceTypesByFilter(ctx context.Context, param map[string]interface{}) (ret []*bhv20230418.SourceType, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = bhv20230418.NewDescribeSourceTypesRequest()
		response = bhv20230418.NewDescribeSourceTypesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeSourceTypes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.SourceTypeSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe source types failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.SourceTypeSet
	return
}

func (me *BhService) DescribeBhAccessWhiteListRuleById(ctx context.Context, ruleId string) (ret *bhv20230418.AccessWhiteListRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeAccessWhiteListRulesRequest()
	response := bhv20230418.NewDescribeAccessWhiteListRulesResponse()
	request.IdSet = []*uint64{helper.StrToUint64Point(ruleId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeAccessWhiteListRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.AccessWhiteListRuleSet == nil || len(result.Response.AccessWhiteListRuleSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe access white list rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.AccessWhiteListRuleSet[0]
	return
}

func (me *BhService) DescribeBhAccessWhiteListConfigById(ctx context.Context) (ret *bhv20230418.DescribeAccessWhiteListRulesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := bhv20230418.NewDescribeAccessWhiteListRulesRequest()
	response := bhv20230418.NewDescribeAccessWhiteListRulesResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeAccessWhiteListRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe access white list rules failed, Response is nil."))
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

func (me *BhService) DescribeBhDeviceById(ctx context.Context, deviceId string) (ret *bhv20230418.Device, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeDevicesRequest()
	response := bhv20230418.NewDescribeDevicesResponse()
	request.IdSet = []*uint64{helper.StrToUint64Point(deviceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeDevices(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DeviceSet == nil || len(result.Response.DeviceSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe device failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.DeviceSet[0]
	return
}

func (me *BhService) DescribeBhAssetSyncFlagConfigById(ctx context.Context) (ret *bhv20230418.AssetSyncFlags, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeAssetSyncFlagRequest()
	response := bhv20230418.NewDescribeAssetSyncFlagResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeAssetSyncFlag(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.AssetSyncFlags == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe asset sync flag failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.AssetSyncFlags
	return
}

func (me *BhService) DescribeBhResourceById(ctx context.Context, resourceId string) (ret *bhv20230418.Resource, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeResourcesRequest()
	response := bhv20230418.NewDescribeResourcesResponse()
	request.ResourceIds = []*string{&resourceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeResources(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ResourceSet == nil || len(result.Response.ResourceSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe resource failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ResourceSet[0]
	return
}

func (me *BhService) DescribeBhReconnectionSettingConfigById(ctx context.Context) (ret *bhv20230418.SecuritySetting, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeSecuritySettingRequest()
	response := bhv20230418.NewDescribeSecuritySettingResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeSecuritySetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.SecuritySetting == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe security setting failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.SecuritySetting
	return
}

func (me *BhService) DescribeBhDepartments(ctx context.Context) (ret *bhv20230418.Departments, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeDepartmentsRequest()
	response := bhv20230418.NewDescribeDepartmentsResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeDepartments(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Departments == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe departments failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Departments
	return
}

func (me *BhService) DescribeBhUserById(ctx context.Context, userId string) (ret *bhv20230418.User, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeUsersRequest()
	response := bhv20230418.NewDescribeUsersResponse()
	request.IdSet = []*uint64{helper.StrToUint64Point(userId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeUsers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.UserSet == nil || len(result.Response.UserSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe users failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.UserSet[0]
	return
}

func (me *BhService) DescribeBhUserGroupById(ctx context.Context, userGroupId string) (ret *bhv20230418.Group, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeUserGroupsRequest()
	response := bhv20230418.NewDescribeUserGroupsResponse()
	request.IdSet = []*uint64{helper.StrToUint64Point(userGroupId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeUserGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.GroupSet == nil || len(result.Response.GroupSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe user groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.GroupSet[0]
	return
}

func (me *BhService) DescribeBhUserDirectoryById(ctx context.Context, directoryId string) (ret *bhv20230418.UserDirectory, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeUserDirectoryRequest()
	response := bhv20230418.NewDescribeUserDirectoryResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset        uint64 = 0
		limit         uint64 = 100
		directoryList []*bhv20230418.UserDirectory
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBhV20230418Client().DescribeUserDirectory(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.UserDirSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe user directory failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.UserDirSet) < 1 {
			break
		}

		directoryList = append(directoryList, response.Response.UserDirSet...)
		if len(response.Response.UserDirSet) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range directoryList {
		if item.Id != nil && helper.UInt64ToStr(*item.Id) == directoryId {
			ret = item
			return
		}
	}

	return
}

func (me *BhService) DescribeBhDevicesByFilter(ctx context.Context, param map[string]interface{}) (ret []*bhv20230418.Device, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = bhv20230418.NewDescribeDevicesRequest()
		response = bhv20230418.NewDescribeDevicesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "IdSet" {
			request.IdSet = v.([]*uint64)
		}
		if k == "Name" {
			request.Name = v.(*string)
		}
		if k == "Ip" {
			request.Ip = v.(*string)
		}
		if k == "ApCodeSet" {
			request.ApCodeSet = v.([]*string)
		}
		if k == "Kind" {
			request.Kind = v.(*uint64)
		}
		if k == "AuthorizedUserIdSet" {
			request.AuthorizedUserIdSet = v.([]*uint64)
		}
		if k == "ResourceIdSet" {
			request.ResourceIdSet = v.([]*string)
		}
		if k == "KindSet" {
			request.KindSet = v.([]*uint64)
		}
		if k == "ManagedAccount" {
			request.ManagedAccount = v.(*string)
		}
		if k == "DepartmentId" {
			request.DepartmentId = v.(*string)
		}
		if k == "AccountIdSet" {
			request.AccountIdSet = v.([]*uint64)
		}
		if k == "ProviderTypeSet" {
			request.ProviderTypeSet = v.([]*uint64)
		}
		if k == "CloudDeviceStatusSet" {
			request.CloudDeviceStatusSet = v.([]*uint64)
		}
		if k == "TagFilters" {
			request.TagFilters = v.([]*bhv20230418.TagFilter)
		}
		if k == "Filters" {
			request.Filters = v.([]*bhv20230418.Filter)
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
			result, e := me.client.UseBhV20230418Client().DescribeDevices(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.DeviceSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe devices failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.DeviceSet) < 1 {
			break
		}

		ret = append(ret, response.Response.DeviceSet...)
		if len(response.Response.DeviceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BhService) DescribeBhAuthModeSettingConfigById(ctx context.Context) (ret *bhv20230418.SecuritySetting, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bhv20230418.NewDescribeSecuritySettingRequest()
	response := bhv20230418.NewDescribeSecuritySettingResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBhV20230418Client().DescribeSecuritySetting(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.SecuritySetting == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe security setting failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.SecuritySetting
	return
}
