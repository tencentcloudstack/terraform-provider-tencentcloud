package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type ClbService struct {
	client *connectivity.TencentCloudClient
}

func (me *ClbService) DescribeLoadBalancerById(ctx context.Context, clbId string) (clbInstance *clb.LoadBalancer, errRet error) {
	logId := getLogId(ctx)
	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = []*string{&clbId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LoadBalancerSet) < 1 {
		errRet = fmt.Errorf("loadBalancer id is not found")
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
			request.LoadBalancerIds = []*string{stringToPointer(v.(string))}
		}

		if k == "network_type" {
			request.LoadBalancerType = stringToPointer(v.(string))
		}
		if k == "clb_name" {
			request.LoadBalancerName = stringToPointer(v.(string))
		}
		if k == "project_id" {
			projectId := int64(v.(int))
			request.ProjectId = &projectId
		}

	}

	offset := int64(0)
	pageSize := int64(100)
	clbs = make([]*clb.LoadBalancer, 0)
	for {
		request.Offset = &(offset)
		request.Limit = &(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
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
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
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
	request.LoadBalancerId = stringToPointer(clbId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
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
			request.ListenerIds = []*string{stringToPointer(listenerId)}
			request.LoadBalancerId = stringToPointer(clbId)
		}
		if k == "clb_id" {
			if clbId == "" {
				clbId = v.(string)
				request.LoadBalancerId = stringToPointer(clbId)
			}
		}
		if k == "protocol" {
			request.Protocol = stringToPointer(v.(string))
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
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
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
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteListener(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
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
		errRet = fmt.Errorf("Listener id and clb id can not be null")
		return
	}
	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerIds = []*string{&listenerId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
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
	request.LoadBalancerId = stringToPointer(clbId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
		return
	}
	clbListener := response.Response.Listeners[0]
	var ruleOutput clb.RuleOutput
	find_flag := false
	for _, rule := range clbListener.Rules {
		if *rule.Domain == domain && *rule.Url == url {
			ruleOutput = *rule
			find_flag = true
			break
		}
	}
	if !find_flag {
		errRet = fmt.Errorf("rule not found!")
		return
	} else {
		clbRule = &ruleOutput
	}
	return
}

func (me *ClbService) DeleteRuleById(ctx context.Context, clbId string, listenerId string, locationId string) error {
	logId := getLogId(ctx)
	request := clb.NewDeleteRuleRequest()
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)
	request.LocationIds = []*string{&locationId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteRule(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
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
	request.LoadBalancerId = stringToPointer(clbId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
		return
	}
	clbListener := response.Response.Listeners[0]
	protocol := clbListener.Protocol
	port := clbListener.Port

	aRequest := clb.NewDescribeTargetsRequest()
	aRequest.ListenerIds = []*string{&listenerId}
	aRequest.LoadBalancerId = stringToPointer(clbId)
	aRequest.Protocol = protocol
	aRequest.Port = port
	ratelimit.Check(request.GetAction())
	aResponse, aErr := me.client.UseClbClient().DescribeTargets(aRequest)

	if aErr != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, aRequest.GetAction(), aRequest.ToJsonString(), aErr.Error())
		errRet = aErr
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, aRequest.GetAction(), aRequest.ToJsonString(), aResponse.ToJsonString())

	if len(aResponse.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
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
		errRet = fmt.Errorf("Listener id and clb id can not be null")
		return
	}
	request.LoadBalancerId = stringToPointer(clbId)
	request.ListenerIds = []*string{&listenerId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
		return
	}
	clbListener := response.Response.Listeners[0]
	protocol := clbListener.Protocol
	port := clbListener.Port

	aRequest := clb.NewDescribeTargetsRequest()
	aRequest.ListenerIds = []*string{&listenerId}
	aRequest.LoadBalancerId = stringToPointer(clbId)
	aRequest.Protocol = protocol
	aRequest.Port = port
	ratelimit.Check(request.GetAction())
	aResponse, aErr := me.client.UseClbClient().DescribeTargets(aRequest)

	if aErr != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, aRequest.GetAction(), aRequest.ToJsonString(), aErr.Error())
		errRet = aErr
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, aRequest.GetAction(), aRequest.ToJsonString(), aResponse.ToJsonString())

	if len(aResponse.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
		return
	} else {
		clbAttachments = append(clbAttachments, aResponse.Response.Listeners[0])
	}

	return
}

func (me *ClbService) DeleteAttachmentById(ctx context.Context, clbId string, listenerId string, locationId string, targets []interface{}) error {
	logId := getLogId(ctx)
	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)
	for _, inst_ := range targets {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["port"], inst["weight"]))
	}
	if locationId != "" {
		request.LocationId = stringToPointer(locationId)
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeregisterTargets(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
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
		errRet = fmt.Errorf("rewriteInfo id is not exist!%s", rewriteId)
		return
	}
	sourceLocId := items[0]
	targetLocId := items[1]
	sourceListenerId := items[2]
	targetListenerId := items[3]
	clbId := items[4]
	result := make(map[string]string)
	request := clb.NewDescribeRewriteRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.SourceListenerIds = []*string{&sourceListenerId}
	request.SourceLocationIds = []*string{&sourceLocId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeRewrite(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RewriteSet) < 1 {
		errRet = fmt.Errorf("redirection id %s is not found", rewriteId)
		return
	}

	ruleOutput := response.Response.RewriteSet[0]
	if ruleOutput.RewriteTarget != nil {
		if *ruleOutput.RewriteTarget.TargetListenerId == targetListenerId && *ruleOutput.RewriteTarget.TargetLocationId == targetLocId {
			result["source_listener_rule_id"] = sourceLocId
			result["target_listener_rule_id"] = targetLocId
			result["source_listener_id"] = sourceListenerId
			result["target_listener_id"] = targetListenerId
			result["clb_id"] = clbId
			rewriteInfo = &result
		} else {
			errRet = fmt.Errorf("redirection id %s is not found", rewriteId)
		}
	} else {
		errRet = fmt.Errorf("redirection id %s is not found", rewriteId)
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
		if k == "source_listener_rule_id" {
			sourceLocId = v
		}
		if k == "target_listener_id" {
			targetListenerId = v
		}
		if k == "target_listener_rule_id" {
			targetLocId = v
		}
	}
	request := clb.NewDescribeRewriteRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.SourceListenerIds = []*string{&sourceListenerId}
	request.SourceLocationIds = []*string{&sourceLocId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DescribeRewrite(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
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
		result["source_listener_rule_id"] = sourceLocId
		result["target_listener_rule_id"] = *ruleOutput.RewriteTarget.TargetLocationId
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
		return fmt.Errorf("rewriteInfo id is not exist!!! %s", rewriteId)

	}
	sourceLocId := items[0]
	targetLocId := items[1]
	sourceListenerId := items[2]
	targetListenerId := items[3]
	clbId := items[4]

	request := clb.NewDeleteRewriteRequest()
	request.LoadBalancerId = stringToPointer(clbId)
	request.SourceListenerId = stringToPointer(sourceListenerId)
	request.TargetListenerId = stringToPointer(targetListenerId)
	var rewriteInfo clb.RewriteLocationMap
	rewriteInfo.SourceLocationId = stringToPointer(sourceLocId)
	rewriteInfo.TargetLocationId = stringToPointer(targetLocId)
	request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseClbClient().DeleteRewrite(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
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
	if v, ok := d.GetOk("health_check_switch"); ok {
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

	if v, ok := d.GetOk("health_check_http_code"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_code can only be set with protocol HTTP/HTTPS.")
			return
		} else {
			healthSetFlag = true
			vv := int64(v.(int))
			healthCheck.HttpCode = &vv
		}
	}

	if v, ok := d.GetOk("health_check_http_path"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_path can only be set with protocol HTTP/HTTPS")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckPath = stringToPointer(v.(string))
		}
	}

	if v, ok := d.GetOk("health_check_http_domain"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_domain can only be set with protocol HTTP/HTTPS")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckDomain = stringToPointer(v.(string))
		}
	}

	if v, ok := d.GetOk("health_check_http_method"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_method can only be set with protocol HTTP/HTTPS")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckMethod = stringToPointer(v.(string))
		}

	}

	if healthSetFlag {
		if !(((protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP || protocol == CLB_LISTENER_PROTOCOL_TCPSSL) && applyType == HEALTH_APPLY_TYPE_LISTENER) || ((protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) && applyType == HEALTH_APPLY_TYPE_RULE)) {
			healthSetFlag = false
			errRet = fmt.Errorf("health para can only be set with TCP/UDP/TCP_SSL listener or rule of HTTP/HTTPS listener")
			return
		}
		healthCheckPara = &healthCheck
	}
	return

}

func checkCertificateInputPara(ctx context.Context, d *schema.ResourceData) (certificateSetFlag bool, certPara *clb.CertificateInput, errRet error) {
	certificateSetFlag = false
	var certificateInput clb.CertificateInput
	certificateSSLMode := ""
	certificateId := ""
	certificateCaId := ""

	if v, ok := d.GetOk("certificate_ssl_mode"); ok {
		certificateSetFlag = true
		certificateSSLMode = v.(string)
		certificateInput.SSLMode = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateSetFlag = true
		certificateId = v.(string)
		certificateInput.CertId = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("certificate_ca_id"); ok {
		certificateSetFlag = true
		certificateCaId = v.(string)
		certificateInput.CertCaId = stringToPointer(v.(string))
	}

	if certificateSetFlag && certificateId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("certificatedId is null")
		return
	}

	if certificateSetFlag && certificateSSLMode == CERT_SSL_MODE_MUT && certificateCaId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("Certificate_ca_key is null and the ssl mode is 'MUTUAL' ")
		return
	}

	certPara = &certificateInput

	return
}
func waitForTaskFinish(requestId string, meta *clb.Client) (err error) {
	taskQueryRequest := clb.NewDescribeTaskStatusRequest()
	taskQueryRequest.TaskId = &requestId
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		taskResponse, e := meta.DescribeTaskStatus(taskQueryRequest)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *taskResponse.Response.Status == int64(CLB_TASK_EXPANDING) {
			return resource.RetryableError(fmt.Errorf("clb task status is %d, requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId))
		}
		return nil
	})
	return
}

func flattenClbTagsMapping(tags []*clb.TagInfo) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.TagKey] = *tag.TagValue
	}
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
