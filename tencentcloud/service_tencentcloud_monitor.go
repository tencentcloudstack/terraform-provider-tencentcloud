package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"gopkg.in/yaml.v2"
)

type MonitorService struct {
	client *connectivity.TencentCloudClient
}

func (me *MonitorService) CheckCanCreateMysqlROInstance(ctx context.Context, mysqlId string) (can bool, errRet error) {

	logId := getLogId(ctx)

	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		errRet = fmt.Errorf("Can not load  time zone `Asia/Chongqing`, reason %s", err.Error())
		return
	}

	request := monitor.NewGetMonitorDataRequest()

	request.Namespace = helper.String("QCE/CDB")
	request.MetricName = helper.String("RealCapacity")
	request.Period = helper.Uint64(60)

	now := time.Now()
	request.StartTime = helper.String(now.Add(-5 * time.Minute).In(loc).Format("2006-01-02T15:04:05+08:00"))
	request.EndTime = helper.String(now.In(loc).Format("2006-01-02T15:04:05+08:00"))

	request.Instances = []*monitor.Instance{
		{
			Dimensions: []*monitor.Dimension{{
				Name:  helper.String("InstanceId"),
				Value: &mysqlId,
			}},
		},
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().GetMonitorData(request)
	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.DataPoints) == 0 {
		return
	}
	dataPoint := response.Response.DataPoints[0]
	if len(dataPoint.Values) == 0 {
		return
	}
	can = true
	return
}

func (me *MonitorService) FullRegions() (regions []string, errRet error) {
	request := cvm.NewDescribeRegionsRequest()
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err := me.client.UseCvmClient().DescribeRegions(request); err != nil {
			return retryError(err, InternalError)
		} else {
			for _, region := range response.Response.RegionSet {
				regions = append(regions, *region.Region)
			}
		}
		return nil
	}); err != nil {
		errRet = err
		return
	}
	return
}

func (me *MonitorService) DescribePolicyGroupDetailInfo(ctx context.Context, groupId int64) (response *monitor.DescribePolicyGroupInfoResponse, errRet error) {

	var (
		request = monitor.NewDescribePolicyGroupInfoRequest()
		err     error
	)
	request.GroupId = &groupId
	request.Module = helper.String("monitor")

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = me.client.UseMonitorClient().DescribePolicyGroupInfo(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		errRet = err
		return
	}
	return
}

func (me *MonitorService) DescribeAlarmPolicyById(ctx context.Context, policyId string) (info *monitor.AlarmPolicy, errRet error) {

	var (
		request = monitor.NewDescribeAlarmPolicyRequest()
	)
	logId := getLogId(ctx)
	request.Module = helper.String("monitor")
	request.PolicyId = &policyId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DescribeAlarmPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if response.Response.Policy == nil {
		return
	}
	info = response.Response.Policy
	return
}

func (me *MonitorService) DescribeAlarmNoticeById(ctx context.Context, alarmmap map[string]interface{}) (noticeIds []*monitor.AlarmNotice, errRet error) {
	var (
		request  = monitor.NewDescribeAlarmNoticesRequest()
		response *monitor.DescribeAlarmNoticesResponse
		err      error
	)
	request.Module = helper.String("monitor")
	request.PageNumber = helper.IntInt64(1)
	request.PageSize = helper.IntInt64(200)
	request.Order = alarmmap["order"].(*string)
	if v, ok := alarmmap["ownerUid"]; ok {
		request.OwnerUid = v.(*int64)
	}
	if v, ok := alarmmap["name"]; ok {
		request.Name = v.(*string)
	}
	if v, ok := alarmmap["receiver_type"]; ok {
		request.ReceiverType = v.(*string)
	}

	if v, ok := alarmmap["userIdArr"]; ok {
		request.UserIds = v.([]*int64)
	}
	if v, ok := alarmmap["groupArr"]; ok {
		request.GroupIds = v.([]*int64)
	}
	if v, ok := alarmmap["noticeArr"]; ok {
		request.NoticeIds = v.([]*string)
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = me.client.UseMonitorClient().DescribeAlarmNotices(request); err != nil {
			return retryError(err, InternalError)
		}
		noticeIds = response.Response.Notices
		return nil
	}); err != nil {
		return
	}
	return
}

func (me *MonitorService) DescribePolicyGroup(ctx context.Context, groupId int64) (info *monitor.DescribePolicyGroupListGroup, errRet error) {

	var (
		request       = monitor.NewDescribePolicyGroupListRequest()
		offset  int64 = 0
		limit   int64 = 20
		finish  bool
	)
	request.Module = helper.String("monitor")
	request.Offset = &offset
	request.Limit = &limit

	for {
		if finish {
			break
		}
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := me.client.UseMonitorClient().DescribePolicyGroupList(request)
			if err != nil {
				return retryError(err, InternalError)
			}
			if len(response.Response.GroupList) < int(limit) {
				finish = true
			}
			for _, v := range response.Response.GroupList {
				if *v.GroupId == groupId {
					info = v
					return nil
				}
			}
			return nil
		}); err != nil {
			errRet = err
			return
		}
		if info != nil {
			return
		}
		offset = offset + limit
	}
	return
}

func (me *MonitorService) DescribeBindingPolicyObjectList(ctx context.Context, groupId int64) (objects []*monitor.DescribeBindingPolicyObjectListInstance, errRet error) {

	var (
		requestList  = monitor.NewDescribeBindingPolicyObjectListRequest()
		responseList *monitor.DescribeBindingPolicyObjectListResponse
		offset       int64 = 0
		limit        int64 = 100
		finish       bool
		err          error
	)

	requestList.GroupId = &groupId
	requestList.Module = helper.String("monitor")
	requestList.Offset = &offset
	requestList.Limit = &limit

	for {
		if finish {
			break
		}
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(requestList.GetAction())
			if responseList, err = me.client.UseMonitorClient().DescribeBindingPolicyObjectList(requestList); err != nil {
				return retryError(err, InternalError)
			}
			objects = append(objects, responseList.Response.List...)
			if len(responseList.Response.List) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			errRet = err
			return
		}
		offset = offset + limit
	}

	return
}

func (me *MonitorService) DescribeBindingAlarmPolicyObjectList(ctx context.Context, policyId string) (
	objects []*monitor.DescribeBindingPolicyObjectListInstance, errRet error) {

	var (
		requestList  = monitor.NewDescribeBindingPolicyObjectListRequest()
		responseList *monitor.DescribeBindingPolicyObjectListResponse
		offset       int64 = 0
		limit        int64 = 100
		finish       bool
		err          error
	)
	requestList.GroupId = helper.Int64(0)
	requestList.PolicyId = &policyId
	requestList.Module = helper.String("monitor")
	requestList.Offset = &offset
	requestList.Limit = &limit

	for {
		if finish {
			break
		}
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(requestList.GetAction())
			if responseList, err = me.client.UseMonitorClient().DescribeBindingPolicyObjectList(requestList); err != nil {
				return retryError(err, InternalError)
			}
			objects = append(objects, responseList.Response.List...)
			if len(responseList.Response.List) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			errRet = err
			return
		}
		offset = offset + limit
	}

	return
}

// DescribeMonitorTmpInstance tmp
func (me *MonitorService) DescribeMonitorTmpInstance(ctx context.Context, tmpInstanceId string) (tmpInstance *monitor.PrometheusInstancesItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&tmpInstanceId}

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*monitor.PrometheusInstancesItem, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMonitorClient().DescribePrometheusInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	tmpInstance = instances[0]
	return
}

func (me *MonitorService) IsolateMonitorTmpInstanceById(ctx context.Context, tmpInstanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewTerminatePrometheusInstancesRequest()
	request.InstanceIds = []*string{&tmpInstanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().TerminatePrometheusInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DeleteMonitorTmpInstanceById(ctx context.Context, tmpInstanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDestroyPrometheusInstanceRequest()
	request.InstanceId = &tmpInstanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DestroyPrometheusInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorTmpCvmAgent(ctx context.Context, instanceId string, tmpCvmAgentId string) (tmpCvmAgent *monitor.PrometheusAgent, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusAgentsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.AgentIds = []*string{&tmpCvmAgentId}

	response, err := me.client.UseMonitorClient().DescribePrometheusAgents(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AgentSet) < 1 {
		return
	}
	tmpCvmAgent = response.Response.AgentSet[0]
	return
}

func (me *MonitorService) DescribeMonitorTmpScrapeJob(ctx context.Context, tmpScrapeJobId string) (tmpScrapeJob *monitor.PrometheusScrapeJob, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusScrapeJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ids := strings.Split(tmpScrapeJobId, FILED_SP)

	request.JobIds = []*string{&ids[0]}
	request.InstanceId = &ids[1]
	request.AgentId = &ids[2]

	response, err := me.client.UseMonitorClient().DescribePrometheusScrapeJobs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ScrapeJobSet) < 1 {
		return
	}
	tmpScrapeJob = response.Response.ScrapeJobSet[0]
	return
}

func (me *MonitorService) DeleteMonitorTmpScrapeJobById(ctx context.Context, tmpScrapeJobId string) (errRet error) {
	logId := getLogId(ctx)

	ids := strings.Split(tmpScrapeJobId, FILED_SP)
	request := monitor.NewDeletePrometheusScrapeJobsRequest()
	request.JobIds = []*string{&ids[0]}
	request.InstanceId = &ids[1]
	request.AgentId = &ids[2]

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusScrapeJobs(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DeleteMonitorAlarmNoticeById(ctx context.Context, Id string) (errRet error) {
	request := monitor.NewDeleteAlarmNoticesRequest()
	request.Module = helper.String("monitor")
	noticeId := Id
	var n = []*string{&noticeId}
	request.NoticeIds = n

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseMonitorClient().DeleteAlarmNotices(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return
}

func (me *MonitorService) DescribeMonitorTmpExporterIntegration(ctx context.Context, tmpExporterIntegrationId string) (tmpExporterIntegration *monitor.IntegrationConfiguration, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeExporterIntegrationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ids := strings.Split(tmpExporterIntegrationId, FILED_SP)
	if ids[0] != "" {
		request.Name = &ids[0]
	}
	request.InstanceId = &ids[1]
	kubeType, _ := strconv.Atoi(ids[2])
	request.KubeType = helper.IntInt64(kubeType)
	request.ClusterId = &ids[3]
	request.Kind = &ids[4]

	response, err := me.client.UseMonitorClient().DescribeExporterIntegrations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.IntegrationSet) < 1 {
		return
	}
	tmpExporterIntegration = response.Response.IntegrationSet[0]
	return
}

func (me *MonitorService) DeleteMonitorTmpExporterIntegrationById(ctx context.Context, tmpExporterIntegrationId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteExporterIntegrationRequest()
	ids := strings.Split(tmpExporterIntegrationId, FILED_SP)

	request.Name = &ids[0]
	request.InstanceId = &ids[1]
	kubeType, _ := strconv.Atoi(ids[2])
	request.KubeType = helper.IntInt64(kubeType)
	request.ClusterId = &ids[3]
	request.Kind = &ids[4]

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteExporterIntegration(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorTmpAlertRuleById(ctx context.Context, instanceId string, tmpAlertRuleId string) (instance *monitor.PrometheusRuleSet, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeAlertRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.RuleId = &tmpAlertRuleId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DescribeAlertRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AlertRuleSet) < 1 {
		return
	}
	instance = response.Response.AlertRuleSet[0]

	return
}

func (me *MonitorService) DeleteMonitorTmpAlertRule(ctx context.Context, instanceId string, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteAlertRulesRequest()
	request.InstanceId = &instanceId
	request.RuleIds = []*string{&ruleId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteAlertRules(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorRecordingRuleById(ctx context.Context, instanceId string, recordingRuleId string) (instance *monitor.RecordingRuleSet, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeRecordingRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.RuleId = &recordingRuleId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DescribeRecordingRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RecordingRuleSet) < 1 {
		return
	}
	instance = response.Response.RecordingRuleSet[0]

	return
}

func (me *MonitorService) DeleteMonitorRecordingRule(ctx context.Context, instanceId string, recordingRuleId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteRecordingRulesRequest()
	request.InstanceId = &instanceId
	request.RuleIds = []*string{&recordingRuleId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteRecordingRules(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorGrafanaInstance(ctx context.Context, instanceId string) (grafanaInstance *monitor.GrafanaInstanceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeGrafanaInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&instanceId}
	request.Offset = helper.IntInt64(0)
	request.Limit = helper.IntInt64(10)

	response, err := me.client.UseMonitorClient().DescribeGrafanaInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.Instances) < 1 {
		return
	}
	grafanaInstance = response.Response.Instances[0]

	return
}

func (me *MonitorService) DeleteMonitorGrafanaInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteGrafanaInstanceRequest()

	request.InstanceIDs = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteGrafanaInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) CleanGrafanaInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewCleanGrafanaInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().CleanGrafanaInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorGrafanaIntegration(ctx context.Context, integrationId, instanceId string) (grafanaIntegration *monitor.GrafanaIntegrationConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeGrafanaIntegrationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.IntegrationId = &integrationId
	request.InstanceId = &instanceId

	response, err := me.client.UseMonitorClient().DescribeGrafanaIntegrations(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil || len(response.Response.IntegrationSet) < 1 {
		return
	}
	grafanaIntegration = response.Response.IntegrationSet[0]
	return
}

func (me *MonitorService) DeleteMonitorGrafanaIntegrationById(ctx context.Context, integrationId, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteGrafanaIntegrationRequest()

	request.IntegrationId = &integrationId
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteGrafanaIntegration(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorGrafanaNotificationChannel(ctx context.Context, channelId, instanceId string) (grafanaNotificationChannel *monitor.GrafanaChannel, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeGrafanaChannelsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Offset = helper.IntInt64(0)
	request.Limit = helper.IntInt64(10)
	request.ChannelIds = []*string{&channelId}
	request.InstanceId = &instanceId

	response, err := me.client.UseMonitorClient().DescribeGrafanaChannels(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.NotificationChannelSet) < 1 {
		return
	}
	grafanaNotificationChannel = response.Response.NotificationChannelSet[0]
	return
}

func (me *MonitorService) DeleteMonitorGrafanaNotificationChannelById(ctx context.Context, channelId, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteGrafanaNotificationChannelRequest()

	request.ChannelIDs = []*string{&channelId}
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteGrafanaNotificationChannel(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorSsoAccount(ctx context.Context, instanceId, userId string) (ssoAccount *monitor.GrafanaAccountInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribeSSOAccountRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.UserId = &userId

	response, err := me.client.UseMonitorClient().DescribeSSOAccount(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, v := range response.Response.AccountSet {
		if *v.UserId == userId {
			ssoAccount = v
			return
		}
	}

	return
}

func (me *MonitorService) DeleteMonitorSsoAccountById(ctx context.Context, instanceId, userId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeleteSSOAccountRequest()

	request.InstanceId = &instanceId
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeleteSSOAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeMonitorGrafanaPlugin(ctx context.Context, instanceId, pluginId string) (grafanaPlugin *monitor.GrafanaPlugin, errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = monitor.NewDescribeInstalledPluginsRequest()
		response *monitor.DescribeInstalledPluginsResponse
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.PluginId = &pluginId

	response, err := me.client.UseMonitorClient().DescribeInstalledPlugins(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	for _, v := range response.Response.PluginSet {
		if *v.PluginId == pluginId {
			grafanaPlugin = v
			return
		}
	}
	return
}

func (me *MonitorService) DeleteMonitorGrafanaPluginById(ctx context.Context, instanceId, pluginId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewUninstallGrafanaPluginsRequest()

	request.InstanceId = &instanceId
	request.PluginIds = []*string{&pluginId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().UninstallGrafanaPlugins(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeTkeTmpAlertPolicy(ctx context.Context, instanceId, tmpAlertPolicyId string) (tmpAlertPolicy *monitor.PrometheusAlertPolicyItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusAlertPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Filters = append(request.Filters, &monitor.Filter{
		Name:   helper.String("ID"),
		Values: []*string{&tmpAlertPolicyId},
	})

	response, err := me.client.UseMonitorClient().DescribePrometheusAlertPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AlertRules) < 1 {
		return
	}
	tmpAlertPolicy = response.Response.AlertRules[0]
	return
}

func (me *MonitorService) DeleteTkeTmpAlertPolicyById(ctx context.Context, instanceId, tmpAlertPolicyId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeletePrometheusAlertPolicyRequest()
	request.InstanceId = &instanceId
	request.AlertIds = []*string{&tmpAlertPolicyId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusAlertPolicy(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeTmpTkeClusterAgentsById(ctx context.Context, instanceId, clusterId, clusterType string) (agents *monitor.PrometheusAgentOverview, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusClusterAgentsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseMonitorClient().DescribePrometheusClusterAgents(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Agents) < 1 {
			break
		}
		for _, v := range response.Response.Agents {
			if *v.ClusterId == clusterId && *v.ClusterType == clusterType {
				return v, nil
			}
		}
		if len(response.Response.Agents) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	return
}

func (me *MonitorService) DeletePrometheusClusterAgent(ctx context.Context, instanceId, clusterId, clusterType string) (errRet error) {
	logId := getLogId(ctx)
	request := monitor.NewDeletePrometheusClusterAgentRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.Agents = append(request.Agents, &monitor.PrometheusAgentInfo{
		ClusterId:   &clusterId,
		ClusterType: &clusterType,
	})

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusClusterAgent(request)
	if err != nil {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeTkeTmpConfigById(ctx context.Context, configId string) (respParams *monitor.DescribePrometheusConfigResponseParams, errRet error) {
	logId := getLogId(ctx)
	request := monitor.NewDescribePrometheusConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, ids [%s], request body [%s], reason[%s]\n",
				logId, "query object", configId, request.ToJsonString(), errRet.Error())
		}
	}()

	ids, err := me.parseConfigId(configId)
	if err != nil {
		errRet = err
		return
	}

	request.ClusterId = &ids.ClusterId
	request.ClusterType = &ids.ClusterType
	request.InstanceId = &ids.InstanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DescribePrometheusConfig(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail,ids [%s], request body [%s], reason[%s]\n",
			logId, request.GetAction(), configId, request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success,ids [%s], request body [%s], response body [%s]\n",
		logId, request.GetAction(), configId, request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.RequestId == nil {
		return nil, fmt.Errorf("response is invalid,%s", response.ToJsonString())
	}

	respParams = response.Response
	return
}

func (me *MonitorService) DeleteTkeTmpConfigByName(ctx context.Context, configId string, ServiceMonitors []*string, PodMonitors []*string, RawJobs []*string) (errRet error) {
	logId := getLogId(ctx)
	request := monitor.NewDeletePrometheusConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,ids [%s], request body [%s], reason[%s]\n",
				logId, "delete object", configId, request.ToJsonString(), errRet.Error())
		}
	}()

	ids, err := me.parseConfigId(configId)
	if err != nil {
		errRet = err
		return
	}

	request.ClusterId = &ids.ClusterId
	request.ClusterType = &ids.ClusterType
	request.InstanceId = &ids.InstanceId

	if len(ServiceMonitors) > 0 {
		request.ServiceMonitors = ServiceMonitors
	}

	if len(PodMonitors) > 0 {
		request.PodMonitors = PodMonitors
	}

	if len(RawJobs) > 0 {
		request.RawJobs = RawJobs
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, ids [%s], request body [%s], response body [%s]\n",
		logId, request.GetAction(), configId, request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) parseConfigId(configId string) (ret *PrometheusConfigIds, err error) {
	idSplit := strings.Split(configId, FILED_SP)
	if len(idSplit) != 3 {
		return nil, fmt.Errorf("id is broken,%s", configId)
	}

	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]
	if instanceId == "" || clusterType == "" || clusterId == "" {
		return nil, fmt.Errorf("id is broken,%s", configId)
	}

	ret = &PrometheusConfigIds{instanceId, clusterType, clusterId}
	return
}

func (me *MonitorService) DescribeTmpTkeTemplateById(ctx context.Context, templateId string) (template *monitor.PrometheusTemp, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusTempRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&monitor.Filter{
			Name:   helper.String("ID"),
			Values: []*string{&templateId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances := make([]*monitor.PrometheusTemp, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseMonitorClient().DescribePrometheusTemp(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Templates) < 1 {
			break
		}
		instances = append(instances, response.Response.Templates...)
		if len(response.Response.Templates) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}

	for _, v := range instances {
		if *v.TemplateId == templateId {
			template = v
			return
		}
	}

	return
}

func (me *MonitorService) DeleteTmpTkeTemplate(ctx context.Context, tempId string) (errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDeletePrometheusTempRequest()
	request.TemplateId = &tempId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusTemp(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DeletePrometheusRecordRuleYaml(ctx context.Context, id, name string) (errRet error) {
	logId := getLogId(ctx)
	request := monitor.NewDeletePrometheusRecordRuleYamlRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &id
	request.Names = []*string{&name}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DeletePrometheusRecordRuleYaml(request)
	if err != nil {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribePrometheusRecordRuleByName(ctx context.Context, id, name string) (
	ret *monitor.DescribePrometheusRecordRulesResponse, errRet error) {

	logId := getLogId(ctx)
	request := monitor.NewDescribePrometheusRecordRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &id
	if name != "" {
		request.Filters = []*monitor.Filter{
			{
				Name:   helper.String("Name"),
				Values: []*string{&name},
			},
		}
	}

	response, err := me.client.UseMonitorClient().DescribePrometheusRecordRules(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return response, nil
}

func (me *MonitorService) DescribeTkeTmpGlobalNotification(ctx context.Context, instanceId string) (tmpNotification *monitor.PrometheusNotificationItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusGlobalNotificationRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseMonitorClient().DescribePrometheusGlobalNotification(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Notification != nil && response.Response.RequestId != nil {
		tmpNotification = response.Response.Notification
		return
	}

	return
}

func (me *MonitorService) ModifyTkeTmpGlobalNotification(ctx context.Context, instanceId string, notification monitor.PrometheusNotificationItem) (response *monitor.ModifyPrometheusGlobalNotificationResponse, errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewModifyPrometheusGlobalNotificationRequest()
	request.InstanceId = &instanceId
	request.Notification = &notification

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().ModifyPrometheusGlobalNotification(request)
	if err != nil {
		errRet = err
		return nil, err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribePrometheusTempSync(ctx context.Context, templateId string) (targets []*monitor.PrometheusTemplateSyncTarget, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = monitor.NewDescribePrometheusTempSyncRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.TemplateId = &templateId
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMonitorClient().DescribePrometheusTempSync(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success,ids [%s], request body [%s], response body [%s]\n",
		logId, request.GetAction(), templateId, request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.RequestId == nil {
		return nil, fmt.Errorf("response is invalid, %s", response.ToJsonString())
	}

	if len(response.Response.Targets) < 1 {
		return
	}

	targets = response.Response.Targets

	return
}

func (me *MonitorService) DescribeMonitorManageGrafanaAttachmentById(ctx context.Context, instanceId string) (manageGrafanaAttachment *monitor.PrometheusInstancesItem, errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDescribePrometheusInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMonitorClient().DescribePrometheusInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	manageGrafanaAttachment = response.Response.InstanceSet[0]
	return
}

func (me *MonitorService) DeleteMonitorManageGrafanaAttachmentById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	resp, err := me.DescribeMonitorManageGrafanaAttachmentById(ctx, instanceId)
	if err != nil {
		errRet = err
		return
	}

	request := monitor.NewUnbindPrometheusManagedGrafanaRequest()
	request.InstanceId = &instanceId
	request.GrafanaId = resp.GrafanaInstanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMonitorClient().UnbindPrometheusManagedGrafana(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MonitorService) DescribeTkeTmpBasicConfigById(ctx context.Context, clusterId, clusterType, instanceId string) (respParams *monitor.DescribePrometheusConfigResponseParams, errRet error) {
	logId := getLogId(ctx)
	request := monitor.NewDescribePrometheusConfigRequest()
	request.InstanceId = &instanceId
	request.ClusterType = &clusterType
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMonitorClient().DescribePrometheusConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.RequestId == nil {
		return nil, fmt.Errorf("response is invalid,%s", response.ToJsonString())
	}

	respParams = response.Response
	return
}

func (me *MonitorService) GetConfigType(name string, respParams *monitor.DescribePrometheusConfigResponseParams) (configType string, config *monitor.PrometheusConfigItem, err error) {
	for _, v := range respParams.ServiceMonitors {
		if *v.Name == name {
			configType = "service_monitors"
			config = v
			return
		}
	}

	for _, v := range respParams.PodMonitors {
		if *v.Name == name {
			configType = "pod_monitors"
			config = v
			return
		}
	}

	for _, v := range respParams.RawJobs {
		if *v.Name == name {
			configType = "raw_jobs"
			config = v
			return
		}
	}
	err = fmt.Errorf("[ERROR] name [%v] configuration does not exist", name)
	return
}

type PrometheusConfig struct {
	Config *string
	Regex  []string
}

func (r *PrometheusConfig) UnmarshalToMap() (map[interface{}]interface{}, error) {
	var configMap map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(*r.Config), &configMap)
	if err != nil {
		log.Printf("[ERROR] yaml Unmarshal fail [%v]\n", err)
		return nil, err
	}
	return configMap, nil
}

func (r *PrometheusConfig) MarshalToYaml(config *map[interface{}]interface{}) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("[ERROR] yaml Marshal fail [%v]\n", err)
		return "", err
	}
	return string(data), nil
}

func (r *PrometheusConfig) SetRegex(configs []interface{}) (*[]interface{}, error) {
	setStatus := false
	regex := strings.Join(r.Regex, "|")
	for k, v := range configs {
		metricRelabelings := v.(map[interface{}]interface{})["metric_relabel_configs"]
		if metricRelabelings == nil {
			if v.(map[interface{}]interface{})["metricRelabelings"] != nil {
				metricRelabelings = v.(map[interface{}]interface{})["metricRelabelings"]
			} else {
				metricRelabelings = []interface{}{}
			}
		}
		metricRelabelingList := []interface{}{}
		for _, vv := range metricRelabelings.([]interface{}) {
			metricRelabeling := vv.(map[interface{}]interface{})
			sourceLabels := metricRelabeling["source_labels"]
			if sourceLabels == nil {
				sourceLabels = metricRelabeling["sourceLabels"]
			}
			if metricRelabeling["action"] == "keep" && sourceLabels.([]interface{})[0] == "__name__" {
				if regex == "" {
					metricRelabeling = nil
				} else {
					metricRelabeling["regex"] = regex
					setStatus = true
				}
			}
			if metricRelabeling["action"] == "drop" || metricRelabeling == nil {
			} else {
				metricRelabelingList = append(metricRelabelingList, metricRelabeling)
			}
		}

		if k == (len(configs)-1) && regex != "" && !setStatus {
			metricRelabeling := map[interface{}]interface{}{
				"source_labels": []string{"__name__"},
				"regex":         regex,
				"replacement":   "$1",
				"action":        "keep",
			}
			metricRelabelingList = append(metricRelabelingList, metricRelabeling)
		}
		if len(metricRelabelingList) > 0 {
			v.(map[interface{}]interface{})["metric_relabel_configs"] = metricRelabelingList
		}
	}
	return &configs, nil
}

func (me *MonitorService) DescribeMonitorTmpGrafanaConfigById(ctx context.Context, instanceId string) (tmpGrafanaConfig *monitor.DescribeGrafanaConfigResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := monitor.NewDescribeGrafanaConfigRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMonitorClient().DescribeGrafanaConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	tmpGrafanaConfig = response.Response
	return
}
