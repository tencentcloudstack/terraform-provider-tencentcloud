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
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	
	request.WatermarkId = helper.Int64Uint64(helper.StrToInt64(watermarkId))

	response, err := me.client.UseCssClient().DescribeLiveWatermark(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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

func (me *CssService) DeleteCssWatermarkById(ctx context.Context, watermarkId string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveWatermarkRequest()

	request.WatermarkId = helper.Int64(helper.StrToInt64(watermarkId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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

func (me *CssService) DescribeCssWatermarkRule(ctx context.Context, domainName, appName, streamName, watermarkId string) (watermarkRule *css.RuleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLiveWatermarkRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	// request.DomainName = &domainName
	// request.AppName = &appName
	// request.StreamName = &streamName
	// request.WatermarkId = &watermarkId

	response, err := me.client.UseCssClient().DescribeLiveWatermarkRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Rules) < 1 {
		return
	}
	watermarkRule = response.Response.Rules[0]
	return
}

func (me *CssService) DeleteCssWatermarkRuleById(ctx context.Context, domainName, appName, streamName, watermarkId string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLiveWatermarkRuleRequest()

	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName
	// request.WatermarkId = &watermarkId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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

func (me *CssService) DescribeCssPullStreamTask(ctx context.Context, taskId string) (pullStreamTask *css.PullStreamTaskInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = css.NewDescribeLivePullStreamTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.TaskId = &taskId

	response, err := me.client.UseCssClient().DescribeLivePullStreamTasks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.TaskInfos) < 1 {
		return
	}
	pullStreamTask = response.Response.TaskInfos[0]
	return
}

func (me *CssService) DeleteCssPullStreamTaskById(ctx context.Context, taskId string) (errRet error) {
	logId := getLogId(ctx)

	request := css.NewDeleteLivePullStreamTaskRequest()

	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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
