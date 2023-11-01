package tencentcloud

import (
	"context"
	"fmt"
	"log"

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

func (me *DlcService) DeleteDlcUsersToWorkGroupAttachmentById(ctx context.Context, workGroupId string, userId []string) (errRet error) {
	logId := getLogId(ctx)

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

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DeleteUsersFromWorkGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DlcService) DescribeDlcCheckDataEngineImageCanBeRollbackByFilter(ctx context.Context, param map[string]interface{}) (checkDataEngineImageCanBeRollback *dlc.CheckDataEngineImageCanBeRollbackResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

	request := dlc.NewDescribeAdvancedStoreLocationRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeAdvancedStoreLocation(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	storeLocationConfig = response.Response
	return
}

func (me *DlcService) DescribeDlcDescribeUserTypeByFilter(ctx context.Context, param map[string]interface{}) (describeUserType *string, errRet error) {
	var (
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
		request = dlc.NewDescribeUserRolesRequest()
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
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDlcClient().DescribeUserRoles(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.UserRoles) < 1 {
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
		logId   = getLogId(ctx)
		request = dlc.NewDescribeUserInfoRequest()
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

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DescribeUserInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.UserInfo == nil {
		return
	}
	describeUserInfo = response.Response.UserInfo
	return
}

func (me *DlcService) DescribeDlcDataEngineByName(ctx context.Context, dataEngineName string) (dataEngine *dlc.DataEngineInfo, errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDescribeDataEnginesRequest()
	item := &dlc.Filter{
		Name:   helper.String("data-engine-name"),
		Values: []*string{helper.String(dataEngineName)}}
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

func (me *DlcService) DeleteDlcDataEngineByName(ctx context.Context, dataEngineName string) (errRet error) {
	logId := getLogId(ctx)

	request := dlc.NewDeleteDataEngineRequest()
	request.DataEngineNames = []*string{&dataEngineName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().DeleteDataEngine(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *DlcService) DescribeDlcDataEngineById(ctx context.Context, dataEngineId string) (dataEngine *dlc.DataEngineInfo, errRet error) {
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
		logId   = getLogId(ctx)
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDlcClient().UnbindWorkGroupsFromUser(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *DlcService) DescribeDlcUserDataEngineConfigById(ctx context.Context, dataEngineId string) (userDataEngineConfig *dlc.DataEngineConfigInstanceInfo, errRet error) {
	logId := getLogId(ctx)

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
		logId   = getLogId(ctx)
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
