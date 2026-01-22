package wedata

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

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

func (me *WedataService) DescribeWedataProjectsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.Project, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectIds" {
			request.ProjectIds = v.([]*string)
		}

		if k == "ProjectName" {
			request.ProjectName = v.(*string)
		}

		if k == "Status" {
			request.Status = v.(*int64)
		}

		if k == "ProjectModel" {
			request.ProjectModel = v.(*string)
		}
	}

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 50
	)
	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListProjects(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNumber += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataDataSourcesByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.DataSource, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDataSourcesRequest()
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

		if k == "Name" {
			request.Name = v.(*string)
		}

		if k == "DisplayName" {
			request.DisplayName = v.(*string)
		}

		if k == "Type" {
			request.Type = v.([]*string)
		}

		if k == "Creator" {
			request.Creator = v.(*string)
		}
	}

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 50
	)
	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListDataSources(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNumber += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataProjectById(ctx context.Context, projectId string) (ret *wedatav20250806.Project, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetProjectRequest()
	response := wedatav20250806.NewGetProjectResponse()
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetProject(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata project failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataDataSourceById(ctx context.Context, projectId, datasourceId string) (ret *wedatav20250806.DataSource, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetDataSourceRequest()
	response := wedatav20250806.NewGetDataSourceResponse()
	request.ProjectId = &projectId
	datasourceIdInt64 := helper.StrToInt64Point(datasourceId)
	request.Id = datasourceIdInt64

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetDataSource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata data source failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataProjectMemberById(ctx context.Context, projectId, userUin string) (ret []*wedatav20250806.ProjectUserRole, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListProjectMembersRequest()
	response := wedatav20250806.NewListProjectMembersResponse()
	request.ProjectId = &projectId
	request.UserUin = &userUin

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().ListProjectMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata project member failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data.Items
	return
}

func (me *WedataService) DescribeWedataProjectRolesByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.SystemRole, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListProjectRolesRequest()
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
		if k == "RoleDisplayName" {
			request.RoleDisplayName = v.(*string)
		}
	}

	var (
		pageNumber int64 = 1
		pageSize   int64 = 50
	)
	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListProjectRoles(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNumber += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataTenantRolesByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.SystemRole, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTenantRolesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "RoleDisplayName" {
			request.RoleDisplayName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseWedataV20250806Client().ListTenantRoles(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataResourceGroupById(ctx context.Context, resourceGroupId string) (ret []*wedatav20250806.ExecutorResourceGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListResourceGroupsRequest()
	response := wedatav20250806.NewListResourceGroupsResponse()
	request.Id = &resourceGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().ListResourceGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata resource groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data.Items
	return
}

func (me *WedataService) DescribeWedataResourceGroupToProjectAttachmentById(ctx context.Context, resourceGroupId, projectId string) (ret *wedatav20250806.BindProject, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListResourceGroupsRequest()
	response := wedatav20250806.NewListResourceGroupsResponse()
	request.Id = &resourceGroupId
	request.ProjectIds = helper.Strings([]string{projectId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().ListResourceGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata resource groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	if response.Response.Data.Items == nil {
		return
	}

	item := response.Response.Data.Items[0]
	for _, item := range item.AssociateProjects {
		if item.ProjectId != nil && *item.ProjectId == projectId {
			ret = item
			break
		}
	}

	return
}

func (me *WedataService) DescribeWedataResourceGroupMetricsByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ResourceGroupMetrics, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = wedatav20250806.NewGetResourceGroupMetricsRequest()
		response = wedatav20250806.NewGetResourceGroupMetricsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceGroupId" {
			request.ResourceGroupId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*uint64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*uint64)
		}
		if k == "MetricType" {
			request.MetricType = v.(*string)
		}
		if k == "Granularity" {
			request.Granularity = v.(*uint64)
		}
	}

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetResourceGroupMetrics(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata resource group metrics failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
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

func (me *WedataService) DescribeWedataSqlScriptRunsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.JobDto, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListSQLScriptRunsRequest()
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

		if k == "ScriptId" {
			request.ScriptId = v.(*string)
		}

		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}

		if k == "ExecuteUserUin" {
			request.ExecuteUserUin = v.(*string)
		}

		if k == "StartTime" {
			request.StartTime = v.(*string)
		}

		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseWedataV20250806Client().ListSQLScriptRuns(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataSqlFolderById(ctx context.Context, projectId, folderId string) (ret *wedatav20250806.SQLFolderNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListSQLFolderContentsRequest()
	response := wedatav20250806.NewListSQLFolderContentsResponse()
	request.ProjectId = helper.String(projectId)
	request.ParentFolderPath = helper.String(folderId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().ListSQLFolderContents(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata sql folder contents failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data[0]
	return
}

func (me *WedataService) DescribeWedataGetSqlFolderById(ctx context.Context, projectId, folderId string) (ret *wedatav20250806.SQLFolderNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetSQLFolderRequest()
	response := wedatav20250806.NewGetSQLFolderResponse()
	request.ProjectId = helper.String(projectId)
	request.FolderId = helper.String(folderId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetSQLFolder(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe wedata get sql folder failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataSqlScriptById(ctx context.Context, projectId, scriptId string) (ret *wedatav20250806.SQLScript, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetSQLScriptRequest()
	response := wedatav20250806.NewGetSQLScriptResponse()
	request.ProjectId = &projectId
	request.ScriptId = &scriptId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetSQLScript(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Get wedata sql folder script failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataCodeFolderById(ctx context.Context, projectId, folderId string) (ret *wedatav20250806.CodeFolderNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListCodeFolderContentsRequest()
	response := wedatav20250806.NewListCodeFolderContentsResponse()
	request.ProjectId = &projectId
	request.ParentFolderPath = &folderId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().ListCodeFolderContents(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Get wedata code folder contents failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data[0]
	return
}

func (me *WedataService) DescribeWedataGetCodeFolderById(ctx context.Context, projectId, folderId string) (ret *wedatav20250806.CodeFolderNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetCodeFolderRequest()
	response := wedatav20250806.NewGetCodeFolderResponse()
	request.ProjectId = &projectId
	request.FolderId = &folderId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetCodeFolder(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Get wedata get code folder failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataCodeFileById(ctx context.Context, projectId, codeFileId string) (ret *wedatav20250806.CodeFile, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetCodeFileRequest()
	response := wedatav20250806.NewGetCodeFileResponse()
	request.ProjectId = &projectId
	request.CodeFileId = &codeFileId
	request.IncludeContent = helper.Bool(true)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetCodeFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Get wedata code file failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.Data
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

func (me *WedataService) DescribeWedataListLineageByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.LineageNodeInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListLineageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceUniqueId" {
			request.ResourceUniqueId = v.(*string)
		}
		if k == "ResourceType" {
			request.ResourceType = v.(*string)
		}
		if k == "Direction" {
			request.Direction = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListLineage(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListColumnLineageByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.LineageNodeInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListColumnLineageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TableUniqueId" {
			request.TableUniqueId = v.(*string)
		}
		if k == "Direction" {
			request.Direction = v.(*string)
		}
		if k == "ColumnName" {
			request.ColumnName = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListColumnLineage(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListProcessLineageByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.LineagePair, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListProcessLineageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProcessId" {
			request.ProcessId = v.(*string)
		}
		if k == "ProcessType" {
			request.ProcessType = v.(*string)
		}
		if k == "Platform" {
			request.Platform = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListProcessLineage(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListCatalogByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.CatalogInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListCatalogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ParentCatalogId" {
			request.ParentCatalogId = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListCatalog(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListDatabaseByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.DatabaseInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDatabaseRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CatalogName" {
			request.CatalogName = v.(*string)
		}
		if k == "DatasourceId" {
			request.DatasourceId = v.(*int64)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListDatabase(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListSchemaByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.SchemaInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListSchemaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CatalogName" {
			request.CatalogName = v.(*string)
		}
		if k == "DatasourceId" {
			request.DatasourceId = v.(*int64)
		}
		if k == "DatabaseName" {
			request.DatabaseName = v.(*string)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 200
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListSchema(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataListTableByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TableInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTableRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CatalogName" {
			request.CatalogName = v.(*string)
		}
		if k == "DatasourceId" {
			request.DatasourceId = v.(*int64)
		}
		if k == "DatabaseName" {
			request.DatabaseName = v.(*string)
		}
		if k == "SchemaName" {
			request.SchemaName = v.(*string)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
	}

	var (
		pageNum  int64 = 1
		pageSize int64 = 500
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseWedataV20250806Client().ListTable(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataGetTableByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.TableInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTableRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TableGuid" {
			request.TableGuid = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseWedataV20250806Client().GetTable(request)
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

func (me *WedataService) DescribeWedataGetTableColumnsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.ColumnInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTableColumnsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TableGuid" {
			request.TableGuid = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseWedataV20250806Client().GetTableColumns(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataLineageAttachmentById(ctx context.Context, sourceResourceUniqueId, sourceResourceType, sourcePlatform, targetResourceUniqueId, targetResourceType, targetPlatform, processId, processType, processPlatform string) (has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListLineageRequest()
	response := wedatav20250806.NewListLineageResponse()
	request.ResourceUniqueId = &sourceResourceUniqueId
	request.ResourceType = &sourceResourceType
	request.Platform = &sourcePlatform
	request.Direction = helper.String("OUTPUT")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		pageNum  int64 = 1
		pageSize int64 = 100
		items          = []*wedatav20250806.LineageNodeInfo{}
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWedataV20250806Client().ListLineage(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe list lineage failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		items = append(items, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	for _, item := range items {
		if item.Resource != nil {
			if item.Resource.ResourceUniqueId != nil && *item.Resource.ResourceUniqueId == targetResourceUniqueId &&
				(item.Resource.ResourceType != nil && *item.Resource.ResourceType == targetResourceType) ||
				(item.Resource.ResourceType == nil && targetResourceType == "WEDATA") {
				if item.Relation != nil {
					if item.Relation.Processes != nil && len(item.Relation.Processes) > 0 {
						for _, process := range item.Relation.Processes {
							if process.ProcessId != nil && *process.ProcessId == processId &&
								process.ProcessType != nil && *process.ProcessType == processType &&
								process.Platform != nil && *process.Platform == processPlatform {
								has = true
								return
							}
						}
					}
				}
			}
		}
	}

	return
}

func (me *WedataService) DescribeWedataTriggerWorkflowById(ctx context.Context, projectId, workflowId string) (ret *wedatav20250806.TriggerWorkflowDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetTriggerWorkflowRequest()
	request.ProjectId = &projectId
	request.WorkflowId = &workflowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTriggerWorkflow(request)
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

func (me *WedataService) DescribeWedataTriggerTaskById(ctx context.Context, projectId, taskId string) (ret *wedatav20250806.TriggerTask, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewGetTriggerTaskRequest()
	request.ProjectId = &projectId
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTriggerTask(request)
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

func (me *WedataService) DescribeWedataTriggerTaskCodeByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTriggerTaskCodeResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTriggerTaskCodeRequest()
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

	response, err := me.client.UseWedataV20250806Client().GetTriggerTaskCode(request)
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

func (me *WedataService) DescribeWedataTriggerWorkflowsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerWorkflowInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTriggerWorkflowsRequest()
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
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListTriggerWorkflows(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNumber += 1
	}

	return
}

func (me *WedataService) DescribeWedataTriggerTaskVersionsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerTaskVersion, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTriggerTaskVersionsRequest()
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
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize
		response, err := me.client.UseWedataV20250806Client().ListTriggerTaskVersions(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNumber += 1
	}

	return
}

func (me *WedataService) DescribeWedataTriggerTaskVersionByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTriggerTaskVersionResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTriggerTaskVersionRequest()
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

	response, err := me.client.UseWedataV20250806Client().GetTriggerTaskVersion(request)
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

func (me *WedataService) DescribeWedataUpstreamTriggerTasksByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerTaskDependDto, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListUpstreamTriggerTasksRequest()
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
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize

		response, err := me.client.UseWedataV20250806Client().ListUpstreamTriggerTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}

		if response.Response.Data.Items != nil {
			ret = append(ret, response.Response.Data.Items...)
		}

		if response.Response.Data.TotalPageNumber == nil || pageNumber >= *response.Response.Data.TotalPageNumber {
			break
		}

		pageNumber++
	}

	return
}

func (me *WedataService) DescribeWedataDownstreamTriggerTasksByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerTaskDependDto, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListDownstreamTriggerTasksRequest()
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
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize

		response, err := me.client.UseWedataV20250806Client().ListDownstreamTriggerTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}

		if response.Response.Data.Items != nil {
			ret = append(ret, response.Response.Data.Items...)
		}

		if response.Response.Data.TotalPageNumber == nil || pageNumber >= *response.Response.Data.TotalPageNumber {
			break
		}

		pageNumber++
	}

	return
}

func (me *WedataService) DescribeWedataOpsTriggerWorkflowsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerWorkflowBrief, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListOpsTriggerWorkflowsRequest()
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
			request.Filters = v.([]*wedatav20250806.Filter)
		}
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedatav20250806.OrderField)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize

		response, err := me.client.UseWedataV20250806Client().ListOpsTriggerWorkflows(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}

		if response.Response.Data.Items != nil {
			ret = append(ret, response.Response.Data.Items...)
		}

		if response.Response.Data.TotalPageNumber == nil || pageNumber >= *response.Response.Data.TotalPageNumber {
			break
		}

		pageNumber++
	}

	return
}

func (me *WedataService) DescribeWedataOpsTriggerWorkflowByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetOpsTriggerWorkflowResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetOpsTriggerWorkflowRequest()
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
		if k == "WorkflowExecutionId" {
			request.WorkflowExecutionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetOpsTriggerWorkflow(request)
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

func (me *WedataService) DescribeWedataTriggerWorkflowRunsByFilter(ctx context.Context, param map[string]interface{}) (ret []*wedatav20250806.TriggerWorkflowRunBrief, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListTriggerWorkflowRunsRequest()
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
			request.Filters = v.([]*wedatav20250806.Filter)
		}
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedatav20250806.OrderField)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNumber uint64 = 1
		pageSize   uint64 = 100
	)

	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize

		response, err := me.client.UseWedataV20250806Client().ListTriggerWorkflowRuns(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}

		if response.Response.Data.Items != nil {
			ret = append(ret, response.Response.Data.Items...)
		}

		if response.Response.Data.TotalPageNumber == nil || pageNumber >= *response.Response.Data.TotalPageNumber {
			break
		}

		pageNumber++
	}

	return
}

func (me *WedataService) DescribeWedataTriggerWorkflowRunByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTriggerWorkflowRunResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTriggerWorkflowRunRequest()
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
		if k == "WorkflowExecutionId" {
			request.WorkflowExecutionId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*wedatav20250806.Filter)
		}
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedatav20250806.OrderField)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTriggerWorkflowRun(request)
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

func (me *WedataService) DescribeWedataTriggerTaskRunByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.GetTriggerTaskRunResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewGetTriggerTaskRunRequest()
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
		if k == "TaskExecutionId" {
			request.TaskExecutionId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().GetTriggerTaskRun(request)
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

func (me *WedataService) DescribeWedataQualityRuleById(ctx context.Context, projectId, ruleId string) (ret *wedatav20250806.QualityRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListQualityRulesRequest()
	request.ProjectId = &projectId
	request.Filters = []*wedatav20250806.Filter{
		{
			Name:   helper.String("RuleIds"),
			Values: []*string{&ruleId},
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListQualityRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Data.Items == nil || len(response.Response.Data.Items) == 0 {
		return
	}

	ret = response.Response.Data.Items[0]
	return
}

func (me *WedataService) DescribeWedataQualityRuleGroupById(ctx context.Context, projectId, ruleGroupId string) (ret *wedatav20250806.QualityRuleGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListQualityRuleGroupsRequest()
	request.ProjectId = &projectId
	request.Filters = []*wedatav20250806.Filter{
		{
			Name:   helper.String("RuleGroupId"),
			Values: []*string{&ruleGroupId},
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListQualityRuleGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Data == nil || response.Response.Data.Items == nil || len(response.Response.Data.Items) == 0 {
		return
	}

	ret = response.Response.Data.Items[0]
	return
}

func (me *WedataService) DescribeWedataQualityRuleGroupExecResultsByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListQualityRuleGroupExecResultsByPageResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = wedatav20250806.NewListQualityRuleGroupExecResultsByPageRequest()
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
			request.Filters = v.([]*wedatav20250806.Filter)
		}
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedatav20250806.OrderField)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataV20250806Client().ListQualityRuleGroupExecResultsByPage(request)
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

func (me *WedataService) DescribeWedataQualityRuleTemplatesByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.ListQualityRuleTemplatesResponseParams, errRet error) {
	var (
		logId      = tccommon.GetLogId(ctx)
		request    = wedatav20250806.NewListQualityRuleTemplatesRequest()
		allItems   []*wedatav20250806.QualityRuleTemplate
		pageNumber uint64 = 1
		pageSize   uint64 = 20
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		switch k {
		case "ProjectId":
			request.ProjectId = v.(*string)
		case "OrderFields":
			request.OrderFields = v.([]*wedatav20250806.OrderField)
		case "Filters":
			request.Filters = v.([]*wedatav20250806.Filter)
		}
	}

	// 分页循环获取所有数据
	for {
		request.PageNumber = &pageNumber
		request.PageSize = &pageSize

		ratelimit.Check(request.GetAction())

		response, err := me.client.UseWedataV20250806Client().ListQualityRuleTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Data == nil {
			break
		}

		if response.Response.Data.Items != nil {
			allItems = append(allItems, response.Response.Data.Items...)
		}

		if response.Response.Data.TotalCount == nil ||
			uint64(len(allItems)) >= *response.Response.Data.TotalCount {
			ret = &wedatav20250806.ListQualityRuleTemplatesResponseParams{
				Data: &wedatav20250806.QualityRuleTemplatePage{
					TotalCount: response.Response.Data.TotalCount,
					Items:      allItems,
				},
				RequestId: response.Response.RequestId,
			}
			break
		}

		pageNumber++
	}

	return
}

func (me *WedataService) DescribeWedataAuthorizeDataSourceById(ctx context.Context, dataSourceId string) (ret *wedatav20250806.AuthInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewDescribeDataSourceAuthorityRequest()
	response := wedatav20250806.NewDescribeDataSourceAuthorityResponse()
	request.Id = helper.StrToUint64Point(dataSourceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().DescribeDataSourceAuthority(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataWorkflowPermissionsById(ctx context.Context, projectId, entityId, entityType string) (ret []*wedatav20250806.WorkflowPermission, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListWorkflowPermissionsRequest()
	response := wedatav20250806.NewListWorkflowPermissionsResponse()
	request.ProjectId = &projectId
	request.EntityId = &entityId
	request.EntityType = &entityType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		pageNum  uint64 = 1
		pageSize uint64 = 100
	)

	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWedataV20250806Client().ListWorkflowPermissions(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe workflow permissions failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.Data == nil || len(response.Response.Data.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Items...)
		if len(response.Response.Data.Items) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataCodePermissionsById(ctx context.Context, projectId string) (ret []*wedatav20250806.ExploreFilePrivilegeItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := wedatav20250806.NewListCodePermissionsRequest()
	response := wedatav20250806.NewListCodePermissionsResponse()
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		pageNum  int64 = 1
		pageSize int64 = 100
	)

	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseWedataV20250806Client().ListCodePermissions(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe code permissions failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.Data == nil || response.Response.Data.Rows == nil || len(response.Response.Data.Rows) < 1 {
			break
		}

		ret = append(ret, response.Response.Data.Rows...)
		if len(response.Response.Data.Rows) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataWorkflowMaxPermissionByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.WorkflowMaxPermission, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = wedatav20250806.NewGetMyWorkflowMaxPermissionRequest()
		response = wedatav20250806.NewGetMyWorkflowMaxPermissionResponse()
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
		if k == "EntityId" {
			request.EntityId = v.(*string)
		}
		if k == "EntityType" {
			request.EntityType = v.(*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetMyWorkflowMaxPermission(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe workflow max permission failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataCodeMaxPermissionByFilter(ctx context.Context, param map[string]interface{}) (ret *wedatav20250806.CodeStudioMaxPermission, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = wedatav20250806.NewGetMyCodeMaxPermissionRequest()
		response = wedatav20250806.NewGetMyCodeMaxPermissionResponse()
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
		if k == "ResourceId" {
			request.ResourceId = v.(*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetMyCodeMaxPermission(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe code max permission failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataDataBackfillPlanById(ctx context.Context, projectId, dataBackfillPlanId string) (ret *wedatav20250806.DataBackfill, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = wedatav20250806.NewGetDataBackfillPlanRequest()
		response = wedatav20250806.NewGetDataBackfillPlanResponse()
	)

	request.ProjectId = &projectId
	request.DataBackfillPlanId = &dataBackfillPlanId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseWedataV20250806Client().GetDataBackfillPlan(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe data back fill plan failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.Data
	return
}

func (me *WedataService) WedataOpsAsyncJobRefresh(ctx context.Context, projectId string, asyncId string) resource.StateRefreshFunc {
	var req *wedatav20250806.GetOpsAsyncJobRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}

		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}

			req = wedatav20250806.NewGetOpsAsyncJobRequest()
			req.ProjectId = helper.String(projectId)
			req.AsyncId = helper.String(asyncId)
		}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().GetOpsAsyncJobWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}

		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}

		state := fmt.Sprintf("%v", *resp.Response.Data.Status)
		return resp.Response, state, nil
	}
}
