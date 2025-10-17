package wedata

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WedataService struct {
	client *connectivity.TencentCloudClient
}

func (me *WedataService) DescribeWedataRuleTemplateById(ctx context.Context, projectId string, ruleTemplateId string) (ruleTemplate *wedata.RuleTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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

func (me *WedataService) DescribeWedataOpsWorkflowsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.OpsWorkflow, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListOpsWorkflowsRequest()
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
		if k == "FolderId" {
			request.FolderId = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "OwnerUin" {
			request.OwnerUin = v.(*string)
		}
		if k == "WorkflowType" {
			request.WorkflowType = v.(*string)
		}
		if k == "KeyWord" {
			request.KeyWord = v.(*string)
		}
		if k == "SortItem" {
			request.SortItem = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
		if k == "CreateUserUin" {
			request.CreateUserUin = v.(*string)
		}
		if k == "ModifyTime" {
			request.ModifyTime = v.(*string)
		}
		if k == "CreateTime" {
			request.CreateTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 1 // page number starts from 1
		limit  uint64 = 100
	)
	for {
		request.PageNumber = &offset
		request.PageSize = &limit
		response, err := me.client.UseWedataV20250806Client().ListOpsWorkflows(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *WedataService) DescribeWedataOpsWorkflowByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetOpsWorkflowResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetOpsWorkflowRequest()
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
		if k == "WorkflowId" {
			request.WorkflowId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsWorkflow(request)
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

func (me *WedataService) DescribeWedataTaskInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListTaskInstancesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTaskInstancesRequest()
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
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
		if k == "InstanceType" {
			request.InstanceType = v.(*uint64)
		}
		if k == "InstanceState" {
			request.InstanceState = v.(*string)
		}
		if k == "TaskTypeId" {
			request.TaskTypeId = v.(*uint64)
		}
		if k == "CycleType" {
			request.CycleType = v.(*string)
		}
		if k == "OwnerUin" {
			request.OwnerUin = v.(*string)
		}
		if k == "FolderId" {
			request.FolderId = v.(*string)
		}
		if k == "WorkflowId" {
			request.WorkflowId = v.(*string)
		}
		if k == "ExecutorGroupId" {
			request.ExecutorGroupId = v.(*string)
		}
		if k == "ScheduleTimeFrom" {
			request.ScheduleTimeFrom = v.(*string)
		}
		if k == "ScheduleTimeTo" {
			request.ScheduleTimeTo = v.(*string)
		}
		if k == "StartTimeFrom" {
			request.StartTimeFrom = v.(*string)
		}
		if k == "StartTimeTo" {
			request.StartTimeTo = v.(*string)
		}
		if k == "LastUpdateTimeFrom" {
			request.LastUpdateTimeFrom = v.(*string)
		}
		if k == "LastUpdateTimeTo" {
			request.LastUpdateTimeTo = v.(*string)
		}
		if k == "SortColumn" {
			request.SortColumn = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListTaskInstances(request)
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

func (me *WedataService) DescribeWedataTaskInstanceByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTaskInstanceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTaskInstanceRequest()
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
		if k == "InstanceKey" {
			request.InstanceKey = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTaskInstance(request)
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

func (me *WedataService) DescribeWedataTaskInstanceLogByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTaskInstanceLogResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTaskInstanceLogRequest()
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
		if k == "InstanceKey" {
			request.InstanceKey = v.(*string)
		}
		if k == "LifeRoundNum" {
			request.LifeRoundNum = v.(*uint64)
		}
		if k == "LogLevel" {
			request.LogLevel = v.(*string)
		}
		if k == "NextCursor" {
			request.NextCursor = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTaskInstanceLog(request)
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

func (me *WedataService) DescribeWedataTaskInstanceExecutionsByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListTaskInstanceExecutionsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTaskInstanceExecutionsRequest()
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
		if k == "InstanceKey" {
			request.InstanceKey = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListTaskInstanceExecutions(request)
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

func (me *WedataService) DescribeWedataUpstreamTaskInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListUpstreamTaskInstancesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListUpstreamTaskInstancesRequest()
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
		if k == "InstanceKey" {
			request.InstanceKey = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListUpstreamTaskInstances(request)
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

func (me *WedataService) DescribeWedataDownstreamTaskInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListDownstreamTaskInstancesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDownstreamTaskInstancesRequest()
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
		if k == "InstanceKey" {
			request.InstanceKey = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListDownstreamTaskInstances(request)
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

func (me *WedataService) DescribeWedataOpsTaskOwnerById(ctx context.Context, projectId, taskId string) (ret *wedatav20250806.Task, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetOpsTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataOpsAsyncJobByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetOpsAsyncJobResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetOpsAsyncJobRequest()
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
		if k == "AsyncId" {
			request.AsyncId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsAsyncJob(request)
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

func (me *WedataService) DescribeWedataOpsAlarmRuleById(ctx context.Context, projectId, alarmRuleId string) (ret *wedatav20250806.AlarmRuleData, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetOpsAlarmRuleRequest()
	request.ProjectId = &projectId
	request.AlarmRuleId = &alarmRuleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsAlarmRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataOpsAlarmRulesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListOpsAlarmRulesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListOpsAlarmRulesRequest()
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
		if k == "MonitorObjectType" {
			request.MonitorObjectType = v.(*int64)
		}
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
		if k == "AlarmType" {
			request.AlarmType = v.(*string)
		}
		if k == "AlarmLevel" {
			request.AlarmLevel = v.(*int64)
		}
		if k == "AlarmRecipientId" {
			request.AlarmRecipientId = v.(*string)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "CreateUserUin" {
			request.CreateUserUin = v.(*string)
		}
		if k == "CreateTimeFrom" {
			request.CreateTimeFrom = v.(*string)
		}
		if k == "CreateTimeTo" {
			request.CreateTimeTo = v.(*string)
		}
		if k == "UpdateTimeFrom" {
			request.UpdateTimeFrom = v.(*string)
		}
		if k == "UpdateTimeTo" {
			request.UpdateTimeTo = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListOpsAlarmRules(request)
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

func (me *WedataService) DescribeWedataOpsTaskCodeByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetOpsTaskCodeResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetOpsTaskCodeRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsTaskCode(request)
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

func (me *WedataService) DescribeWedataOpsUpstreamTasksByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListUpstreamOpsTasksResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListUpstreamOpsTasksRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListUpstreamOpsTasks(request)
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

func (me *WedataService) DescribeWedataOpsDownstreamTasksByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListDownstreamOpsTasksResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDownstreamOpsTasksRequest()
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
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListDownstreamOpsTasks(request)
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

func (me *WedataService) DescribeWedataOpsTasksByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListOpsTasksResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListOpsTasksRequest()
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
		if k == "TaskTypeId" {
			request.TaskTypeId = v.(*string)
		}
		if k == "WorkflowId" {
			request.WorkflowId = v.(*string)
		}
		if k == "WorkflowName" {
			request.WorkflowName = v.(*string)
		}
		if k == "OwnerUin" {
			request.OwnerUin = v.(*string)
		}
		if k == "FolderId" {
			request.FolderId = v.(*string)
		}
		if k == "SourceServiceId" {
			request.SourceServiceId = v.(*string)
		}
		if k == "TargetServiceId" {
			request.TargetServiceId = v.(*string)
		}
		if k == "ExecutorGroupId" {
			request.ExecutorGroupId = v.(*string)
		}
		if k == "CycleType" {
			request.CycleType = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListOpsTasks(request)
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

func (me *WedataService) DescribeWedataOpsTaskById(ctx context.Context, projectId, taskId string) (ret *wedatav20250806.Task, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetOpsTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataOpsAlarmMessageByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetAlarmMessageResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetAlarmMessageRequest()
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
		if k == "AlarmMessageId" {
			request.AlarmMessageId = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetAlarmMessage(request)
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

func (me *WedataService) DescribeWedataOpsAlarmMessagesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListAlarmMessagesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListAlarmMessagesRequest()
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
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "AlarmLevel" {
			request.AlarmLevel = v.(*uint64)
		}
		if k == "AlarmRecipientId" {
			request.AlarmRecipientId = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListAlarmMessages(request)
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

func (me *WedataService) DescribeWedataDataBackfillInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListDataBackfillInstancesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDataBackfillInstancesRequest()
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
		if k == "DataBackfillPlanId" {
			request.DataBackfillPlanId = v.(*string)
		}
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListDataBackfillInstances(request)
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

func (me *WedataService) DescribeWedataDataBackfillPlanByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetDataBackfillPlanResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetDataBackfillPlanRequest()
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
		if k == "DataBackfillPlanId" {
			request.DataBackfillPlanId = v.(*string)
		}
		if k == "TimeZone" {
			request.TimeZone = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetDataBackfillPlan(request)
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

func (me *WedataService) DescribeWedataDataSourceListByFilter(ctx context.Context, param map[string]interface{}) (dataSourceList []*wedata.DataSourceInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		pageNum  uint64 = 1
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
		logId   = tccommon.GetLogId(ctx)
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
		pageNum  uint64 = 1
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

func (me *WedataService) DescribeWedataDatasourceById(ctx context.Context, ownerProjectId, datasourceId string) (datasource *wedata.DataSourceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedata.NewDescribeDataSourceListRequest()
	request.PageNumber = common.Uint64Ptr(1)
	request.PageSize = common.Uint64Ptr(1)
	request.Filters = []*wedata.Filter{
		{
			Name:   common.StringPtr("ownerProjectId"),
			Values: common.StringPtrs([]string{ownerProjectId}),
		},
		{
			Name:   common.StringPtr("ID"),
			Values: common.StringPtrs([]string{datasourceId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeDataSourceList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data.Rows) != 1 {
		return
	}

	datasource = response.Response.Data.Rows[0]
	return
}

func (me *WedataService) DeleteWedataDatasourceById(ctx context.Context, ownerProjectId, datasourceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedata.NewDeleteDataSourcesRequest()
	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.Ids = common.Uint64Ptrs([]uint64{Id})
	request.ProjectId = &ownerProjectId

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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

func (me *WedataService) DescribeWedataIntegrationOfflineTaskById(ctx context.Context, projectId, taskId string) (integrationOfflineTask *wedata.DescribeIntegrationTaskResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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

func (me *WedataService) DescribeWedataResourceFilesByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.ResourceFileItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListResourceFilesRequest()
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
		if k == "ResourceName" {
			request.ResourceName = v.(*string)
		}
		if k == "ParentFolderPath" {
			request.ParentFolderPath = v.(*string)
		}
		if k == "CreateUserUin" {
			request.CreateUserUin = v.(*string)
		}
		if k == "ModifyTimeStart" {
			request.ModifyTimeStart = v.(*string)
		}
		if k == "ModifyTimeEnd" {
			request.ModifyTimeEnd = v.(*string)
		}
		if k == "CreateTimeStart" {
			request.CreateTimeStart = v.(*string)
		}
		if k == "CreateTimeEnd" {
			request.CreateTimeEnd = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)
	for {
		request.PageNumber = helper.Uint64(pageNumber)
		request.PageSize = helper.Uint64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListResourceFiles(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataWorkflowFoldersByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.WorkflowFolder, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListWorkflowFoldersRequest()
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
		if k == "ParentFolderPath" {
			request.ParentFolderPath = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = helper.Uint64(pageNumber)
		request.PageSize = helper.Uint64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListWorkflowFolders(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataWorkflowsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.WorkflowInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListWorkflowsRequest()
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
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "ParentFolderPath" {
			request.ParentFolderPath = v.(*string)
		}
		if k == "WorkflowType" {
			request.WorkflowType = v.(*string)
		}
		if k == "BundleId" {
			request.BundleId = v.(*string)
		}
		if k == "OwnerUin" {
			request.OwnerUin = v.(*string)
		}
		if k == "CreateUserUin" {
			request.CreateUserUin = v.(*string)
		}
		if k == "ModifyTime" {
			request.ModifyTime = v.([]*string)
		}
		if k == "CreateTime" {
			request.CreateTime = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber int64 = 1
		pageSize   int64 = 100
	)

	for {
		request.PageNumber = helper.Int64(pageNumber)
		request.PageSize = helper.Int64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListWorkflows(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}
	return
}

func (me *WedataService) DescribeWedataTasksByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TaskBaseAttribute, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTasksRequest()
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
		if k == "TaskName" {
			request.TaskName = v.(*string)
		}
		if k == "WorkflowId" {
			request.WorkflowId = v.(*string)
		}
		if k == "OwnerUin" {
			request.OwnerUin = v.(*string)
		}
		if k == "TaskTypeId" {
			request.TaskTypeId = v.(*int64)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "Submit" {
			request.Submit = v.(*bool)
		}
		if k == "BundleId" {
			request.BundleId = v.(*string)
		}
		if k == "CreateUserUin" {
			request.CreateUserUin = v.(*string)
		}
		if k == "ModifyTime" {
			request.ModifyTime = v.([]*string)
		}
		if k == "CreateTime" {
			request.CreateTime = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber int64 = 1
		pageSize   int64 = 100
	)

	for {
		request.PageNumber = helper.Int64(pageNumber)
		request.PageSize = helper.Int64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataTaskVersionsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TaskVersion, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTaskVersionsRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
		if k == "TaskVersionType" {
			request.TaskVersionType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = helper.Uint64(pageNumber)
		request.PageSize = helper.Uint64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListTaskVersions(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataUpstreamTasksByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TaskDependDto, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListUpstreamTasksRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = helper.Uint64(pageNumber)
		request.PageSize = helper.Uint64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListUpstreamTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataDownstreamTasksByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TaskDependDto, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDownstreamTasksRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = helper.Uint64(pageNumber)
		request.PageSize = helper.Uint64(pageSize)
		response, err := me.client.UseWedataV20250806Client().ListDownstreamTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataTaskCodeByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.TaskCodeResult, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTaskCodeRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTaskCode(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataTaskVersionByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTaskVersionResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTaskVersionRequest()
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
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
		if k == "VersionId" {
			request.VersionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTaskVersion(request)
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

func NewWedataService(client *connectivity.TencentCloudClient) WedataService {
	return WedataService{client: client}
}

func (me *WedataService) DescribeWedataWorkflowFolders(ctx context.Context, projectId, folderId, parentFolderPath string) (folders []*wedatav20250806.WorkflowFolder, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListWorkflowFoldersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ProjectId = helper.String(projectId)
	request.ParentFolderPath = helper.String(parentFolderPath)
	ratelimit.Check(request.GetAction())

	var (
		pageNum  uint64 = 1
		pageSize uint64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListWorkflowFolders(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		folders = append(folders, response.Response.Data.Items...)
		for _, item := range response.Response.Data.Items {
			if folderId != "" && item.FolderId != nil && *item.FolderId == folderId {
				folders = append(folders, item)
			}
		}
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataWorkflowById(ctx context.Context, projectId, workflowId string) (ret *wedatav20250806.WorkflowDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetWorkflowRequest()
	request.ProjectId = helper.String(projectId)
	request.WorkflowId = helper.String(workflowId)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataResourceFileById(ctx context.Context, projectId, resourceId string) (ret *wedatav20250806.ResourceFile, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetResourceFileRequest()
	request.ProjectId = helper.String(projectId)
	request.ResourceId = helper.String(resourceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetResourceFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataResourceFolderById(ctx context.Context, projectId, folderId, parentFolderPath string) (folders []*wedatav20250806.ResourceFolder, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListResourceFoldersRequest()
	request.ProjectId = helper.String(projectId)
	request.ParentFolderPath = helper.String(parentFolderPath)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		pageNum  uint64 = 1
		pageSize uint64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListResourceFolders(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		folders = append(folders, response.Response.Data.Items...)
		for _, item := range response.Response.Data.Items {
			if folderId != "" && item.FolderId != nil && *item.FolderId == folderId {
				folders = append(folders, item)
			}
		}
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}
	return
}

func (me *WedataService) DescribeWedataTaskById(ctx context.Context, projectId, taskId string) (ret *wedatav20250806.GetTaskResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetTaskRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ProjectId = helper.String(projectId)
	request.TaskId = helper.String(taskId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}
