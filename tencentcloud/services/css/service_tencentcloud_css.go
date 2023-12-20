package css

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCssService(client *connectivity.TencentCloudClient) CssService {
	return CssService{client: client}
}

type CssService struct {
	client *connectivity.TencentCloudClient
}

func (me *CssService) DescribeCssWatermark(ctx context.Context, watermarkId string) (watermark *css.WatermarkInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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

func (me *CssService) DescribeCssWatermarkRuleAttachment(ctx context.Context, domainName, appName, streamName, watermarkId string) (watermarkRuleAttachment *css.RuleInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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

	for _, rule := range response.Response.Rules {
		if *rule.DomainName == domainName && *rule.AppName == appName && *rule.StreamName == streamName && helper.Int64ToStr(*rule.TemplateId) == watermarkId {
			watermarkRuleAttachment = rule
			return
		}
	}

	return
}

func (me *CssService) DetachCssWatermarkRuleAttachment(ctx context.Context, domainName, appName, streamName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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

func (me *CssService) DescribeCssBackupStreamByFilter(ctx context.Context, param map[string]interface{}) (backupStream []*css.BackupStreamGroupInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeBackupStreamListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StreamName" {
			request.StreamName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DescribeBackupStreamList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.StreamInfoList) < 1 {
		return
	}

	backupStream = response.Response.StreamInfoList

	return
}

func (me *CssService) DescribeCssBackupStreamById(ctx context.Context, pushDomainName string, appName string, streamName string) (backupStream *css.BackupStreamDetailData, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeBackupStreamListRequest()
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeBackupStreamList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.StreamInfoList) < 1 {
		return
	}

	for _, v := range response.Response.StreamInfoList[0].BackupList {
		if *v.AppName == appName && *v.DomainName == pushDomainName && *v.MasterFlag == 1 {
			backupStream = v
		}
	}

	return
}

func (me *CssService) DescribeCssBackupStreamByStreamName(ctx context.Context, streamName string) (backupStream *css.BackupStreamGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeBackupStreamListRequest()
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeBackupStreamList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.StreamInfoList) < 1 {
		return
	}

	backupStream = response.Response.StreamInfoList[0]

	return
}

func (me *CssService) DescribeCssWatermarksByFilter(ctx context.Context, param map[string]interface{}) (watermarks []*css.WatermarkInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeLiveWatermarksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveWatermarks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.WatermarkList) < 1 {
		return
	}

	watermarks = response.Response.WatermarkList

	return
}

func (me *CssService) DescribeCssDeliverLogDownListByFilter(ctx context.Context, param map[string]interface{}) (deliverLogDownList []*css.PushLogInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeDeliverLogDownListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCssClient().DescribeDeliverLogDownList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.LogInfoList) < 1 {
		return
	}

	deliverLogDownList = response.Response.LogInfoList

	return
}

func (me *CssService) DescribeCssStreamMonitorListByFilter(ctx context.Context, param map[string]interface{}) (streamMonitorList []*css.LiveStreamMonitorInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeLiveStreamMonitorListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Index = &offset
		request.Count = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCssClient().DescribeLiveStreamMonitorList(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LiveStreamMonitors) < 1 {
			break
		}
		streamMonitorList = append(streamMonitorList, response.Response.LiveStreamMonitors...)
		if len(response.Response.LiveStreamMonitors) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	return
}

func (me *CssService) DescribeCssXp2pDetailInfoListByFilter(ctx context.Context, param map[string]interface{}) (xp2pDetailInfoList []*css.XP2PDetailInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeLiveXP2PDetailInfoListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "QueryTime" {
			request.QueryTime = v.(*string)
		}
		if k == "Type" {
			request.Type = v.([]*string)
		}
		if k == "StreamNames" {
			request.StreamNames = v.([]*string)
		}
		if k == "Dimension" {
			request.Dimension = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveXP2PDetailInfoList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DataInfoList) < 1 {
		return
	}

	xp2pDetailInfoList = response.Response.DataInfoList

	return
}

func (me *CssService) DescribeCssMonitorReportByFilter(ctx context.Context, param map[string]interface{}) (monitorReport *css.DescribeMonitorReportResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeMonitorReportRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "MonitorId" {
			request.MonitorId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeMonitorReport(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	monitorReport = response.Response

	return
}

func (me *CssService) DescribeCssPadTemplatesByFilter(ctx context.Context, param map[string]interface{}) (padTemplates []*css.PadTemplate, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeLivePadTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePadTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Templates) < 1 {
		return
	}

	padTemplates = response.Response.Templates

	return
}

func (me *CssService) DescribeCssPullStreamTaskStatusByFilter(ctx context.Context, param map[string]interface{}) (pullStreamTaskStatus *css.TaskStatusInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeLivePullStreamTaskStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TaskId" {
			request.TaskId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePullStreamTaskStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.TaskStatusInfo == nil {
		return
	}

	pullStreamTaskStatus = response.Response.TaskStatusInfo

	return
}

func (me *CssService) DescribeCssTimeShiftRecordDetailByFilter(ctx context.Context, param map[string]interface{}) (timeShiftRecordDetail []*css.TimeShiftRecord, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeTimeShiftRecordDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "AppName" {
			request.AppName = v.(*string)
		}
		if k == "StreamName" {
			request.StreamName = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}
		if k == "DomainGroup" {
			request.DomainGroup = v.(*string)
		}
		if k == "TransCodeId" {
			request.TransCodeId = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeTimeShiftRecordDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.RecordList) < 1 {
		return
	}

	timeShiftRecordDetail = response.Response.RecordList

	return
}

func (me *CssService) DescribeCssTimeShiftStreamListByFilter(ctx context.Context, param map[string]interface{}) (timeShiftStreamList []*css.TimeShiftStreamInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = css.NewDescribeTimeShiftStreamListRequest()
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
		if k == "StreamName" {
			request.StreamName = v.(*string)
		}
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "DomainGroup" {
			request.DomainGroup = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeTimeShiftStreamList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.StreamList) < 1 {
		return
	}

	timeShiftStreamList = response.Response.StreamList

	return
}

func (me *CssService) DescribeCssCallbackRuleById(ctx context.Context, templateId int64, domainName string) (callbackRule *css.CallBackRuleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveCallbackRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveCallbackRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	for _, v := range response.Response.Rules {
		if *v.DomainName == domainName && *v.TemplateId == templateId {
			callbackRule = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssCallbackRuleById(ctx context.Context, domainName string, appName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveCallbackRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveCallbackRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssCallbackTemplateById(ctx context.Context, templateId int64) (callbackTemplate *css.CallBackTemplateInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveCallbackTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveCallbackTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	callbackTemplate = response.Response.Template

	return
}

func (me *CssService) DeleteCssCallbackTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveCallbackTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveCallbackTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssDomainCertById(ctx context.Context, domainName string) (domainCert *css.DomainCertInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveDomainCertRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveDomainCert(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	domainCert = response.Response.DomainCertInfo

	return
}

func (me *CssService) DescribeCssDomainCertBindingsById(ctx context.Context, domainName string) (domainCertBindings *css.LiveDomainCertBindings, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveDomainCertBindingsRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveDomainCertBindings(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LiveDomainCertBindings) < 1 {
		return
	}

	domainCertBindings = response.Response.LiveDomainCertBindings[0]

	return
}

func (me *CssService) DescribeCssDomainRefererById(ctx context.Context, domainName string) (domainReferer *css.RefererAuthConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveDomainRefererRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveDomainReferer(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	domainReferer = response.Response.RefererAuthConfig
	return
}

func (me *CssService) DescribeCssRecordRuleById(ctx context.Context, templateId int64, domainName string) (recordRule *css.RuleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveRecordRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveRecordRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	for _, v := range response.Response.Rules {
		if *v.TemplateId == templateId || *v.DomainName == domainName {
			recordRule = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssRecordRuleById(ctx context.Context, domainName, appName, streamName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveRecordRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveRecordRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssRecordTemplateById(ctx context.Context, templateId int64) (recordTemplate *css.RecordTemplateInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveRecordTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveRecordTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	recordTemplate = response.Response.Template

	return
}

func (me *CssService) DeleteCssRecordTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveRecordTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveRecordTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssSnapshotRuleById(ctx context.Context, templateId int64, domainName string) (snapshotRule *css.RuleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveSnapshotRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveSnapshotRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	for _, v := range response.Response.Rules {
		if *v.TemplateId == templateId && *v.DomainName == domainName {
			snapshotRule = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssSnapshotRuleById(ctx context.Context, domainName, appName, streamName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveSnapshotRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveSnapshotRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssSnapshotTemplateById(ctx context.Context, templateId int64) (snapshotTemplate *css.SnapshotTemplateInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveSnapshotTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	snapshotTemplate = response.Response.Template

	return
}

func (me *CssService) DeleteCssSnapshotTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveSnapshotTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssPadRuleAttachmentById(ctx context.Context, templateId int64, domainName string) (padRuleAttachment *css.RuleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLivePadRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePadRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	for _, v := range response.Response.Rules {
		if *v.TemplateId == templateId && *v.DomainName == domainName {
			padRuleAttachment = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssPadRuleAttachmentById(ctx context.Context, domainName, appName, streamName string, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLivePadRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName
	templateIdUint := uint64(templateId)
	request.TemplateId = &templateIdUint

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLivePadRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssPadTemplateById(ctx context.Context, templateId int64) (padTemplate *css.PadTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLivePadTemplateRequest()
	templateIdUint := uint64(templateId)
	request.TemplateId = &templateIdUint

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLivePadTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	padTemplate = response.Response.Template
	return
}

func (me *CssService) DeleteCssPadTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLivePadTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLivePadTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssStreamMonitorById(ctx context.Context, monitorId string) (streamMonitor *css.LiveStreamMonitorInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveStreamMonitorRequest()
	request.MonitorId = &monitorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveStreamMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	streamMonitor = response.Response.LiveStreamMonitor

	return
}

func (me *CssService) DeleteCssStreamMonitorById(ctx context.Context, monitorId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveStreamMonitorRequest()
	request.MonitorId = &monitorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveStreamMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssTimeshiftRuleAttachmentById(ctx context.Context, templateId int64, domainName string) (timeshiftRuleAttachment *css.RuleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveTimeShiftRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveTimeShiftRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	for _, v := range response.Response.Rules {
		if *v.TemplateId == templateId && *v.DomainName == domainName {
			timeshiftRuleAttachment = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssTimeshiftRuleAttachmentById(ctx context.Context, domainName string, appName string, streamName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveTimeShiftRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveTimeShiftRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) DescribeCssTimeshiftTemplateById(ctx context.Context, templateId int64) (timeshiftTemplate *css.TimeShiftTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDescribeLiveTimeShiftTemplatesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DescribeLiveTimeShiftTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Templates) < 1 {
		return
	}

	for _, v := range response.Response.Templates {
		if *v.TemplateId == uint64(templateId) {
			timeshiftTemplate = v
			return
		}
	}

	return
}

func (me *CssService) DeleteCssTimeshiftTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewDeleteLiveTimeShiftTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().DeleteLiveTimeShiftTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) CssRestartPushTaskStateRefreshFunc(taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		paramMap := make(map[string]interface{})
		paramMap["TaskId"] = helper.String(taskId)

		instance, err := me.DescribeCssPullStreamTaskStatusByFilter(ctx, paramMap)

		if err != nil {
			return nil, "", err
		}

		return instance, *instance.RunStatus, nil
	}
}

func (me *CssService) DeleteCssStartStreamMonitorById(ctx context.Context, monitorId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := css.NewStopLiveStreamMonitorRequest()
	request.MonitorId = &monitorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCssClient().StopLiveStreamMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CssService) CssStartStreamMonitorStateRefreshFunc(monitorId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeCssStreamMonitorById(ctx, monitorId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.UInt64ToStr(*object.Status), nil
	}
}
