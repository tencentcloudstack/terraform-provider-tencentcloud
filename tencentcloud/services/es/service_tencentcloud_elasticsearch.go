package es

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewElasticsearchService(client *connectivity.TencentCloudClient) ElasticsearchService {
	return ElasticsearchService{client: client}
}

type ElasticsearchService struct {
	client *connectivity.TencentCloudClient
}

func (me *ElasticsearchService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *es.InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := es.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEsClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if len(response.Response.InstanceList) < 1 {
		return
	}
	instance = response.Response.InstanceList[0]
	return
}

func (me *ElasticsearchService) DescribeInstancesByFilter(ctx context.Context, instanceId, instanceName string,
	tags map[string]string) (instances []*es.InstanceInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := es.NewDescribeInstancesRequest()
	if instanceId != "" {
		request.InstanceIds = []*string{&instanceId}
	}
	if instanceName != "" {
		request.InstanceNames = []*string{&instanceName}
	}
	for k, v := range tags {
		tag := es.TagInfo{
			TagKey:   helper.String(k),
			TagValue: helper.String(v),
		}
		request.TagList = append(request.TagList, &tag)
	}

	offset := 0
	pageSize := 100
	instances = make([]*es.InstanceInfo, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseEsClient().DescribeInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ElasticsearchService) DeleteInstance(ctx context.Context, instanceId string) error {
	logId := tccommon.GetLogId(ctx)
	request := es.NewDeleteInstanceRequest()
	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().DeleteInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

// UpdateInstance FIXME: use *Request instead of these suck params
func (me *ElasticsearchService) UpdateInstance(ctx context.Context, instanceId, instanceName, password string, basicSecurityType int64, nodeList []*es.NodeInfo, nodeTypeInfo *es.WebNodeTypeInfo, esAcl *es.EsAcl) error {
	logId := tccommon.GetLogId(ctx)
	request := es.NewUpdateInstanceRequest()
	request.InstanceId = &instanceId
	if instanceName != "" {
		request.InstanceName = &instanceName
	}
	if password != "" {
		request.Password = &password
	}
	if basicSecurityType > 0 {
		request.BasicSecurityType = &basicSecurityType
	}
	if nodeList != nil {
		request.NodeInfoList = nodeList
	}
	if nodeTypeInfo != nil {
		request.WebNodeTypeInfo = nodeTypeInfo
	}
	if esAcl != nil {
		request.EsAcl = esAcl
	}
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpdateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) UpdateInstanceVersion(ctx context.Context, instanceId, version string) error {
	logId := tccommon.GetLogId(ctx)
	request := es.NewUpgradeInstanceRequest()
	request.InstanceId = &instanceId
	request.EsVersion = &version

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpgradeInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) UpdateInstanceLicense(ctx context.Context, instanceId, licenseType string) error {
	logId := tccommon.GetLogId(ctx)
	request := es.NewUpgradeLicenseRequest()
	request.InstanceId = &instanceId
	request.LicenseType = &licenseType

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpgradeLicense(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) DescribeElasticsearchIndexByName(ctx context.Context, instanceId, indexType, indexName string) (index *es.IndexMetaField, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDescribeIndexMetaRequest()
	request.InstanceId = &instanceId
	request.IndexName = &indexName
	request.IndexType = &indexType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeIndexMeta(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil {
		index = response.Response.IndexMetaField
	}
	return
}

func (me *ElasticsearchService) DeleteElasticsearchIndexByName(ctx context.Context, instanceId, indexType, indexName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDeleteIndexRequest()
	request.InstanceId = &instanceId
	request.IndexType = &indexType
	request.IndexName = &indexName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DeleteIndex(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ElasticsearchService) DescribeElasticsearchLogstashById(ctx context.Context, instanceId string) (logstash *es.LogstashInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDescribeLogstashInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeLogstashInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceList) < 1 {
		return
	}

	logstash = response.Response.InstanceList[0]
	return
}

func (me *ElasticsearchService) DeleteElasticsearchLogstashById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDeleteLogstashInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DeleteLogstashInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ElasticsearchService) ElasticsearchLogstashStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeElasticsearchLogstashById(ctx, instanceId)
		log.Printf("object: %v, err: %v", object, err)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return &es.LogstashInstanceInfo{}, "-99", nil
		}
		return object, helper.Int64ToStr(*object.Status), nil
	}
}

func (me *ElasticsearchService) DescribeElasticsearchLogstashPipelineById(ctx context.Context, instanceId, pipelineId string) (logstashPipeline *es.LogstashPipelineInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDescribeLogstashPipelinesRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeLogstashPipelines(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LogstashPipelineList) < 1 {
		return
	}

	for _, item := range response.Response.LogstashPipelineList {
		if *item.PipelineId == pipelineId {
			logstashPipeline = item
			break
		}
	}
	return
}

func (me *ElasticsearchService) ElasticsearchLogstashPipelineStateRefreshFunc(instanceId, pipelineId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeElasticsearchLogstashPipelineById(ctx, instanceId, pipelineId)
		log.Printf("object: %v, err: %v", object, err)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return &es.LogstashPipelineInfo{}, "-99", nil
		}
		return object, helper.Int64ToStr(*object.Status), nil
	}
}

func (me *ElasticsearchService) DeleteElasticsearchLogstashPipelineById(ctx context.Context, instanceId, pipelineId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDeleteLogstashPipelinesRequest()
	request.InstanceId = &instanceId
	request.PipelineIds = []*string{&pipelineId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DeleteLogstashPipelines(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ElasticsearchService) UpdateLogstashInstance(ctx context.Context, instanceId string, params map[string]interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := es.NewUpdateLogstashInstanceRequest()
	request.InstanceId = &instanceId

	for k, v := range params {
		if k == "instance_name" {
			request.InstanceName = helper.String(v.(string))
		}
		if k == "node_num" {
			request.NodeNum = helper.IntUint64(v.(int))
		}
		if k == "node_type" {
			request.NodeType = helper.String(v.(string))
		}
		if k == "disk_size" {
			request.DiskSize = helper.IntUint64(v.(int))
		}
		if k == "operation_duration" {
			operationDurationUpdated := v.(es.OperationDurationUpdated)
			request.OperationDuration = &operationDurationUpdated
		}
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseEsClient().UpdateLogstashInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update elasticsearch logstash failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 3*tccommon.ReadRetryTimeout, time.Second, me.ElasticsearchLogstashStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}

func (me *ElasticsearchService) UpdateJdk(ctx context.Context, instanceId string, params map[string]interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := es.NewUpdateJdkRequest()
	request.InstanceId = helper.String(instanceId)
	if v, ok := params["Jdk"]; ok {
		request.Jdk = helper.String(v.(string))
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().UpdateJdk(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return err
}

func (me *ElasticsearchService) DescribeElasticsearchInstanceLogsByFilter(ctx context.Context, param map[string]interface{}) (elasticsearchInstanceLogs []*es.InstanceLog, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeInstanceLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "LogType" {
			request.LogType = v.(*uint64)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEsClient().DescribeInstanceLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceLogList) < 1 {
			break
		}
		elasticsearchInstanceLogs = append(elasticsearchInstanceLogs, response.Response.InstanceLogList...)
		if len(response.Response.InstanceLogList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ElasticsearchService) DescribeElasticsearchInstanceOperationsByFilter(ctx context.Context, param map[string]interface{}) (instanceOperations []*elasticsearch.Operation, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeInstanceOperationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEsClient().DescribeInstanceOperations(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Operations) < 1 {
			break
		}
		instanceOperations = append(instanceOperations, response.Response.Operations...)
		if len(response.Response.Operations) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ElasticsearchService) DescribeElasticsearchLogstashInstanceLogsByFilter(ctx context.Context, param map[string]interface{}) (logstashInstanceLogs []*elasticsearch.InstanceLog, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeLogstashInstanceLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "LogType" {
			request.LogType = v.(*uint64)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEsClient().DescribeLogstashInstanceLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceLogList) < 1 {
			break
		}
		logstashInstanceLogs = append(logstashInstanceLogs, response.Response.InstanceLogList...)
		if len(response.Response.InstanceLogList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ElasticsearchService) DescribeElasticsearchLogstashInstanceOperationsByFilter(ctx context.Context, param map[string]interface{}) (logstashInstanceOperations []*elasticsearch.Operation, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeLogstashInstanceOperationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEsClient().DescribeLogstashInstanceOperations(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Operations) < 1 {
			break
		}
		logstashInstanceOperations = append(logstashInstanceOperations, response.Response.Operations...)
		if len(response.Response.Operations) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ElasticsearchService) DescribeElasticsearchViewsByFilter(ctx context.Context, param map[string]interface{}) (clusterView *elasticsearch.ClusterView, nodesViews []*elasticsearch.NodeView, kibanasViews []*elasticsearch.KibanaView, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeViewsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeViews(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	clusterView = response.Response.ClusterView
	nodesViews = response.Response.NodesView
	kibanasViews = response.Response.KibanasView

	return
}

func (me *ElasticsearchService) ElasticsearchInstanceRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeInstanceById(ctx, instanceId)
		log.Printf("object: %v, err: %v", object, err)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return &es.InstanceInfo{}, "-99", nil
		}
		return object, helper.Int64ToStr(*object.Status), nil
	}
}

func (me *ElasticsearchService) DescribeElasticsearchDictionariesById(ctx context.Context, instanceId string) (Dictionaries *elasticsearch.DiagnoseResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewDescribeDiagnoseRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeDiagnose(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiagnoseResults) < 1 {
		return
	}

	Dictionaries = response.Response.DiagnoseResults[0]
	return
}

func (me *ElasticsearchService) UpdateDiagnoseSettings(ctx context.Context, instanceId string, params map[string]interface{}) error {
	logId := tccommon.GetLogId(ctx)
	request := es.NewUpdateDiagnoseSettingsRequest()
	request.InstanceId = helper.String(instanceId)

	for k, v := range params {
		if k == "Status" {
			request.Status = helper.IntInt64(v.(int))
		}
		if k == "CronTime" {
			request.CronTime = helper.String(v.(string))
		}
	}
	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseEsClient().UpdateDiagnoseSettings(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update es Diagnose failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func (me *ElasticsearchService) GetDiagnoseSettingsById(ctx context.Context, instanceId string) (diagnoseSettings *es.GetDiagnoseSettingsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := es.NewGetDiagnoseSettingsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().GetDiagnoseSettings(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	diagnoseSettings = response.Response
	return
}

func (me *ElasticsearchService) DescribeElasticsearchDiagnoseByFilter(ctx context.Context, param map[string]interface{}) (diagnose []*elasticsearch.DiagnoseResult, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeDiagnoseRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = helper.String(v.(string))
		}
		if k == "Date" {
			request.Date = helper.String(v.(string))
		}
		if k == "Limit" {
			request.Limit = helper.IntInt64(v.(int))
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEsClient().DescribeDiagnose(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DiagnoseResults) < 1 {
		return
	}
	diagnose = append(diagnose, response.Response.DiagnoseResults...)

	return
}

func (me *ElasticsearchService) DescribeElasticsearchInstancePluginListByFilter(ctx context.Context, param map[string]interface{}) (InstancePluginList []*elasticsearch.DescribeInstancePluginInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = es.NewDescribeInstancePluginListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
		if k == "PluginType" {
			request.PluginType = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEsClient().DescribeInstancePluginList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.PluginList) < 1 {
			break
		}
		InstancePluginList = append(InstancePluginList, response.Response.PluginList...)
		if len(response.Response.PluginList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ElasticsearchService) DescribeElasticsearchDescribeIndexListByFilter(ctx context.Context, param map[string]interface{}) (DescribeIndexList []*elasticsearch.IndexMetaField, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = elasticsearch.NewDescribeIndexListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "IndexType" {
			request.IndexType = v.(*string)
		}
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "IndexName" {
			request.IndexName = v.(*string)
		}
		if k == "Username" {
			request.Username = v.(*string)
		}
		if k == "Password" {
			request.Password = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "IndexStatusList" {
			request.IndexStatusList = v.([]*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
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
		response, err := me.client.UseEsClient().DescribeIndexList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.IndexMetaFields) < 1 {
			break
		}
		DescribeIndexList = append(DescribeIndexList, response.Response.IndexMetaFields...)
		if len(response.Response.IndexMetaFields) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
