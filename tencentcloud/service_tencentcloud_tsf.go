package tencentcloud

import (
	"context"
	"fmt"
	"log"

	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TsfService struct {
	client *connectivity.TencentCloudClient
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

// func (me *TsfService) DeleteTsfApiRateLimitRuleById(ctx context.Context, apiId, ruleId string) (errRet error) {
// 	logId := getLogId(ctx)

// 	request := tsf.NewDeleteApiRateLimitRuleRequest()
// 	request.ApiId = &apiId

// 	defer func() {
// 		if errRet != nil {
// 			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
// 		}
// 	}()

// 	ratelimit.Check(request.GetAction())

// 	response, err := me.client.UseTsfClient().DeleteApiRateLimitRule(request)
// 	if err != nil {
// 		errRet = err
// 		return
// 	}
// 	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

// 	return
// }

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

func (me *TsfService) DeleteTsfApplicationReleaseConfigById(ctx context.Context, configId string) (errRet error) {
	logId := getLogId(ctx)

	request := tsf.NewRevocationConfigRequest()
	request.ConfigReleaseId = &configId

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
