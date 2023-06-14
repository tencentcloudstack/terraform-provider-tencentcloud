package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TsfService struct {
	client *connectivity.TencentCloudClient
}

func (me *TsfService) DescribeTsfClusterById(ctx context.Context, clusterId string) (cluster *tsf.ClusterV2, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeClustersRequest()
	request.ClusterIdList = []*string{&clusterId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*tsf.ClusterV2, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeClusters(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		instances = append(instances, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	cluster = instances[0]
	return
}

func (me *TsfService) DeleteTsfClusterById(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteClusterRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteCluster(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) TsfClusterStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeTsfClusterById(ctx, clusterId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.ClusterStatus), nil
	}
}

func (me *TsfService) DescribeTsfApplicationConfigById(ctx context.Context, configId, configName string) (applicationConfig *tsf.Config, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeConfigsRequest()
	request.ConfigId = &configId
	if configId != "" {
		request.ConfigId = &configId
	}
	if configName != "" {
		request.ConfigName = &configName
	}
	if configId == "" && configName == "" {
		errRet = fmt.Errorf("`configId` and `configName` cannot both be empty")
		return
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	if len(response.Response.Result.Content) < 1 {
		return
	}

	applicationConfig = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationConfigById(ctx context.Context, configId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteConfigRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfMicroserviceById(ctx context.Context, namespaceId, microserviceId, microserviceName string) (microservice *tsf.Microservice, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeMicroservicesRequest()
	request.NamespaceId = &namespaceId
	if microserviceId != "" {
		request.MicroserviceIdList = []*string{&microserviceId}
	}
	if microserviceName != "" {
		request.MicroserviceNameList = []*string{&microserviceName}
	}
	if microserviceId == "" && microserviceName == "" {
		errRet = fmt.Errorf("`microserviceId` and `microserviceName` cannot both be empty")
		return
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeMicroservices(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.Content) < 1 {
		return
	}

	microservice = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfMicroserviceById(ctx context.Context, microserviceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteMicroserviceRequest()
	request.MicroserviceId = &microserviceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteMicroservice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApiGroupById(ctx context.Context, groupId string) (apiGroup *tsf.ApiGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeApiGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeApiGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	apiGroup = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfApiGroupById(ctx context.Context, groupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteApiGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteApiGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApiRateLimitRuleById(ctx context.Context, apiId, ruleId string) (apiRateLimitRule *tsf.ApiRateLimitRule, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeApiRateLimitRulesRequest()
	request.ApiId = &apiId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeApiRateLimitRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result) < 1 {
		return
	}

	for _, v := range response.Response.Result {
		if *v.RuleId == ruleId {
			apiRateLimitRule = v
			return
		}
	}

	return
}

func (me *TsfService) DeleteTsfApiRateLimitRuleById(ctx context.Context, apiId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteApiRateLimitRuleRequest()
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteApiRateLimitRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfConfigTemplateById(ctx context.Context, templateId string) (configTemplate *tsf.ConfigTemplate, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeConfigTemplateRequest()
	request.ConfigTemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeConfigTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	configTemplate = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfConfigTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteConfigTemplateRequest()
	request.ConfigTemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteConfigTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfLaneById(ctx context.Context, laneId string) (lane *tsf.LaneInfo, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeLanesRequest()
	request.LaneIdList = []*string{&laneId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeLanes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
		return
	}

	lane = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfLaneById(ctx context.Context, laneId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteLaneRequest()
	request.LaneId = &laneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteLane(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfLaneRuleById(ctx context.Context, ruleId string) (laneRule *tsf.LaneRule, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeLaneRulesRequest()
	request.RuleId = &ruleId
	request.Limit = helper.IntInt64(10)
	request.Offset = helper.IntInt64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeLaneRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
		return
	}

	laneRule = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfLaneRuleById(ctx context.Context, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteLaneRuleRequest()
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteLaneRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfNamespaceById(ctx context.Context, namespaceId string) (namespace *tsf.Namespace, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeSimpleNamespacesRequest()
	request.NamespaceId = &namespaceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeSimpleNamespaces(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
		return
	}

	namespace = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfNamespaceById(ctx context.Context, namespaceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteNamespaceRequest()
	request.NamespaceId = &namespaceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteNamespace(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *TsfService) DescribeTsfPathRewriteById(ctx context.Context, pathRewriteId string) (pathRewrite *tsf.PathRewrite, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribePathRewriteRequest()
	request.PathRewriteId = &pathRewriteId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribePathRewrite(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	pathRewrite = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfPathRewriteById(ctx context.Context, pathRewriteId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeletePathRewritesRequest()
	request.PathRewriteIds = []*string{&pathRewriteId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeletePathRewrites(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfRepositoryById(ctx context.Context, repositoryId string) (repository *tsf.RepositoryInfo, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeRepositoryRequest()
	request.RepositoryId = &repositoryId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeRepository(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	repository = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfRepositoryById(ctx context.Context, repositoryId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteRepositoryRequest()
	request.RepositoryId = &repositoryId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteRepository(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfTaskById(ctx context.Context, taskId string) (task *tsf.TaskRecord, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeTaskDetailRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeTaskDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	task = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfTaskById(ctx context.Context, taskId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteTaskRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfUnitRuleById(ctx context.Context, id string) (unitRule *tsf.UnitRule, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeUnitRuleRequest()
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeUnitRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	unitRule = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfUnitRuleById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteUnitRuleRequest()
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteUnitRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *TsfService) DescribeTsfContainGroupById(ctx context.Context, groupId string) (containGroup *tsf.ContainerGroupDetail, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeContainerGroupDetailRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeContainerGroupDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	containGroup = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfContainGroupById(ctx context.Context, groupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteContainerGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteContainerGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationReleaseConfigById(ctx context.Context, configId string, groupId string) (applicationReleaseConfig *tsf.ConfigRelease, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeConfigReleasesRequest()
	request.ConfigId = &configId
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeConfigReleases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
		return
	}

	applicationReleaseConfig = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationReleaseConfigById(ctx context.Context, configReleaseId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewRevocationConfigRequest()
	request.ConfigReleaseId = &configReleaseId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().RevocationConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfGroupById(ctx context.Context, groupId string) (group *tsf.VmGroup, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	group = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfGroupById(ctx context.Context, groupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationById(ctx context.Context, applicationId string) (application *tsf.ApplicationForPage, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeApplicationRequest()
	request.ApplicationId = &applicationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeApplication(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	application = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfApplicationById(ctx context.Context, applicationId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteApplicationRequest()
	request.ApplicationId = &applicationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteApplication(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationFileConfigReleaseById(ctx context.Context, configId string, groupId string) (applicationFileConfigRelease *tsf.FileConfigRelease, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeFileConfigReleasesRequest()
	request.ConfigId = &configId
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeFileConfigReleases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.Content) < 1 {
		return
	}

	applicationFileConfigRelease = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationFileConfigReleaseById(ctx context.Context, configId string, groupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewRevokeFileConfigRequest()

	applicationfileConfig, err := me.DescribeTsfApplicationFileConfigReleaseById(ctx, configId, groupId)
	if err != nil {
		log.Printf("[CRITAL]%s Describe tsf applicationFileConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	request.ConfigReleaseId = applicationfileConfig.ConfigReleaseId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().RevokeFileConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationPublicConfigById(ctx context.Context, configId string) (applicationPublicConfig *tsf.Config, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribePublicConfigRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribePublicConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	applicationPublicConfig = response.Response.Result
	return
}

func (me *TsfService) DeleteTsfApplicationPublicConfigById(ctx context.Context, configId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeletePublicConfigRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeletePublicConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationPublicConfigReleaseById(ctx context.Context, configId, namespaceId string) (applicationPublicConfigRelease *tsf.ConfigRelease, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribePublicConfigReleasesRequest()
	request.ConfigId = &configId
	request.NamespaceId = &namespaceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribePublicConfigReleases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
		return
	}

	applicationPublicConfigRelease = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationPublicConfigReleaseById(ctx context.Context, configId, namespaceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewRevocationPublicConfigRequest()

	applicationPublicConfig, err := me.DescribeTsfApplicationPublicConfigReleaseById(ctx, configId, namespaceId)
	if err != nil {
		log.Printf("[CRITAL]%s Describe tsf applicationfileConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	request.ConfigReleaseId = applicationPublicConfig.ConfigReleaseId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().RevocationPublicConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfInstancesAttachmentById(ctx context.Context, clusterId string, instanceId string) (instance *tsf.Instance, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeClusterInstancesRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeClusterInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		for _, v := range response.Response.Result.Content {
			if *v.InstanceId == instanceId {
				instance = v
				return
			}
		}
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TsfService) DeleteTsfInstancesAttachmentById(ctx context.Context, clusterId string, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewRemoveInstancesRequest()
	request.ClusterId = &clusterId
	request.InstanceIdList = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().RemoveInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfApplicationByFilter(ctx context.Context, param map[string]interface{}) (application *tsf.TsfPageApplication, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeApplicationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ApplicationType" {
			request.ApplicationType = v.(*string)
		}
		if k == "MicroserviceType" {
			request.MicroserviceType = v.(*string)
		}
		if k == "ApplicationResourceTypeList" {
			request.ApplicationResourceTypeList = v.([]*string)
		}
		if k == "ApplicationIdList" {
			request.ApplicationIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset  int64 = 0
		limit   int64 = 20
		total   int64
		content = make([]*tsf.ApplicationForPage, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeApplications(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		content = append(content, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	application = &tsf.TsfPageApplication{
		TotalCount: &total,
		Content:    content,
	}

	return
}

func (me *TsfService) DescribeTsfApplicationConfigByFilter(ctx context.Context, param map[string]interface{}) (applicationConfig *tsf.TsfPageConfig, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeConfigsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
		if k == "ConfigId" {
			request.ConfigId = v.(*string)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
		if k == "ConfigName" {
			request.ConfigName = v.(*string)
		}
		if k == "ConfigVersion" {
			request.ConfigVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		config = make([]*tsf.Config, 0)
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		config = append(config, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	applicationConfig = &tsf.TsfPageConfig{
		TotalCount: &total,
		Content:    config,
	}

	return
}

func (me *TsfService) DescribeTsfApplicationFileConfigByFilter(ctx context.Context, param map[string]interface{}) (applicationFileConfig *tsf.TsfPageFileConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeFileConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ConfigId" {
			request.ConfigId = v.(*string)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
		if k == "ConfigName" {
			request.ConfigName = v.(*string)
		}
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
		if k == "ConfigVersion" {
			request.ConfigVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		config = make([]*tsf.FileConfig, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeFileConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		config = append(config, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	applicationFileConfig = &tsf.TsfPageFileConfig{
		TotalCount: &total,
		Content:    config,
	}

	return
}

func (me *TsfService) DescribeTsfApplicationPublicConfigByFilter(ctx context.Context, param map[string]interface{}) (applicationPublicConfig *tsf.TsfPageConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribePublicConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ConfigId" {
			request.ConfigId = v.(*string)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
		if k == "ConfigName" {
			request.ConfigName = v.(*string)
		}
		if k == "ConfigVersion" {
			request.ConfigVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		config = make([]*tsf.Config, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribePublicConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		config = append(config, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	applicationPublicConfig = &tsf.TsfPageConfig{
		TotalCount: &total,
		Content:    config,
	}

	return
}

func (me *TsfService) DescribeTsfClusterByFilter(ctx context.Context, param map[string]interface{}) (cluster *tsf.TsfPageCluster, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeSimpleClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterIdList" {
			request.ClusterIdList = v.([]*string)
		}
		if k == "ClusterType" {
			request.ClusterType = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "DisableProgramAuthCheck" {
			request.DisableProgramAuthCheck = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset   int64 = 0
		limit    int64 = 20
		total    int64
		clusters = make([]*tsf.Cluster, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeSimpleClusters(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		clusters = append(clusters, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	cluster = &tsf.TsfPageCluster{
		TotalCount: &total,
		Content:    clusters,
	}

	return
}

func (me *TsfService) DescribeTsfMicroserviceByFilter(ctx context.Context, param map[string]interface{}) (microservice *tsf.TsfPageMicroservice, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeMicroservicesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "NamespaceId" {
			request.NamespaceId = v.(*string)
		}
		if k == "Status" {
			request.Status = v.([]*string)
		}
		if k == "MicroserviceIdList" {
			request.MicroserviceIdList = v.([]*string)
		}
		if k == "MicroserviceNameList" {
			request.MicroserviceNameList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		micro  = make([]*tsf.Microservice, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeMicroservices(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		micro = append(micro, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	microservice = &tsf.TsfPageMicroservice{
		TotalCount: &total,
		Content:    micro,
	}

	return
}

func (me *TsfService) DescribeTsfUnitRulesByFilter(ctx context.Context, param map[string]interface{}) (unitRule *tsf.TsfPageUnitRuleV2, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeUnitRulesV2Request()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayInstanceId" {
			request.GatewayInstanceId = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		rules  = make([]*tsf.UnitRule, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeUnitRulesV2(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		rules = append(rules, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	unitRule = &tsf.TsfPageUnitRuleV2{
		TotalCount: &total,
		Content:    rules,
	}

	return
}

func (me *TsfService) DescribeTsfConfigSummaryByFilter(ctx context.Context, param map[string]interface{}) (configSummary *tsf.TsfPageConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeConfigSummaryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "ConfigTagList" {
			request.ConfigTagList = v.([]*string)
		}
		if k == "DisableProgramAuthCheck" {
			request.DisableProgramAuthCheck = v.(*bool)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		config = make([]*tsf.Config, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeConfigSummary(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		config = append(config, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	configSummary = &tsf.TsfPageConfig{
		TotalCount: &total,
		Content:    config,
	}

	return
}

func (me *TsfService) DescribeTsfDeliveryConfigByGroupIdByFilter(ctx context.Context, param map[string]interface{}) (deliveryConfigByGroupID *tsf.SimpleKafkaDeliveryConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeDeliveryConfigByGroupIdRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeDeliveryConfigByGroupId(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	deliveryConfigByGroupID = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfDeliveryConfigsByFilter(ctx context.Context, param map[string]interface{}) (deliveryConfigs *tsf.DeliveryConfigBindGroups, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeDeliveryConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset    int64 = 0
		limit     int64 = 20
		total     int64
		bindGroup = make([]*tsf.DeliveryConfigBindGroup, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeDeliveryConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}

		total = *response.Response.Result.TotalCount
		bindGroup = append(bindGroup, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	deliveryConfigs = &tsf.DeliveryConfigBindGroups{
		TotalCount: &total,
		Content:    bindGroup,
	}

	return
}

func (me *TsfService) DescribeTsfPublicConfigSummaryByFilter(ctx context.Context, param map[string]interface{}) (publicConfigSummary *tsf.TsfPageConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribePublicConfigSummaryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "ConfigTagList" {
			request.ConfigTagList = v.([]*string)
		}
		if k == "DisableProgramAuthCheck" {
			request.DisableProgramAuthCheck = v.(*bool)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset        int64 = 0
		limit         int64 = 20
		total         int64
		configSummary = make([]*tsf.Config, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribePublicConfigSummary(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		configSummary = append(configSummary, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	publicConfigSummary = &tsf.TsfPageConfig{
		TotalCount: &total,
		Content:    configSummary,
	}

	return
}

func (me *TsfService) DescribeTsfApiGroupByFilter(ctx context.Context, param map[string]interface{}) (apiGroupInfo *tsf.TsfPageApiGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeApiGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "GroupType" {
			request.GroupType = v.(*string)
		}
		if k == "AuthType" {
			request.AuthType = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "GatewayInstanceId" {
			request.GatewayInstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
		info   = make([]*tsf.ApiGroupInfo, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeApiGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		info = append(info, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	apiGroupInfo = &tsf.TsfPageApiGroupInfo{
		TotalCount: &total,
		Content:    info,
	}

	return
}

func (me *TsfService) DescribeTsfApplicationAttributeByFilter(ctx context.Context, param map[string]interface{}) (applicationAttribute *tsf.ApplicationAttribute, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeApplicationAttributeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeApplicationAttribute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	applicationAttribute = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfBusinessLogConfigsByFilter(ctx context.Context, param map[string]interface{}) (businessLogConfigs *tsf.TsfPageBusinessLogConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeBusinessLogConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "DisableProgramAuthCheck" {
			request.DisableProgramAuthCheck = v.(*bool)
		}
		if k == "ConfigIdList" {
			request.ConfigIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
		total  int64
		info   = make([]*tsf.BusinessLogConfig, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeBusinessLogConfigs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		info = append(info, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	businessLogConfigs = &tsf.TsfPageBusinessLogConfig{
		TotalCount: &total,
		Content:    info,
	}

	return
}

func (me *TsfService) DescribeTsfApiDetailByFilter(ctx context.Context, param map[string]interface{}) (apiDetail *tsf.ApiDetailResponse, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeApiDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "MicroserviceId" {
			request.MicroserviceId = v.(*string)
		}
		if k == "Path" {
			request.Path = v.(*string)
		}
		if k == "Method" {
			request.Method = v.(*string)
		}
		if k == "PkgVersion" {
			request.PkgVersion = v.(*string)
		}
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeApiDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	apiDetail = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfMicroserviceApiVersionByFilter(ctx context.Context, param map[string]interface{}) (microserviceApiVersion []*tsf.ApiVersionArray, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeApiVersionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "MicroserviceId" {
			request.MicroserviceId = v.(*string)
		}
		if k == "Path" {
			request.Path = v.(*string)
		}
		if k == "Method" {
			request.Method = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeApiVersions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Result == nil {
		return
	}

	microserviceApiVersion = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfBindApiGroupById(ctx context.Context, groupId string, gatewayDeployGroupId string) (bindApiGroup *tsf.GatewayDeployGroup, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeGroupBindedGatewaysRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeGroupBindedGateways(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
			break
		}

		for _, v := range response.Response.Result.Content {
			if *v.DeployGroupId == gatewayDeployGroupId {
				bindApiGroup = v
				return
			}
		}

		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TsfService) DeleteTsfBindApiGroupById(ctx context.Context, groupId string, gatewayDeployGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewUnbindApiGroupRequest()
	request.GroupGatewayList = []*tsf.GatewayGroupIds{
		{
			GatewayDeployGroupId: &gatewayDeployGroupId,
			GroupId:              &groupId,
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().UnbindApiGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfRepositoryByFilter(ctx context.Context, param map[string]interface{}) (repositoryList *tsf.RepositoryList, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeRepositoriesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "RepositoryType" {
			request.RepositoryType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset     uint64 = 0
		limit      uint64 = 20
		total      int64
		repository = make([]*tsf.RepositoryInfo, 0)
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeRepositories(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		repository = append(repository, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	repositoryList = &tsf.RepositoryList{
		TotalCount: &total,
		Content:    repository,
	}

	return
}

func (me *TsfService) DescribeTsfApplicationFileConfigById(ctx context.Context, configId string) (applicationFileConfig *tsf.FileConfig, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeFileConfigsRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeFileConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Result.Content) < 1 {
		return
	}

	applicationFileConfig = response.Response.Result.Content[0]
	return
}

func (me *TsfService) DeleteTsfApplicationFileConfigById(ctx context.Context, configId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteFileConfigRequest()
	request.ConfigId = &configId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteFileConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfEnableUnitRuleById(ctx context.Context, id string) (enableUnitRuleAttachment *tsf.UnitRule, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeUnitRuleRequest()
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeUnitRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	enableUnitRuleAttachment = response.Response.Result
	return
}

func (me *TsfService) DescribeTsfDescribePodInstancesByFilter(ctx context.Context, param map[string]interface{}) (describePodInstances *tsf.GroupPodResult, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribePodInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
		if k == "PodNameList" {
			request.PodNameList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset    int64 = 0
		limit     int64 = 20
		instances []*tsf.GroupPod
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribePodInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		instances = append(instances, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	describePodInstances = &tsf.GroupPodResult{
		TotalCount: helper.IntInt64(len(instances)),
		Content:    instances,
	}

	return
}

func (me *TsfService) DescribeTsfGatewayAllGroupApisByFilter(ctx context.Context, param map[string]interface{}) (gatewayAllGroupApis *tsf.GatewayVo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeGatewayAllGroupApisRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayDeployGroupId" {
			request.GatewayDeployGroupId = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeGatewayAllGroupApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	gatewayAllGroupApis = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfGroupGatewaysByFilter(ctx context.Context, param map[string]interface{}) (groupGateways *tsf.TsfPageApiGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeGroupGatewaysRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayDeployGroupId" {
			request.GatewayDeployGroupId = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset   int64 = 0
		limit    int64 = 20
		gateways []*tsf.ApiGroupInfo
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeGroupGateways(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		gateways = append(gateways, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	groupGateways = &tsf.TsfPageApiGroupInfo{
		TotalCount: helper.IntInt64(len(gateways)),
		Content:    gateways,
	}

	return
}

func (me *TsfService) DescribeTsfGroupInstancesByFilter(ctx context.Context, param map[string]interface{}) (groupInstances *tsf.TsfPageInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeGroupInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset    int64 = 0
		limit     int64 = 20
		instances []*tsf.Instance
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeGroupInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		instances = append(instances, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	groupInstances = &tsf.TsfPageInstance{
		TotalCount: helper.IntInt64(len(instances)),
		Content:    instances,
	}

	return
}

func (me *TsfService) DescribeTsfUsableUnitNamespacesByFilter(ctx context.Context, param map[string]interface{}) (usableUnitNamespaces *tsf.TsfPageUnitNamespace, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeUsableUnitNamespacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset     int64 = 0
		limit      int64 = 20
		namespaces []*tsf.UnitNamespace
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeUsableUnitNamespaces(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		namespaces = append(namespaces, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	usableUnitNamespaces = &tsf.TsfPageUnitNamespace{
		TotalCount: helper.IntInt64(len(namespaces)),
		Content:    namespaces,
	}

	return
}

func (me *TsfService) DescribeTsfGroupConfigReleaseByFilter(ctx context.Context, param map[string]interface{}) (groupConfigRelease *tsf.GroupRelease, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeGroupReleaseRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GroupId" {
			request.GroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTsfClient().DescribeGroupRelease(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	groupConfigRelease = response.Response.Result
	return
}

func (me *TsfService) DescribeTsfDeployVmGroupById(ctx context.Context, groupId string) (deployVmGroup *tsf.VmGroup, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	deployVmGroup = response.Response.Result
	return
}

func (me *TsfService) DescribeTsfReleaseApiGroupById(ctx context.Context, groupId string) (releaseApiGroup *tsf.ApiGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeApiGroupRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeApiGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	releaseApiGroup = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfStartContainerGroupById(ctx context.Context, groupId string) (startContainerGroup *tsf.ContainerGroupOther, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeContainerGroupAttributeRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeContainerGroupAttribute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	startContainerGroup = response.Response.Result

	return
}

func (me *TsfService) DescribeTsfStartGroupById(ctx context.Context, groupId string) (startGroup *tsf.VmGroupOther, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeGroupAttributeRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeGroupAttribute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	startGroup = response.Response.Result
	return
}

func (me *TsfService) DescribeTsfUnitNamespaceById(ctx context.Context, gatewayInstanceId, namespaceId string) (unitNamespace *tsf.UnitNamespace, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeUnitNamespacesRequest()
	request.GatewayInstanceId = &gatewayInstanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeUnitNamespaces(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	for _, v := range response.Response.Result.Content {
		if *v.NamespaceId == namespaceId {
			unitNamespace = v
			return
		}
	}

	return
}

func (me *TsfService) DeleteTsfUnitNamespaceById(ctx context.Context, gatewayInstanceId, unitNamespace string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDeleteUnitNamespacesRequest()
	request.GatewayInstanceId = &gatewayInstanceId
	request.UnitNamespaceList = []*string{&unitNamespace}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DeleteUnitNamespaces(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TsfService) DescribeTsfDeployContainerGroupById(ctx context.Context, groupId string) (deployContainerGroup *tsf.ContainerGroupDeploy, errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewDescribeContainerGroupDeployInfoRequest()
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTsfClient().DescribeContainerGroupDeployInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil {
		return
	}

	deployContainerGroup = response.Response.Result
	return
}

func (me *TsfService) DescribeTsfDescriptionContainerGroupByFilter(ctx context.Context, param map[string]interface{}) (descriptionContainerGroup *tsf.ContainGroupResult, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeContainerGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "NamespaceId" {
			request.NamespaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		group  []*tsf.ContainGroup
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeContainerGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		group = append(group, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	descriptionContainerGroup = &tsf.ContainGroupResult{
		TotalCount: helper.IntInt64(len(group)),
		Content:    group,
	}

	return
}

func (me *TsfService) DescribeTsfGroupsByFilter(ctx context.Context, param map[string]interface{}) (groups *tsf.TsfPageVmGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
		if k == "ApplicationId" {
			request.ApplicationId = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}
		if k == "NamespaceId" {
			request.NamespaceId = v.(*string)
		}
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "GroupResourceTypeList" {
			request.GroupResourceTypeList = v.([]*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "GroupIdList" {
			request.GroupIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		group  []*tsf.VmGroupSimple
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		group = append(group, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	groups = &tsf.TsfPageVmGroup{
		TotalCount: helper.IntInt64(len(group)),
		Content:    group,
	}

	return
}

func (me *TsfService) DescribeTsfMsApiListByFilter(ctx context.Context, param map[string]interface{}) (msApiList *tsf.TsfApiListResponse, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tsf.NewDescribeMsApiListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "MicroserviceId" {
			request.MicroserviceId = v.(*string)
		}
		if k == "SearchWord" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		msApi  []*tsf.MsApiArray
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTsfClient().DescribeMsApiList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.Content) < 1 {
			break
		}
		msApi = append(msApi, response.Response.Result.Content...)
		if len(response.Response.Result.Content) < int(limit) {
			break
		}

		offset += limit
	}

	msApiList = &tsf.TsfApiListResponse{
		TotalCount: helper.IntInt64(len(msApi)),
		Content:    msApi,
	}

	return
}
