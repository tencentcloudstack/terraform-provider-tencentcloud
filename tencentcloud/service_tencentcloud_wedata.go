package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WedataService struct {
	client *connectivity.TencentCloudClient
}

func (me *WedataService) DescribeWedataRuleTemplateById(ctx context.Context, projectId string, ruleTemplateId string) (ruleTemplate *wedata.RuleTemplate, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
	request.TemplateId = helper.StrToUint64Point(ruleTemplateId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Data != nil {

		ruleTemplate = response.Response.Data
	}

	return
}

func (me *WedataService) DeleteWedataRuleTemplateById(ctx context.Context, projectId, ruleTemplateId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
	request.Ids = []*uint64{helper.StrToUint64Point(ruleTemplateId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataRuleTemplatesByFilter(ctx context.Context, param map[string]interface{}) (ruleTemplates []*wedata.RuleTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeRuleTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*uint64)
		}
		if k == "SourceObjectType" {
			request.SourceObjectType = v.(*uint64)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}
		if k == "SourceEngineTypes" {
			request.SourceEngineTypes = v.([]*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ruleTemplates = response.Response.Data

	return
}

func (me *WedataService) DescribeWedataDataSourceListByFilter(ctx context.Context, param map[string]interface{}) (dataSourceList []*wedata.DataSourceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedata.OrderField)
		}

		if k == "Filters" {
			request.Filters = v.([]*wedata.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNum  uint64 = 0
		pageSize uint64 = 20
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataClient().DescribeDataSourceList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.Rows) < 1 {
			break
		}

		dataSourceList = append(dataSourceList, response.Response.Data.Rows...)
		if len(response.Response.Data.Rows) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataDataSourceInfoListByFilter(ctx context.Context, param map[string]interface{}) (dataSourceInfoList []*wedata.DatasourceBaseInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceInfoListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.(*wedata.Filter)
		}

		if k == "OrderFields" {
			request.OrderFields = v.(*wedata.OrderField)
		}

		if k == "Type" {
			request.Type = v.(*string)
		}

		if k == "DatasourceName" {
			request.DatasourceName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNum  uint64 = 0
		pageSize uint64 = 20
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataClient().DescribeDataSourceInfoList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DatasourceSet) < 1 {
			break
		}

		dataSourceInfoList = append(dataSourceInfoList, response.Response.DatasourceSet...)
		if len(response.Response.DatasourceSet) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataDataSourceWithoutInfoByFilter(ctx context.Context, param map[string]interface{}) (dataSourceWithoutInfo []*wedata.DataSourceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceWithoutInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedata.OrderField)
		}

		if k == "Filters" {
			request.Filters = v.([]*wedata.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeDataSourceWithoutInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data) < 1 {
		return
	}

	dataSourceWithoutInfo = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataDatasourceById(ctx context.Context, datasourceId string) (datasource *wedata.DataSourceInfo, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeDatasourceRequest()
	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.Id = &Id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeDatasource(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	datasource = response.Response.Data
	return
}

func (me *WedataService) DeleteWedataDatasourceById(ctx context.Context, datasourceId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteDataSourcesRequest()
	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.Ids = common.Uint64Ptrs([]uint64{Id})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteDataSources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataFunctionById(ctx context.Context, functionId, funcType, funcName, projectId string) (function *wedata.OrganizationalFunction, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeOrganizationalFunctionsRequest()
	request.Type = common.StringPtr("FUNC_DEVELOP")
	request.Name = &funcName
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeOrganizationalFunctions(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Content) < 1 {
		return
	}

	for _, item := range response.Response.Content {
		if item.FuncId != nil {
			if *item.FuncId == functionId {
				function = item
				break
			}
		}
	}

	return
}

func (me *WedataService) DeleteWedataFunctionById(ctx context.Context, functionId, projectId, clusterIdentifier string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteCustomFunctionRequest()
	request.FunctionId = &functionId
	request.ProjectId = &projectId
	request.ClusterIdentifier = &clusterIdentifier

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteCustomFunction(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataResourceById(ctx context.Context, projectId, filePath, resourceId string) (resourceInfo *wedata.ResourcePathTree, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeResourceManagePathTreesRequest()
	request.ProjectId = &projectId
	request.FilePath = &filePath

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeResourceManagePathTrees(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data) < 1 {
		return
	}

	for _, item := range response.Response.Data {
		if *item.ResourceId == resourceId {
			resourceInfo = item
			break
		}
	}

	return
}

func (me *WedataService) DeleteWedataResourceById(ctx context.Context, projectId, resourceId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteResourceFileRequest()
	request.ProjectId = &projectId
	request.ResourceId = &resourceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteResourceFile(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataScriptById(ctx context.Context, projectId, filePath string) (fileInfo *wedata.UserFileInfo, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewGetFileInfoRequest()
	request.ProjectId = &projectId
	request.FilePath = &filePath

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().GetFileInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	fileInfo = response.Response.UserFileInfo
	return
}

func (me *WedataService) DeleteWedataScriptById(ctx context.Context, projectId, resourceId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteFileRequest()
	request.ProjectId = &projectId
	request.ResourceId = &resourceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteFile(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataDqRuleById(ctx context.Context, projectId, ruleId string) (dqRule *wedata.Rule, errRet error) {
	logId := getLogId(ctx)
	request := wedata.NewDescribeRuleRequest()
	request.ProjectId = &projectId
	ruleIdInt, _ := strconv.ParseUint(ruleId, 10, 64)
	request.RuleId = &ruleIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	dqRule = response.Response.Data
	return
}

func (me *WedataService) DeleteWedataDqRuleById(ctx context.Context, projectId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleRequest()
	request.ProjectId = &projectId
	ruleIdInt, _ := strconv.ParseUint(ruleId, 10, 64)
	request.RuleId = &ruleIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataDqRuleTemplateById(ctx context.Context, projectId, templateId string) (dqRuleTemplate *wedata.RuleTemplate, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeRuleTemplateRequest()
	request.ProjectId = &projectId
	TemplateIdInt, _ := strconv.ParseUint(templateId, 10, 64)
	request.TemplateId = &TemplateIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	dqRuleTemplate = response.Response.Data
	return
}

func (me *WedataService) DeleteWedataDqRuleTemplateById(ctx context.Context, projectId, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleTemplateRequest()
	request.ProjectId = &projectId
	TemplateIdInt, _ := strconv.ParseUint(templateId, 10, 64)
	request.Ids = common.Uint64Ptrs([]uint64{TemplateIdInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataBaselineById(ctx context.Context, projectId, baselineId string) (baseline *wedata.BaselineDetailResponse, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeBaselineByIdRequest()
	request.ProjectId = &projectId
	request.BaselineId = &baselineId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeBaselineById(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	baseline = response.Response.Data
	return
}

func (me *WedataService) DeleteWedataBaselineById(ctx context.Context, projectId, baselineId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteBaselineRequest()
	request.ProjectId = &projectId
	request.BaselineId = &baselineId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteBaseline(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataIntegrationOfflineTaskById(ctx context.Context, projectId, taskId string) (integrationOfflineTask *wedata.DescribeIntegrationTaskResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeIntegrationTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId
	request.TaskId = &taskId
	request.TaskType = common.Uint64Ptr(202)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeIntegrationTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	integrationOfflineTask = response.Response
	return
}

func (me *WedataService) DeleteWedataIntegrationOfflineTaskById(ctx context.Context, projectId, taskId string) (errRet error) {
	logId := getLogId(ctx)
	request := wedata.NewDeleteOfflineTaskRequest()
	request.OperatorName = common.StringPtr("")
	request.ProjectId = &projectId
	request.TaskId = &taskId
	request.VirtualFlag = common.BoolPtr(false)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteOfflineTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataIntegrationRealtimeTaskById(ctx context.Context, projectId, taskId string) (integrationRealtimeTask *wedata.DescribeIntegrationTaskResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeIntegrationTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId
	request.TaskType = common.Uint64Ptr(201)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeIntegrationTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	integrationRealtimeTask = response.Response
	return
}

func (me *WedataService) DeleteWedataIntegrationRealtimeTaskById(ctx context.Context, projectId, taskId string) (errRet error) {
	logId := getLogId(ctx)
	request := wedata.NewDeleteIntegrationTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteIntegrationTask(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataIntegrationTaskNodeById(ctx context.Context, projectId, nodeId string) (integrationTaskNode *wedata.DescribeIntegrationNodeResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeIntegrationNodeRequest()
	request.ProjectId = &projectId
	request.Id = &nodeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeIntegrationNode(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	integrationTaskNode = response.Response
	return
}

func (me *WedataService) DeleteWedataIntegrationTaskNodeById(ctx context.Context, projectId, nodeId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteIntegrationNodeRequest()
	request.ProjectId = &projectId
	request.Id = &nodeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteIntegrationNode(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
