package tencentcloud

import (
	"context"
	"fmt"
	"log"

	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CssService struct {
	client *connectivity.TencentCloudClient
}

func (me *CssService) DescribeCssWatermark(ctx context.Context, watermarkId string) (watermark *css.WatermarkInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveWatermarkRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.WatermarkId = helper.StrToUint64Point(watermarkId)

	response, err := me.client.UseCssClient().DescribeLiveWatermark(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.Watermark == nil {
		return
	}
	watermark = response.Response.Watermark
	return
}

func (me *CssService) DescribeCssWatermarks(ctx context.Context) (marks []*css.WatermarkInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveWatermarksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseCssClient().DescribeLiveWatermarks(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || *response.Response.TotalNum < 1 {
		return
	}
	marks = response.Response.WatermarkList
	return
}

func (me *CssService) DeleteCssWatermarkById(ctx context.Context, watermarkId *int64) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveWatermarkRequest()

	request.WatermarkId = watermarkId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DeleteLiveWatermark(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssPullStreamTask(ctx context.Context, taskId string) (tasks []*css.PullStreamTaskInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLivePullStreamTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if taskId != "" {
		request.TaskId = &taskId
	}

	response, err := me.client.UseCssClient().DescribeLivePullStreamTasks(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || *response.Response.TotalNum < 1 {
		return
	}
	tasks = response.Response.TaskInfos
	return
}

func (me *CssService) DeleteCssPullStreamTaskById(ctx context.Context, taskId, operator *string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLivePullStreamTaskRequest()

	request.TaskId = taskId
	request.Operator = operator

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DeleteLivePullStreamTask(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssLiveTranscodeTemplate(ctx context.Context, templateId *int64) (liveTranscodeTemplate *css.TemplateInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveTranscodeTemplateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.TemplateId = templateId

	response, err := me.client.UseCssClient().DescribeLiveTranscodeTemplate(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.Template == nil {
		return
	}
	liveTranscodeTemplate = response.Response.Template
	return
}

func (me *CssService) DescribeCssLiveTranscodeTemplates(ctx context.Context) (temps []*css.TemplateInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveTranscodeTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseCssClient().DescribeLiveTranscodeTemplates(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || len(response.Response.Templates) < 1 {
		return
	}
	temps = response.Response.Templates
	return
}

func (me *CssService) DeleteCssLiveTranscodeTemplateById(ctx context.Context, templateId *int64) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveTranscodeTemplateRequest()

	request.TemplateId = templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DeleteLiveTranscodeTemplate(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssLiveTranscodeRuleAttachment(ctx context.Context, domainName, templateId *string) (rules []*css.RuleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveTranscodeRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if domainName != nil {
		request.DomainNames = []*string{domainName}
	}

	if templateId != nil {
		request.TemplateIds = []*int64{helper.StrToInt64Point(*templateId)}
	}

	response, err := me.client.UseCssClient().DescribeLiveTranscodeRules(request)
	if err != nil {
		log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Rules) < 1 {
		return
	}
	rules = response.Response.Rules
	return
}

func (me *CssService) DeleteCssLiveTranscodeRuleAttachmentById(ctx context.Context, domainName, appName, streamName, templateId *string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveTranscodeRuleRequest()

	request.DomainName = domainName
	request.AppName = appName
	request.StreamName = streamName
	request.TemplateId = helper.Int64(helper.StrToInt64(*templateId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DeleteLiveTranscodeRule(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssDomainById(ctx context.Context, name string) (domain *css.DomainInfo, errRet error) {
	logId := getLogId(ctx)

	request := css.NewDescribeLiveDomainRequest()
	request.DomainName = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	result := response.Response.DomainInfo
	if result != nil {
		domain = result
		return
	}

	return
}

func (me *CssService) DeleteCssDomainById(ctx context.Context, name *string, dtype *uint64) (errRet error) {
	logId := getLogId(ctx)

	if name == nil || dtype == nil {
		return fmt.Errorf("DeleteCssDomainById: the required parameters name and type are nil!")
	}
	request := css.NewDeleteLiveDomainRequest()
	request.DomainName = name
	request.DomainType = dtype

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssDomainsByFilter(ctx context.Context, param map[string]interface{}) (domains []*css.DomainInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveDomainsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DomainStatus" {
			request.DomainStatus = v.(*uint64)
		}

		if k == "DomainType" {
			request.DomainType = v.(*uint64)
		}

		if k == "IsDelayLive" {
			request.IsDelayLive = v.(*uint64)
		}

		if k == "DomainPrefix" {
			request.DomainPrefix = v.(*string)
		}

		if k == "PlayType" {
			request.PlayType = v.(*uint64)
		}

	}
	ratelimit.Check(request.GetAction())

	var currNumber uint64 = 1
	var pageSize uint64 = 20

	for {
		request.PageNum = &currNumber
		request.PageSize = &pageSize

		response, err := me.client.UseCssClient().DescribeLiveDomains(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DomainList) < 1 {
			break
		}

		domains = append(domains, response.Response.DomainList...)
		if len(response.Response.DomainList) < int(pageSize) {
			break
		}
		currNumber++
	}
	return
}

func (me *CssService) DescribeCssPlayDomainCertAttachmentById(ctx context.Context, domainName string, cloudCertId string) (playDomainCertAttachment *css.LiveDomainCertBindings, errRet error) {
	logId := getLogId(ctx)

	request := css.NewDescribeLiveDomainCertBindingsRequest()
	request.DomainName = &domainName

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
	instances := make([]*css.LiveDomainCertBindings, 0)
	for {
		request.Offset = &offset
		request.Length = &limit
		response, err := me.client.UseCssClient().DescribeLiveDomainCertBindings(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LiveDomainCertBindings) < 1 {
			break
		}
		instances = append(instances, response.Response.LiveDomainCertBindings...)
		if len(response.Response.LiveDomainCertBindings) < int(limit) {
			break
		}
		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	if *instances[0].CloudCertId != cloudCertId {
		return nil, fmt.Errorf("The CloudCertId[%s] of API [%s] does not equal to specified cloudCertId:[%s]", *instances[0].CloudCertId, request.GetAction(), cloudCertId)
	}

	playDomainCertAttachment = instances[0]

	return
}

func (me *CssService) DeleteCssPlayDomainCertAttachmentById(ctx context.Context, domainName string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewUnBindLiveDomainCertRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().UnBindLiveDomainCert(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssPlayAuthKeyConfigById(ctx context.Context, domainName string) (playAuthKeyConfig *css.PlayAuthKeyInfo, errRet error) {
	logId := getLogId(ctx)

	request := css.NewDescribeLivePlayAuthKeyRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePlayAuthKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.PlayAuthKeyInfo == nil {
		return
	}

	playAuthKeyConfig = response.Response.PlayAuthKeyInfo
	return
}

func (me *CssService) DescribeCssPushAuthKeyConfigById(ctx context.Context, domainName string) (pushAuthKeyConfig *css.PushAuthKeyInfo, errRet error) {
	logId := getLogId(ctx)

	request := css.NewDescribeLivePushAuthKeyRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePushAuthKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.PushAuthKeyInfo == nil {
		return
	}

	pushAuthKeyConfig = response.Response.PushAuthKeyInfo
	return
}
