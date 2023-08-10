package tencentcloud

import (
	"context"
	"log"

	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type EbService struct {
	client *connectivity.TencentCloudClient
}

func (me *EbService) DescribeEbSearchByFilter(ctx context.Context, param map[string]interface{}) (ebSearch []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = eb.NewDescribeLogTagValueRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "EventBusId" {
			request.EventBusId = v.(*string)
		}
		if k == "GroupField" {
			request.GroupField = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.([]*eb.LogFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Page = &offset
		request.Limit = &limit
		response, err := me.client.UseEbClient().DescribeLogTagValue(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Results) < 1 {
			break
		}
		ebSearch = append(ebSearch, response.Response.Results...)
		if len(response.Response.Results) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *EbService) DescribeEbSearchLogByFilter(ctx context.Context, param map[string]interface{}) (ebSearch []*eb.SearchLogResult, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = eb.NewSearchLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "EventBusId" {
			request.EventBusId = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.([]*eb.LogFilter)
		}
		if k == "OrderFields" {
			request.OrderFields = v.([]*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Page = &offset
		request.Limit = &limit
		response, err := me.client.UseEbClient().SearchLog(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Results) < 1 {
			break
		}
		ebSearch = append(ebSearch, response.Response.Results...)
		if len(response.Response.Results) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *EbService) DescribeEbEventBusById(ctx context.Context, eventBusId string) (event *eb.GetEventBusResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := eb.NewGetEventBusRequest()
	request.EventBusId = &eventBusId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().GetEventBus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	event = response.Response
	return
}

func (me *EbService) DeleteEbEventBusById(ctx context.Context, eventBusId string) (errRet error) {
	logId := getLogId(ctx)

	request := eb.NewDeleteEventBusRequest()
	request.EventBusId = &eventBusId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().DeleteEventBus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EbService) DescribeEbEventTargetById(ctx context.Context, eventBusId string, ruleId string, targetId string) (eventTarget *eb.Target, errRet error) {
	logId := getLogId(ctx)

	request := eb.NewListTargetsRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().ListTargets(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Targets) < 1 {
		return
	}

	for _, v := range response.Response.Targets {
		if *v.TargetId == targetId {
			eventTarget = v
			return
		}
	}

	return
}

func (me *EbService) DeleteEbEventTargetById(ctx context.Context, eventBusId string, ruleId string, targetId string) (errRet error) {
	logId := getLogId(ctx)

	request := eb.NewDeleteTargetRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TargetId = &targetId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().DeleteTarget(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EbService) DescribeEbEventRuleById(ctx context.Context, eventBusId string, ruleId string) (rule *eb.GetRuleResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := eb.NewGetRuleRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().GetRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	rule = response.Response
	return
}

func (me *EbService) DeleteEbEventRuleById(ctx context.Context, eventBusId string, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := eb.NewDeleteRuleRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().DeleteRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EbService) DescribeEbEventTransformById(ctx context.Context, eventBusId string, ruleId string, transformationId string) (ebTransform *eb.Transformation, errRet error) {
	logId := getLogId(ctx)

	request := eb.NewGetTransformationRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TransformationId = &transformationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().GetTransformation(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Transformations) < 1 {
		return
	}

	ebTransform = response.Response.Transformations[0]
	return
}

func (me *EbService) DeleteEbEventTransformById(ctx context.Context, eventBusId string, ruleId string, transformationId string) (errRet error) {
	logId := getLogId(ctx)

	request := eb.NewDeleteTransformationRequest()
	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TransformationId = &transformationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().DeleteTransformation(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EbService) DescribeEbBusByFilter(ctx context.Context, param map[string]interface{}) (bus []*eb.EventBus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = eb.NewListEventBusesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*eb.Filter)
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
		response, err := me.client.UseEbClient().ListEventBuses(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EventBuses) < 1 {
			break
		}
		bus = append(bus, response.Response.EventBuses...)
		if len(response.Response.EventBuses) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
