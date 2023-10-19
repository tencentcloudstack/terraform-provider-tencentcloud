package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ElasticsearchService struct {
	client *connectivity.TencentCloudClient
}

func (me *ElasticsearchService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *es.InstanceInfo, errRet error) {
	logId := getLogId(ctx)
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

	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
		ctx := contextNil

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
	logId := getLogId(ctx)

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
		ctx := contextNil

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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

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
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseEsClient().UpdateLogstashInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update elasticsearch logstash failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"1"}, 3*readRetryTimeout, time.Second, me.ElasticsearchLogstashStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}
