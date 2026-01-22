package dlc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DlcService struct {
	client *connectivity.TencentCloudClient
}

func (me *DlcService) DescribeDlcWorkGroupById(ctx context.Context, workGroupId string) (workGroup *dlc.WorkGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeWorkGroupsRequest()
	response := dlc.NewDescribeWorkGroupsResponse()
	request.WorkGroupId = helper.Int64(helper.StrToInt64(workGroupId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeWorkGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.WorkGroupSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc work groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	if len(response.Response.WorkGroupSet) < 1 {
		return
	}

	workGroup = response.Response.WorkGroupSet[0]
	return
}

func (me *DlcService) DeleteDlcWorkGroupById(ctx context.Context, workGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDeleteWorkGroupRequest()
	request.WorkGroupIds = []*int64{helper.Int64(helper.StrToInt64(workGroupId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DeleteWorkGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *DlcService) DescribeDlcUserById(ctx context.Context, userId string) (user *dlc.UserInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeUsersRequest()
	response := dlc.NewDescribeUsersResponse()
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeUsers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

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
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDeleteUserRequest()
	request.UserIds = []*string{&userId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DeleteUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *DlcService) DeleteDlcUsersToWorkGroupAttachmentById(ctx context.Context, workGroupId string, userId []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDeleteUsersFromWorkGroupRequest()
	request.AddInfo = &dlc.UserIdSetOfWorkGroupId{
		WorkGroupId: helper.StrToInt64Point(workGroupId),
		UserIds:     helper.Strings(userId),
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DeleteUsersFromWorkGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc users from work group failed, Response is nil."))
		}

		return nil
	})

	return
}

func (me *DlcService) DescribeDlcCheckDataEngineImageCanBeRollbackByFilter(ctx context.Context, param map[string]interface{}) (checkDataEngineImageCanBeRollback *dlc.CheckDataEngineImageCanBeRollbackResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewCheckDataEngineImageCanBeRollbackRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().CheckDataEngineImageCanBeRollback(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	checkDataEngineImageCanBeRollback = response.Response
	return
}

func (me *DlcService) DescribeDlcStoreLocationConfigById(ctx context.Context) (storeLocationConfig *dlc.DescribeAdvancedStoreLocationResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeAdvancedStoreLocationRequest()
	response := dlc.NewDescribeAdvancedStoreLocationResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeAdvancedStoreLocation(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc advanced store location failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	storeLocationConfig = response.Response
	return
}

func (me *DlcService) DescribeDlcDescribeUserTypeByFilter(ctx context.Context, param map[string]interface{}) (describeUserType *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeUserTypeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "UserId" {
			request.UserId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeUserType(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.UserType == nil {
		return
	}

	describeUserType = response.Response.UserType
	return
}
func (me *DlcService) DescribeDlcDescribeUserRolesByFilter(ctx context.Context, param map[string]interface{}) (describeUserRoles []*dlc.UserRole, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = dlc.NewDescribeUserRolesRequest()
		response = dlc.NewDescribeUserRolesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Fuzzy" {
			request.Fuzzy = v.(*string)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "Sorting" {
			request.Sorting = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit

		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseDlcClient().DescribeUserRoles(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe dlc user roles failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			return
		}

		if len(response.Response.UserRoles) < 1 {
			break
		}

		describeUserRoles = append(describeUserRoles, response.Response.UserRoles...)
		if len(response.Response.UserRoles) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
func (me *DlcService) DescribeDlcDescribeUserInfoByFilter(ctx context.Context, param map[string]interface{}) (describeUserInfo *dlc.UserDetailInfo, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = dlc.NewDescribeUserInfoRequest()
		response = dlc.NewDescribeUserInfoResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "UserId" {
			request.UserId = v.(*string)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*dlc.Filter)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "Sorting" {
			request.Sorting = v.(*string)
		}
	}

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeUserInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc user info failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	describeUserInfo = response.Response.UserInfo
	return
}

func (me *DlcService) DescribeDlcDataEngineByName(ctx context.Context, dataEngineName string) (dataEngine *dlc.DataEngineInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeDataEnginesRequest()
	response := dlc.NewDescribeDataEnginesResponse()
	item := &dlc.Filter{
		Name:   helper.String("data-engine-name"),
		Values: []*string{helper.String(dataEngineName)}}
	request.Filters = []*dlc.Filter{item}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeDataEngines(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DataEngines == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc data engine failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	if len(response.Response.DataEngines) < 1 {
		return
	}

	dataEngine = response.Response.DataEngines[0]
	return
}

func (me *DlcService) DeleteDlcDataEngineByName(ctx context.Context, dataEngineName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDeleteDataEngineRequest()
	request.DataEngineNames = []*string{&dataEngineName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DeleteDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete data engine failed, Response is nil."))
		}

		return nil
	})

	return
}

func (me *DlcService) DescribeDlcDataEngineById(ctx context.Context, dataEngineId string) (dataEngine *dlc.DataEngineInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeDataEnginesRequest()
	item := &dlc.Filter{
		Name:   helper.String("engine-id"),
		Values: []*string{helper.String(dataEngineId)}}
	request.Filters = []*dlc.Filter{item}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeDataEngines(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.DataEngines) < 1 {
		return
	}

	dataEngine = response.Response.DataEngines[0]
	return
}
func (me *DlcService) DescribeDlcCheckDataEngineImageCanBeUpgradeByFilter(ctx context.Context, param map[string]interface{}) (checkDataEngineImageCanBeUpgrade *dlc.CheckDataEngineImageCanBeUpgradeResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewCheckDataEngineImageCanBeUpgradeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().CheckDataEngineImageCanBeUpgrade(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	checkDataEngineImageCanBeUpgrade = response.Response
	return
}
func (me *DlcService) DescribeDlcDataEngineImageVersionsByFilter(ctx context.Context, param map[string]interface{}) (describeDataEngineImageVersions []*dlc.DataEngineImageVersion, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeDataEngineImageVersionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "EngineType" {
			request.EngineType = v.(*string)
		}

		if k == "Sort" {
			request.Sort = v.(*string)
		}

		if k == "Asc" {
			request.Asc = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeDataEngineImageVersions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.ImageParentVersions) < 1 {
		return
	}
	describeDataEngineImageVersions = append(describeDataEngineImageVersions, response.Response.ImageParentVersions...)
	return
}
func (me *DlcService) DescribeDlcDataEnginePythonSparkImagesByFilter(ctx context.Context, param map[string]interface{}) (describeDataEnginePythonSparkImages []*dlc.PythonSparkImage, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeDataEnginePythonSparkImagesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ChildImageVersionId" {
			request.ChildImageVersionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeDataEnginePythonSparkImages(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.PythonSparkImages) < 1 {
		return
	}
	describeDataEnginePythonSparkImages = append(describeDataEnginePythonSparkImages, response.Response.PythonSparkImages...)

	return
}
func (me *DlcService) DescribeDlcDescribeEngineUsageInfoByFilter(ctx context.Context, param map[string]interface{}) (describeEngineUsageInfo *dlc.DescribeEngineUsageInfoResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeEngineUsageInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeEngineUsageInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	describeEngineUsageInfo = response.Response
	return
}
func (me *DlcService) DescribeDlcDescribeWorkGroupInfoByFilter(ctx context.Context, param map[string]interface{}) (describeWorkGroupInfo *dlc.WorkGroupDetailInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeWorkGroupInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkGroupId" {
			request.WorkGroupId = v.(*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*dlc.Filter)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "Sorting" {
			request.Sorting = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeWorkGroupInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.WorkGroupInfo == nil {
		return
	}
	describeWorkGroupInfo = response.Response.WorkGroupInfo

	return
}
func (me *DlcService) DlcRestartDataEngineStateRefreshFunc(dataEngineId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dataEngine, err := me.DescribeDlcDataEngineById(context.Background(), dataEngineId)
		if err != nil {
			return nil, "", err
		}

		request := dlc.NewDescribeDataEngineRequest()
		request.DataEngineName = dataEngine.DataEngineName
		response, err := me.client.UseDlcClient().DescribeDataEngine(request)

		if err != nil {
			return nil, "", err
		}

		if response == nil || response.Response == nil || response.Response.DataEngine == nil {
			return nil, "", fmt.Errorf("not found instance")
		}

		return response.Response.DataEngine, helper.Int64ToStr(*response.Response.DataEngine.State), nil
	}
}

func (me *DlcService) DeleteDlcAttachWorkGroupPolicyAttachmentById(ctx context.Context, workGroupId string, requestSet []*dlc.Policy) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDetachWorkGroupPolicyRequest()
	request.WorkGroupId = helper.StrToInt64Point(workGroupId)
	request.PolicySet = requestSet
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DetachWorkGroupPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *DlcService) DeleteDlcBindWorkGroupsToUserById(ctx context.Context, userId string, workIds []*int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewUnbindWorkGroupsFromUserRequest()
	request.AddInfo = &dlc.WorkGroupIdSetOfUserId{
		UserId:       &userId,
		WorkGroupIds: workIds,
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().UnbindWorkGroupsFromUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Unbind dlc work groups from user failed, Response is nil."))
		}

		return nil
	})

	return
}

func (me *DlcService) DescribeDlcUserDataEngineConfigById(ctx context.Context, dataEngineId string) (userDataEngineConfig *dlc.DataEngineConfigInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeUserDataEngineConfigRequest()
	request.Filters = []*dlc.Filter{
		{Name: helper.String("engine-id"), Values: []*string{&dataEngineId}},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeUserDataEngineConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.DataEngineConfigInstanceInfos) < 1 {
		return
	}

	userDataEngineConfig = response.Response.DataEngineConfigInstanceInfos[0]
	return
}

func (me *DlcService) DescribeDlcCheckDataEngineConfigPairsValidityByFilter(ctx context.Context, param map[string]interface{}) (checkDataEngineConfigPairsValidity *dlc.CheckDataEngineConfigPairsValidityResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewCheckDataEngineConfigPairsValidityRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ChildImageVersionId" {
			request.ChildImageVersionId = v.(*string)
		}
		if k == "DataEngineConfigPairs" {
			request.DataEngineConfigPairs = v.([]*dlc.DataEngineConfigPair)
		}
		if k == "ImageVersionId" {
			request.ImageVersionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().CheckDataEngineConfigPairsValidity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	checkDataEngineConfigPairsValidity = response.Response
	return
}

func (me *DlcService) DescribeDlcDescribeUpdatableDataEnginesByFilter(ctx context.Context, param map[string]interface{}) (describeUpdatableDataEngines []*dlc.DataEngineBasicInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeUpdatableDataEnginesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineConfigCommand" {
			request.DataEngineConfigCommand = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeUpdatableDataEngines(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.DataEngineBasicInfos) < 1 {
		return
	}
	describeUpdatableDataEngines = append(describeUpdatableDataEngines, response.Response.DataEngineBasicInfos...)

	return
}
func (me *DlcService) DescribeDlcDescribeDataEngineEventsByFilter(ctx context.Context, param map[string]interface{}) (describeDataEngineEvents []*dlc.HouseEventsInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeDataEngineEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineName" {
			request.DataEngineName = v.(*string)
		}

		if k == "SessionId" {
			request.SessionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDlcClient().DescribeDataEngineEvents(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Events) < 1 {
			break
		}

		describeDataEngineEvents = append(describeDataEngineEvents, response.Response.Events...)
		if len(response.Response.Events) < int(limit) {
			break
		}
		offset += limit
	}

	return
}

func (me *DlcService) DescribeDlcTaskResultByFilter(ctx context.Context, param map[string]interface{}) (ret *dlc.DescribeTaskResultResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeTaskResultRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}

		if k == "NextToken" {
			request.NextToken = v.(*string)
		}

		if k == "MaxResults" {
			request.MaxResults = v.(*int64)
		}

		if k == "IsTransformDataType" {
			request.IsTransformDataType = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeTaskResult(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *DlcService) DescribeDlcEngineNodeSpecificationsByFilter(ctx context.Context, param map[string]interface{}) (ret *dlc.DescribeEngineNodeSpecResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeEngineNodeSpecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineName" {
			request.DataEngineName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeEngineNodeSpec(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *DlcService) DescribeDlcNativeSparkSessionsByFilter(ctx context.Context, param map[string]interface{}) (ret []*dlc.SparkSessionInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeNativeSparkSessionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}

		if k == "ResourceGroupId" {
			request.ResourceGroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeNativeSparkSessions(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.SparkSessionsList) < 1 {
		return
	}

	ret = response.Response.SparkSessionsList
	return
}

func (me *DlcService) DescribeDlcStandardEngineResourceGroupConfigInformationByFilter(ctx context.Context, param map[string]interface{}) (ret *dlc.DescribeStandardEngineResourceGroupConfigInfoResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeStandardEngineResourceGroupConfigInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}

		if k == "Sorting" {
			request.Sorting = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*dlc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeStandardEngineResourceGroupConfigInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *DlcService) DescribeDlcDataEngineNetworkByFilter(ctx context.Context, param map[string]interface{}) (ret []*dlc.EngineNetworkInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeEngineNetworksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}

		if k == "Sorting" {
			request.Sorting = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*dlc.Filter)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDlcClient().DescribeEngineNetworks(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EngineNetworkInfos) < 1 {
			break
		}

		ret = append(ret, response.Response.EngineNetworkInfos...)
		if len(response.Response.EngineNetworkInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DlcService) DescribeDlcDataEngineSessionParametersByFilter(ctx context.Context, param map[string]interface{}) (ret []*dlc.DataEngineImageSessionParameter, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeDataEngineSessionParametersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}

		if k == "DataEngineName" {
			request.DataEngineName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeDataEngineSessionParameters(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.DataEngineParameters) < 1 {
		return
	}

	ret = response.Response.DataEngineParameters
	return
}

func (me *DlcService) DescribeDlcSessionImageVersionByFilter(ctx context.Context, param map[string]interface{}) (ret []*dlc.EngineSessionImage, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dlc.NewDescribeSessionImageVersionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DataEngineId" {
			request.DataEngineId = v.(*string)
		}

		if k == "FrameworkType" {
			request.FrameworkType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDlcClient().DescribeSessionImageVersion(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.EngineSessionImages) < 1 {
		return
	}

	ret = response.Response.EngineSessionImages
	return
}

func (me *DlcService) DescribeDlcUserVpcConnectionById(ctx context.Context, engineNetworkId, userVpcEndpointId string) (ret *dlc.UserVpcConnectionInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeUserVpcConnectionRequest()
	response := dlc.NewDescribeUserVpcConnectionResponse()
	request.EngineNetworkId = &engineNetworkId
	request.UserVpcEndpointIds = helper.Strings([]string{userVpcEndpointId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeUserVpcConnection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc user vpc connection failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc user vpc connection failed, reason:%+v", logId, errRet)
		return
	}

	if len(response.Response.UserVpcConnectionInfos) != 1 {
		return
	}

	ret = response.Response.UserVpcConnectionInfos[0]
	return
}

func (me *DlcService) DescribeDlcStandardEngineResourceGroupById(ctx context.Context, engineResourceGroupName string) (ret *dlc.StandardEngineResourceGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeStandardEngineResourceGroupsRequest()
	response := dlc.NewDescribeStandardEngineResourceGroupsResponse()
	request.Filters = []*dlc.Filter{
		{
			Name:   helper.String("engine-resource-group-name-unique"),
			Values: helper.Strings([]string{engineResourceGroupName}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeStandardEngineResourceGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc standard engine resource groups failed, reason:%+v", logId, errRet)
		return
	}

	if len(response.Response.UserEngineResourceGroupInfos) != 1 {
		return
	}

	ret = response.Response.UserEngineResourceGroupInfos[0]
	return
}

func (me *DlcService) DescribeDlcDataMaskStrategyById(ctx context.Context, strategyId string) (ret *dlc.DataMaskStrategy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeDataMaskStrategiesRequest()
	response := dlc.NewDescribeDataMaskStrategiesResponse()
	request.Filters = []*dlc.Filter{
		{
			Name:   helper.String("strategy-id"),
			Values: helper.Strings([]string{strategyId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeDataMaskStrategies(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc data mask strategies failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc data mask strategies failed, reason:%+v", logId, errRet)
		return
	}

	if len(response.Response.Strategies) != 1 {
		return
	}

	ret = response.Response.Strategies[0]
	return
}

func (me *DlcService) DescribeDlcAttachDataMaskPolicyById(ctx context.Context, catalog, dataBase, table string) (ret *dlc.TableResponseInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeTableRequest()
	response := dlc.NewDescribeTableResponse()
	request.DatasourceConnectionName = &catalog
	request.DatabaseName = &dataBase
	request.TableName = &table

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeTable(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc table failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc table failed, reason:%+v", logId, errRet)
		return
	}

	ret = response.Response.Table
	return
}

func (me *DlcService) DescribeDlcStandardEngineResourceGroupConfigInfoById(ctx context.Context, engineResourceGroupName string) (ret *dlc.StandardEngineResourceGroupConfigInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeStandardEngineResourceGroupConfigInfoRequest()
	response := dlc.NewDescribeStandardEngineResourceGroupConfigInfoResponse()
	request.Filters = []*dlc.Filter{
		{
			Name:   helper.String("engine-resource-group-name"),
			Values: helper.Strings([]string{engineResourceGroupName}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeStandardEngineResourceGroupConfigInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.StandardEngineResourceGroupConfigInfos == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource group config info failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc standard engine resource group config info failed, reason:%+v", logId, errRet)
		return
	}

	if len(response.Response.StandardEngineResourceGroupConfigInfos) < 1 {
		return
	}

	ret = response.Response.StandardEngineResourceGroupConfigInfos[0]
	return
}

func (me *DlcService) DescribeDlcDatasourceHouseAttachmentById(ctx context.Context, datasourceConnectionName string) (ret *dlc.NetworkConnection, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dlc.NewDescribeNetworkConnectionsRequest()
	response := dlc.NewDescribeNetworkConnectionsResponse()
	request.NetworkConnectionName = helper.String(datasourceConnectionName)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDlcClient().DescribeNetworkConnections(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.NetworkConnectionSet == nil || len(result.Response.NetworkConnectionSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc network connections failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe dlc network connections failed, reason:%+v", logId, errRet)
		return
	}

	ret = response.Response.NetworkConnectionSet[0]
	return
}
