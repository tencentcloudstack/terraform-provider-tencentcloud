package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ClbService struct {
	client *connectivity.TencentCloudClient
}

func RuleIdCheck(ruleId string) error {
	//check clb listener rule id first
	//old example file cause wrong usage of listener.id
	items := strings.Split(ruleId, FILED_SP)
	if len(items) > 1 {
		return fmt.Errorf("Unsupported references of rule_id since version 1.47.0, please check your tf content and use `tencentcloud_clb_listener_rule.xxx.rule_id` instead of `tencentcloud_clb_listener_rule.xxx.id`")
	}
	return nil
}

func ListenerIdCheck(listenerId string) error {
	//check clb listener listener id first
	//old example file cause wrong usage of listener.id
	items := strings.Split(listenerId, FILED_SP)
	if len(items) > 1 {
		return fmt.Errorf("Unsupported references of listener_id since version 1.47.0, please check your tf content and use `tencentcloud_clb_listener.xxx.listener_id` instead of `tencentcloud_clb_listener.xxx.id`")
	}
	return nil
}

func (me *ClbService) DescribeLoadBalancerById(ctx context.Context, clbId string) (clbInstance *clb.LoadBalancer, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = []*string{&clbId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LoadBalancerSet) < 1 {
		return
	}
	clbInstance = response.Response.LoadBalancerSet[0]
	return
}

func (me *ClbService) DescribeLoadBalancerByFilter(ctx context.Context, params map[string]interface{}) (clbs []*clb.LoadBalancer, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeLoadBalancersRequest()

	for k, v := range params {
		if k == "clb_id" {
			request.LoadBalancerIds = []*string{helper.String(v.(string))}
		}
		if k == "network_type" {
			request.LoadBalancerType = helper.String(v.(string))
		}
		if k == "clb_name" {
			request.LoadBalancerName = helper.String(v.(string))
		}
		if k == "project_id" {
			projectId := int64(v.(int))
			request.ProjectId = &projectId
		}
	}

	offset := int64(0)
	pageSize := int64(CLB_PAGE_LIMIT)
	clbs = make([]*clb.LoadBalancer, 0)
	for {
		request.Offset = &(offset)
		request.Limit = &(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
		if err != nil {
			errRet = errors.WithStack(err)
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LoadBalancerSet) < 1 {
			break
		}

		clbs = append(clbs, response.Response.LoadBalancerSet...)

		if int64(len(response.Response.LoadBalancerSet)) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClbService) DeleteLoadBalancerById(ctx context.Context, clbId string) error {
	logId := getLogId(ctx)
	request := clb.NewDeleteLoadBalancerRequest()
	request.LoadBalancerIds = []*string{&clbId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteLoadBalancer(request)
	if err != nil {
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return nil
			}
		}
		return errors.WithStack(err)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	requestId := *response.Response.RequestId
	retryErr := waitForTaskFinish(requestId, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (me *ClbService) DescribeListenerById(ctx context.Context, listenerId string, clbId string) (clbListener *clb.Listener, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	request.ListenerIds = []*string{&listenerId}
	request.LoadBalancerId = &clbId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		return
	}
	clbListener = response.Response.Listeners[0]
	return
}

func (me *ClbService) DescribeListenersByFilter(ctx context.Context, params map[string]interface{}) (listeners []*clb.Listener, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	clbId := ""
	for k, v := range params {
		if k == "listener_id" {
			listenerId := v.(string)
			request.ListenerIds = []*string{&listenerId}
			request.LoadBalancerId = &clbId
		}
		if k == "clb_id" {
			if clbId == "" {
				clbId = v.(string)
				request.LoadBalancerId = &clbId
			}
		}
		if k == "protocol" {
			request.Protocol = helper.String(v.(string))
		}
		if k == "port" {
			port := int64(v.(int))
			request.Port = &port
		}
	}

	listeners = make([]*clb.Listener, 0)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	listeners = append(listeners, response.Response.Listeners...)

	return
}

func (me *ClbService) DeleteListenerById(ctx context.Context, clbId string, listenerId string) error {
	logId := getLogId(ctx)
	request := clb.NewDeleteListenerRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = &clbId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteListener(request)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	requestId := *response.Response.RequestId
	retryErr := waitForTaskFinish(requestId, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (me *ClbService) DescribeRulesByFilter(ctx context.Context, params map[string]string) (rules []*clb.RuleOutput, errRet error) {
	logId := getLogId(ctx)
	//listener filter
	clbId := ""
	listenerId := ""
	//rule filter
	domain := ""
	url := ""
	locationId := ""
	scheduler := ""
	for k, v := range params {
		if k == "listener_id" {
			listenerId = v
		}
		if k == "clb_id" {
			clbId = v
		}
		if k == "domain" {
			domain = v
		}
		if k == "url" {
			url = v
		}
		if k == "rule_id" {
			locationId = v
		}
		if k == "scheduler" {
			scheduler = v
		}
	}
	//get listener first
	request := clb.NewDescribeListenersRequest()
	if listenerId == "" || clbId == "" {
		errRet = fmt.Errorf("[CHECK][CLB rule][Describe] check: Listener id and CLB id can not be null")
		return
	}
	request.LoadBalancerId = &clbId
	request.ListenerIds = []*string{&listenerId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	//listener not found, return empty
	if len(response.Response.Listeners) < 1 {
		return
	}
	clbListener := response.Response.Listeners[0]
	rules = make([]*clb.RuleOutput, 0)
	for _, rule := range clbListener.Rules {
		if domain != "" {
			if *rule.Domain != domain {
				continue
			}
		}
		if url != "" {
			if *rule.Url != url {
				continue
			}
		}
		if locationId != "" {
			if *rule.LocationId != locationId {
				continue
			}
		}
		if scheduler != "" {
			if *rule.Scheduler != scheduler {
				continue
			}
		}
		rules = append(rules, rule)
	}

	return
}

func (me *ClbService) DescribeRuleByPara(ctx context.Context, clbId string, listenerId string, domain string, url string) (clbRule *clb.RuleOutput, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	request.ListenerIds = []*string{&listenerId}
	request.LoadBalancerId = &clbId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	//listener not found, return empty
	if len(response.Response.Listeners) < 1 {
		return
	}
	clbListener := response.Response.Listeners[0]
	var ruleOutput clb.RuleOutput
	findFlag := false
	for _, rule := range clbListener.Rules {
		if *rule.Domain == domain && *rule.Url == url {
			ruleOutput = *rule
			findFlag = true
			break
		}
	}
	if !findFlag {
		errRet = fmt.Errorf("[CHECK][CLB rule][Describe] check: rule not found!")
		errRet = errors.WithStack(errRet)
		return
	} else {
		clbRule = &ruleOutput
	}
	return
}

func (me *ClbService) DeleteRuleById(ctx context.Context, clbId string, listenerId string, locationId string) error {
	logId := getLogId(ctx)
	request := clb.NewDeleteRuleRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = &clbId
	request.LocationIds = []*string{&locationId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteRule(request)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	requestId := *response.Response.RequestId
	retryErr := waitForTaskFinish(requestId, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (me *ClbService) DescribeAttachmentByPara(ctx context.Context, clbId string, listenerId string, locationId string) (clbAttachment *clb.ListenerBackend, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	request.ListenerIds = []*string{&listenerId}
	request.LoadBalancerId = &clbId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	//listener not found, return empty
	if len(response.Response.Listeners) < 1 {
		return
	}
	clbListener := response.Response.Listeners[0]
	protocol := clbListener.Protocol
	port := clbListener.Port

	aRequest := clb.NewDescribeTargetsRequest()
	aRequest.ListenerIds = []*string{&listenerId}
	aRequest.LoadBalancerId = &clbId
	aRequest.Protocol = protocol
	aRequest.Port = port
	ratelimit.Check(request.GetAction())
	aResponse, aErr := me.client.UseClbClient().DescribeTargets(aRequest)

	if aErr != nil {
		//in case that the lb is not exist, return empty
		if e, ok := aErr.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(aErr)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, aRequest.GetAction(), aRequest.ToJsonString(), aResponse.ToJsonString())

	if len(aResponse.Response.Listeners) < 1 {
		return
	} else {
		clbAttachment = aResponse.Response.Listeners[0]
	}

	return
}

func (me *ClbService) DescribeAttachmentsByFilter(ctx context.Context, params map[string]string) (clbAttachments []*clb.ListenerBackend, errRet error) {
	logId := getLogId(ctx)
	//listener filter
	clbId := ""
	listenerId := ""
	//rule filter

	for k, v := range params {
		if k == "listener_id" {
			listenerId = v
		}
		if k == "clb_id" {
			clbId = v
		}
	}
	//get listener first
	request := clb.NewDescribeListenersRequest()
	if listenerId == "" || clbId == "" {
		errRet = fmt.Errorf("[CHECK][CLB attachment][Describe] check: Listener id and clb id can not be null")
		errRet = errors.WithStack(errRet)
		return
	}
	request.LoadBalancerId = &clbId
	request.ListenerIds = []*string{&listenerId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		return
	}
	clbListener := response.Response.Listeners[0]
	protocol := clbListener.Protocol
	port := clbListener.Port

	aRequest := clb.NewDescribeTargetsRequest()
	aRequest.ListenerIds = []*string{&listenerId}
	aRequest.LoadBalancerId = &clbId
	aRequest.Protocol = protocol
	aRequest.Port = port
	ratelimit.Check(request.GetAction())
	aResponse, aErr := me.client.UseClbClient().DescribeTargets(aRequest)

	if aErr != nil {
		//in case that the lb is not exist, return empty
		if e, ok := aErr.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(aErr)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, aRequest.GetAction(), aRequest.ToJsonString(), aResponse.ToJsonString())

	if len(aResponse.Response.Listeners) < 1 {
		return
	} else {
		clbAttachments = append(clbAttachments, aResponse.Response.Listeners[0])
	}

	return
}

func (me *ClbService) DeleteAttachmentById(ctx context.Context, clbId string, listenerId string, locationId string, targets []interface{}) error {
	logId := getLogId(ctx)
	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = &listenerId
	request.LoadBalancerId = &clbId
	for _, inst_ := range targets {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}
	if locationId != "" {
		request.LocationId = &locationId
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeregisterTargets(request)
	if err != nil {
		ee, ok := err.(*sdkErrors.TencentCloudSDKError)
		if ok && ee.GetCode() == "InvalidParameter" {
			return nil
		}
		return errors.WithStack(err)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	requestId := *response.Response.RequestId
	retryErr := waitForTaskFinish(requestId, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (me *ClbService) DescribeRedirectionById(ctx context.Context, rewriteId string) (rewriteInfo *map[string]string, errRet error) {
	logId := getLogId(ctx)
	items := strings.Split(rewriteId, "#")
	if len(items) != 5 {
		errRet = fmt.Errorf("[CHECK][CLB redirection][Describe] check: redirection id %s is not with format loc-xxx#loc-xxx#lbl-xxx#lbl-xxx#lb-xxx", rewriteId)
		errRet = errors.WithStack(errRet)
		return
	}
	sourceLocId := items[0]
	targetLocId := items[1]
	sourceListenerId := items[2]
	targetListenerId := items[3]
	clbId := items[4]
	result := make(map[string]string)
	request := clb.NewDescribeRewriteRequest()
	request.LoadBalancerId = &clbId
	request.SourceListenerIds = []*string{&sourceListenerId}
	request.SourceLocationIds = []*string{&sourceLocId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeRewrite(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RewriteSet) < 1 {
		return
	}

	for _, v := range response.Response.RewriteSet {
		//sometimes the response returns all the rules under a certain url, so filter again in the code
		if v.RewriteTarget != nil {
			if *v.RewriteTarget.TargetListenerId == targetListenerId && *v.RewriteTarget.TargetLocationId == targetLocId {
				result["source_rule_id"] = sourceLocId
				result["target_rule_id"] = targetLocId
				result["source_listener_id"] = sourceListenerId
				result["target_listener_id"] = targetListenerId
				result["clb_id"] = clbId
				rewriteInfo = &result
				return
			}
		}
	}

	return
}

func (me *ClbService) DescribeAllAutoRedirections(ctx context.Context, rewriteId string) (rewriteInfos []*map[string]string, errRet error) {
	logId := getLogId(ctx)
	items := strings.Split(rewriteId, "#")
	if len(items) != 5 {
		errRet = fmt.Errorf("[CHECK][CLB redirection][Describe] check: redirection id %s is not with format loc-xxx#loc-xxx#lbl-xxx#lbl-xxx#lb-xxx", rewriteId)
		errRet = errors.WithStack(errRet)
		return
	}
	sourceLocationId := items[0]
	sourceListenerId := items[2]
	clbId := items[4]
	request := clb.NewDescribeRewriteRequest()
	request.LoadBalancerId = &clbId
	request.SourceListenerIds = []*string{&sourceListenerId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeRewrite(request)
	if err != nil {
		//in case that the lb is not exist, return empty
		if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if e.GetCode() == "InvalidParameter.LBIdNotFound" {
				return
			}
		}
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RewriteSet) < 1 {
		return
	}

	//get listener id and domain
	domain := ""
	for _, v := range response.Response.RewriteSet {
		if v.RewriteTarget != nil && sourceListenerId == *v.ListenerId && sourceLocationId == *v.LocationId {
			domain = *v.Domain
			break
		}
	}

	for _, v := range response.Response.RewriteSet {
		//auto rewrite will associate all the url under the domain
		if v.RewriteTarget != nil && sourceListenerId == *v.ListenerId && domain == *v.Domain {
			result := make(map[string]string)
			result["source_rule_id"] = *v.LocationId
			result["target_rule_id"] = *v.RewriteTarget.TargetLocationId
			result["source_listener_id"] = *v.ListenerId
			result["target_listener_id"] = *v.RewriteTarget.TargetListenerId
			result["clb_id"] = clbId
			rewriteInfos = append(rewriteInfos, &result)
		}
	}
	return
}

func (me *ClbService) DescribeRedirectionsByFilter(ctx context.Context, params map[string]string) (rewriteInfos []*map[string]string, errRet error) {
	logId := getLogId(ctx)
	clbId := ""
	sourceListenerId := ""
	sourceLocId := ""
	targetListenerId := ""
	targetLocId := ""
	for k, v := range params {
		if k == "source_listener_id" {
			sourceListenerId = v
		}
		if k == "clb_id" {
			clbId = v
		}
		if k == "source_rule_id" {
			sourceLocId = v
		}
		if k == "target_listener_id" {
			targetListenerId = v
		}
		if k == "target_rule_id" {
			targetLocId = v
		}
	}
	request := clb.NewDescribeRewriteRequest()
	request.LoadBalancerId = &clbId
	request.SourceListenerIds = []*string{&sourceListenerId}
	request.SourceLocationIds = []*string{&sourceLocId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeRewrite(request)
	if err != nil {
		errRet = err
		errRet = errors.WithStack(errRet)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RewriteSet) < 1 {
		return
	}
	ruleOutput := response.Response.RewriteSet[0]
	if ruleOutput.RewriteTarget != nil {
		if targetListenerId != "" && *ruleOutput.RewriteTarget.TargetListenerId != targetListenerId {
			return
		}
		if targetLocId != "" && *ruleOutput.RewriteTarget.TargetLocationId != targetLocId {
			return
		}
		result := make(map[string]string)
		result["source_rule_id"] = sourceLocId
		result["target_rule_id"] = *ruleOutput.RewriteTarget.TargetLocationId
		result["source_listener_id"] = sourceListenerId
		result["target_listener_id"] = *ruleOutput.RewriteTarget.TargetListenerId
		result["clb_id"] = clbId
		rewriteInfos = append(rewriteInfos, &result)
	}
	return
}

func (me *ClbService) DeleteRedirectionById(ctx context.Context, rewriteId string) error {
	logId := getLogId(ctx)
	items := strings.Split(rewriteId, "#")
	if len(items) != 5 {
		errRet := fmt.Errorf("[CHECK][CLB redirection][Describe] check: redirection id %s is not with format loc-xxx#loc-xxx#lbl-xxx#lbl-xxx#lb-xxx", rewriteId)
		errRet = errors.WithStack(errRet)
		return errRet
	}
	sourceLocId := items[0]
	targetLocId := items[1]
	sourceListenerId := items[2]
	targetListenerId := items[3]
	clbId := items[4]

	request := clb.NewDeleteRewriteRequest()
	request.LoadBalancerId = &clbId
	request.SourceListenerId = &sourceListenerId
	request.TargetListenerId = &targetListenerId
	var rewriteInfo clb.RewriteLocationMap
	rewriteInfo.SourceLocationId = &sourceLocId
	rewriteInfo.TargetLocationId = &targetLocId
	request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteRewrite(request)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	requestId := *response.Response.RequestId
	retryErr := waitForTaskFinish(requestId, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func checkHealthCheckPara(ctx context.Context, d *schema.ResourceData, protocol string, applyType string) (healthSetFlag bool, healthCheckPara *clb.HealthCheck, errRet error) {
	var healthCheck clb.HealthCheck
	healthSetFlag = false
	healthCheckPara = &healthCheck
	if v, ok := d.GetOkExists("health_check_switch"); ok {
		healthSetFlag = true
		vv := v.(bool)
		vvv := int64(0)
		if vv {
			vvv = 1
		}
		healthCheck.HealthSwitch = &vvv
	}
	if v, ok := d.GetOk("health_check_time_out"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.TimeOut = &vv
	}
	if v, ok := d.GetOk("health_check_interval_time"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.IntervalTime = &vv
	}
	if v, ok := d.GetOk("health_check_health_num"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.HealthNum = &vv
	}
	if v, ok := d.GetOk("health_check_unhealth_num"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.UnHealthNum = &vv
	}
	if v, ok := d.GetOk("health_check_port"); ok {
		healthSetFlag = true
		healthCheck.CheckPort = helper.Int64(int64(v.(int)))
	}
	var checkType string
	if v, ok := d.GetOk("health_check_type"); ok {
		healthSetFlag = true
		checkType = v.(string)
		healthCheck.CheckType = &checkType
	}
	if v, ok := d.GetOk("health_check_http_code"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS ||
			(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_HTTP)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_code can only be set with protocol HTTP/HTTPS or HTTP of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.HttpCode = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("health_check_http_path"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS ||
			(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_HTTP)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_path can only be set with protocol HTTP/HTTPS or HTTP of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.HttpCheckPath = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_http_domain"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS ||
			(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_HTTP)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_domain can only be set with protocol HTTP/HTTPS or HTTP of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.HttpCheckDomain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_http_method"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS ||
			(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_HTTP)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_method can only be set with protocol HTTP/HTTPS or HTTP of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.HttpCheckMethod = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_http_version"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_HTTP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_version can only be set with protocol HTTP of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.HttpVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_context_type"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_CUSTOM) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_context_type can only be set with protocol CUSTOM of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.ContextType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_send_context"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_CUSTOM) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_send_context can only be set with protocol CUSTOM of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.SendContext = helper.String(v.(string))
	}
	if v, ok := d.GetOk("health_check_recv_context"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP && checkType == HEALTH_CHECK_TYPE_CUSTOM) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_recv_context can only be set with protocol CUSTOM of TCP")
			errRet = errors.WithStack(errRet)
			return
		}
		healthSetFlag = true
		healthCheck.RecvContext = helper.String(v.(string))
	}

	if healthSetFlag {
		if !(((protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP ||
			protocol == CLB_LISTENER_PROTOCOL_TCPSSL) && applyType == HEALTH_APPLY_TYPE_LISTENER) ||
			((protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) &&
				applyType == HEALTH_APPLY_TYPE_RULE)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health para can only be set with TCP/UDP/TCP_SSL listener or rule of HTTP/HTTPS listener")
			errRet = errors.WithStack(errRet)
			return
		}
		if protocol == CLB_LISTENER_PROTOCOL_TCP {
			if checkType == HEALTH_CHECK_TYPE_HTTP && healthCheck.HttpCheckDomain == nil {
				healthCheck.HttpCheckDomain = helper.String("")
			}
			if healthCheck.CheckPort == nil {
				healthCheck.CheckPort = helper.Int64(-1)
			}
			if healthCheck.HttpCheckPath == nil {
				healthCheck.HttpCheckPath = helper.String("")
			}
		}
		healthCheckPara = &healthCheck
	}
	return

}

func checkCertificateInputPara(ctx context.Context, d *schema.ResourceData, meta interface{}) (certificateSetFlag bool, certPara *clb.CertificateInput, errRet error) {
	certificateSetFlag = false
	var certificateInput clb.CertificateInput
	certificateSSLMode := ""
	certificateId := ""
	certificateCaId := ""

	if v, ok := d.GetOk("certificate_ssl_mode"); ok {
		certificateSetFlag = true
		certificateSSLMode = v.(string)
		certificateInput.SSLMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateSetFlag = true
		certificateId = v.(string)
		certificateInput.CertId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("certificate_ca_id"); ok {
		certificateSetFlag = true
		certificateCaId = v.(string)
		certificateInput.CertCaId = helper.String(v.(string))
	}

	if certificateSetFlag && certificateId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("certificated Id is null")
		errRet = errors.WithStack(errRet)
		return
	}

	if certificateSetFlag && certificateSSLMode == CERT_SSL_MODE_MUT && certificateCaId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("certificate_ca_key is null and the ssl mode is 'MUTUAL' ")
		errRet = errors.WithStack(errRet)
		return
	}

	certPara = &certificateInput

	//check type valid
	sslService := SSLService{client: meta.(*TencentCloudClient).apiV3Conn}

	if certificateInput.CertCaId != nil {
		check, err := sslService.checkCertificateType(ctx, *certificateInput.CertCaId, SSL_CERT_TYPE_CA)
		if err != nil {
			certificateSetFlag = false
			errRet = fmt.Errorf("certificated %s check error %s", *certificateInput.CertCaId, err)
			errRet = errors.WithStack(errRet)
			return
		}
		if !check {
			certificateSetFlag = false
			errRet = fmt.Errorf("certificated %s check error cert type is not `%s`", *certificateInput.CertCaId, SSL_CERT_TYPE_CA)
			return
		}
	}
	if certificateInput.CertId != nil {
		check, err := sslService.checkCertificateType(ctx, *certificateInput.CertId, SSL_CERT_TYPE_SERVER)
		if err != nil {
			certificateSetFlag = false
			errRet = fmt.Errorf("certificated %s check error %s", *certificateInput.CertId, err)
			errRet = errors.WithStack(errRet)
			return
		}
		if !check {
			certificateSetFlag = false
			errRet = fmt.Errorf("certificated %s check error cert type is not `%s`", *certificateInput.CertId, SSL_CERT_TYPE_SERVER)
			errRet = errors.WithStack(errRet)
			return
		}
	}
	return
}

func waitForTaskFinish(requestId string, meta *clb.Client) (err error) {
	taskQueryRequest := clb.NewDescribeTaskStatusRequest()
	taskQueryRequest.TaskId = &requestId
	err = resource.Retry(4*readRetryTimeout, func() *resource.RetryError {
		taskResponse, e := meta.DescribeTaskStatus(taskQueryRequest)
		if e != nil {
			return resource.NonRetryableError(errors.WithStack(e))
		}
		if *taskResponse.Response.Status == int64(CLB_TASK_EXPANDING) {
			return resource.RetryableError(errors.WithStack(fmt.Errorf("CLB task status is %d(expanding), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		} else if *taskResponse.Response.Status == int64(CLB_TASK_FAIL) {
			return resource.NonRetryableError(errors.WithStack(fmt.Errorf("CLB task status is %d(failed), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		}
		return nil
	})
	return
}

func flattenBackendList(list []*clb.Backend) (mapping []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		target := map[string]interface{}{
			"instance_id": v.InstanceId,
			"port":        v.Port,
			"weight":      v.Weight,
		}
		result = append(result, target)
	}
	return result
}

func clbNewTarget(instanceId, port, weight interface{}) *clb.Target {
	id := instanceId.(string)
	p := int64(port.(int))
	bk := clb.Target{
		InstanceId: &id,
		Port:       &p,
	}
	if w, ok := weight.(int); ok {
		weight64 := int64(w)
		bk.Weight = &weight64
	}
	return &bk
}

func (me *ClbService) CreateTargetGroup(ctx context.Context, targetGroupName string, vpcId string, port uint64,
	targetGroupInstances []*clb.TargetGroupInstance) (targetGroupId string, err error) {
	var response *clb.CreateTargetGroupResponse

	request := clb.NewCreateTargetGroupRequest()
	request.TargetGroupName = &targetGroupName
	request.TargetGroupInstances = targetGroupInstances
	request.Port = &port
	if vpcId != "" {
		request.VpcId = &vpcId
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseClbClient().CreateTargetGroup(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return
	}
	if response.Response.TargetGroupId == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty ", request.GetAction())
		return
	}
	targetGroupId = *response.Response.TargetGroupId
	return
}

func (me *ClbService) ModifyTargetGroup(ctx context.Context, targetGroupId string, targetGroupName string) (err error) {
	request := clb.NewModifyTargetGroupAttributeRequest()
	request.TargetGroupId = &targetGroupId
	request.TargetGroupName = &targetGroupName

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseClbClient().ModifyTargetGroupAttribute(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (me *ClbService) RegisterTargetInstances(ctx context.Context, targetGroupId, bindIp string, port, weight uint64) (err error) {
	request := clb.NewRegisterTargetGroupInstancesRequest()
	request.TargetGroupId = &targetGroupId
	request.TargetGroupInstances = []*clb.TargetGroupInstance{
		{
			BindIP: &bindIp,
			Port:   &port,
			Weight: &weight,
		},
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseClbClient().RegisterTargetGroupInstances(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *ClbService) DeregisterTargetInstances(ctx context.Context, targetGroupId, bindIp string, port uint64) (err error) {
	request := clb.NewDeregisterTargetGroupInstancesRequest()
	request.TargetGroupId = &targetGroupId
	request.TargetGroupInstances = []*clb.TargetGroupInstance{
		{
			BindIP: &bindIp,
			Port:   &port,
		},
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseClbClient().DeregisterTargetGroupInstances(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *ClbService) DeleteTarget(ctx context.Context, targetGroupId string) error {
	request := clb.NewDeleteTargetGroupsRequest()
	request.TargetGroupIds = []*string{&targetGroupId}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseClbClient().DeleteTargetGroups(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *ClbService) DescribeTargetGroups(ctx context.Context, targetGroupId string, filters map[string]string) (targetGroupInfos []*clb.TargetGroupInfo, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeTargetGroupsRequest()
	if targetGroupId != "" {
		request.TargetGroupIds = []*string{&targetGroupId}
	}
	request.Filters = make([]*clb.Filter, 0, len(filters))
	for k, v := range filters {
		filter := clb.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize = uint64(CLB_PAGE_LIMIT)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClbClient().DescribeTargetGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.TargetGroupSet) < 1 {
			break
		}
		targetGroupInfos = append(targetGroupInfos, response.Response.TargetGroupSet...)
		if len(response.Response.TargetGroupSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClbService) DescribeTargetGroupInstances(ctx context.Context, filters map[string]string) (targetGroupInstances []*clb.TargetGroupBackend, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeTargetGroupInstancesRequest()
	request.Filters = make([]*clb.Filter, 0, len(filters))
	for k, v := range filters {
		filter := clb.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize = uint64(CLB_PAGE_LIMIT)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClbClient().DescribeTargetGroupInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.TargetGroupInstanceSet) < 1 {
			break
		}
		targetGroupInstances = append(targetGroupInstances, response.Response.TargetGroupInstanceSet...)
		if len(response.Response.TargetGroupInstanceSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClbService) AssociateTargetGroups(ctx context.Context, listenerId, clbId, targetGroupId, locationId string) (errRet error) {
	request := clb.NewAssociateTargetGroupsRequest()
	association := clb.TargetGroupAssociation{
		LoadBalancerId: &clbId,
		ListenerId:     &listenerId,
		TargetGroupId:  &targetGroupId,
	}
	if locationId != "" {
		association.LocationId = &locationId
	}
	request.Associations = append(request.Associations, &association)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseClbClient().AssociateTargetGroups(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	errRet = err
	return
}

func (me *ClbService) DescribeAssociateTargetGroups(ctx context.Context, ids []string) (has bool, err error) {
	var (
		logId       = getLogId(ctx)
		targetInfos []*clb.TargetGroupInfo
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetInfos, err = me.DescribeTargetGroups(ctx, ids[0], nil)
		if err != nil {
			return retryError(err, InternalError)
		}
		if targetInfos == nil || len(targetInfos[0].AssociatedRule) == 0 {
			return resource.RetryableError(fmt.Errorf("response is nil"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s DescribeTargetGroups failed, reason:%s ", logId, err.Error())
		return
	}

	for _, info := range targetInfos {
		for _, rule := range info.AssociatedRule {
			var originLocationId string
			originClbId := *rule.LoadBalancerId
			originListenerId := *rule.ListenerId
			if rule.LocationId != nil {
				originLocationId = *rule.LocationId
			}

			if *rule.Protocol == CLB_LISTENER_PROTOCOL_TCP || *rule.Protocol == CLB_LISTENER_PROTOCOL_UDP || *rule.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
				if originListenerId == ids[1] && originClbId == ids[2] {
					return true, nil
				}
			} else if originListenerId == ids[1] && originClbId == ids[2] && originLocationId == ids[3] {
				return true, nil
			}
		}
	}

	return false, nil
}

func (me *ClbService) DisassociateTargetGroups(ctx context.Context, targetGroupId, listenerId, clbId, locationId string) (errRet error) {
	var ruleId *string

	if locationId != "" {
		ruleId = &locationId
	}

	request := clb.NewDisassociateTargetGroupsRequest()
	request.Associations = []*clb.TargetGroupAssociation{
		{
			LoadBalancerId: &clbId,
			ListenerId:     &listenerId,
			TargetGroupId:  &targetGroupId,
			LocationId:     ruleId,
		},
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseClbClient().DisassociateTargetGroups(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	errRet = err
	return
}

func (me *ClbService) ModifyTargetGroupInstancesWeight(ctx context.Context, targetGroupId, bindIp string, port, weight uint64) (errRet error) {
	var instance = clb.TargetGroupInstance{
		BindIP: &bindIp,
		Port:   &port,
		Weight: &weight,
	}
	request := clb.NewModifyTargetGroupInstancesWeightRequest()
	request.TargetGroupId = &targetGroupId
	request.TargetGroupInstances = []*clb.TargetGroupInstance{&instance}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseClbClient().ModifyTargetGroupInstancesWeight(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
