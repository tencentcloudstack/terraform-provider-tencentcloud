package tencentcloud

import (
	"context"
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

func (me *CssService) DescribeCssWatermarkRuleAttachment(ctx context.Context, domainName, appName, streamName, watermarkId string) (watermarkRuleAttachment *css.RuleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveWatermarkRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	// request.DomainName = &domainName
	// request.AppName = &appName
	// request.StreamName = &streamName
	// request.WatermarkId = &watermarkId

	response, err := me.client.UseCssClient().DescribeLiveWatermarkRules(request)
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
	watermarkRuleAttachment = response.Response.Rules[0]
	return
}

func (me *CssService) DetachCssWatermarkRuleAttachment(ctx context.Context, domainName, appName, streamName string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveWatermarkRuleRequest()

	request.DomainName = helper.String(domainName)
	request.AppName = helper.String(appName)
	request.StreamName = helper.String(streamName)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DeleteLiveWatermarkRule(request)
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
