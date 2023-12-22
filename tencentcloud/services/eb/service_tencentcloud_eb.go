package eb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewEbService(client *connectivity.TencentCloudClient) EbService {
	return EbService{client: client}
}

type EbService struct {
	client *connectivity.TencentCloudClient
}

func (me *EbService) DescribeEbSearchByFilter(ctx context.Context, param map[string]interface{}) (ebSearch []*string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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

func (me *EbService) DescribeEbEventConnectorById(ctx context.Context, connectionId string, eventBusId string) (eventConnector *eb.Connection, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := eb.NewListConnectionsRequest()
	request.EventBusId = &eventBusId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().ListConnections(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Connections) < 1 {
		return
	}

	for _, v := range response.Response.Connections {
		if *v.ConnectionId == connectionId {
			eventConnector = v
		}
	}

	return
}

func (me *EbService) DeleteEbEventConnectorById(ctx context.Context, connectionId string, eventBusId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := eb.NewDeleteConnectionRequest()
	request.ConnectionId = &connectionId
	request.EventBusId = &eventBusId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEbClient().DeleteConnection(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EbService) DescribeEbEventRulesByFilter(ctx context.Context, param map[string]interface{}) (eventRules []*eb.Rule, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = eb.NewListRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "EventBusId" {
			request.EventBusId = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
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
		response, err := me.client.UseEbClient().ListRules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Rules) < 1 {
			break
		}
		eventRules = append(eventRules, response.Response.Rules...)
		if len(response.Response.Rules) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *EbService) DescribeEbPlateformByFilter(ctx context.Context, param map[string]interface{}) (plateform []*eb.PlatformEventDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = eb.NewListPlatformEventNamesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProductType" {
			request.ProductType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEbClient().ListPlatformEventNames(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.EventNames) < 1 {
		return
	}

	plateform = response.Response.EventNames

	return
}

func (me *EbService) DescribeEbPlatformEventPatternsByFilter(ctx context.Context, param map[string]interface{}) (platformEventPatterns []*eb.PlatformEventSummary, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = eb.NewListPlatformEventPatternsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProductType" {
			request.ProductType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEbClient().ListPlatformEventPatterns(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.EventPatterns) < 1 {
		return
	}

	platformEventPatterns = response.Response.EventPatterns

	return
}

func (me *EbService) DescribeEbPlatformProductsByFilter(ctx context.Context, param map[string]interface{}) (platformProducts []*eb.PlatformProduct, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = eb.NewListPlatformProductsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEbClient().ListPlatformProducts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PlatformProducts) < 1 {
		return
	}

	platformProducts = response.Response.PlatformProducts

	return
}

func (me *EbService) DescribeEbPlateformEventTemplateByFilter(ctx context.Context, param map[string]interface{}) (plateformEventTemplate *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = eb.NewGetPlatformEventTemplateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "EventType" {
			request.EventType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEbClient().GetPlatformEventTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.EventTemplate == nil {
		return
	}

	plateformEventTemplate = response.Response.EventTemplate

	return
}
