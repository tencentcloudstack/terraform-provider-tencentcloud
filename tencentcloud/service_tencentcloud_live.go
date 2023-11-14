package tencentcloud

import (
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type LiveService struct {
	client *connectivity.TencentCloudClient
}

func (me *LiveService) DescribeLiveCallbackRuleById(ctx context.Context, domainName string, appName string) (callbackRule *live.CallBackRuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveCallbackRulesRequest()
	request.DomainName = &domainName
	request.AppName = &appName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveCallbackRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.CallBackRuleInfo) < 1 {
		return
	}

	callbackRule = response.Response.CallBackRuleInfo[0]
	return
}

func (me *LiveService) DeleteLiveCallbackRuleById(ctx context.Context, domainName string, appName string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveCallbackRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveCallbackRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveCallbackTemplateById(ctx context.Context, templateId string) (callbackTemplate *live.CallBackTemplateInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveCallbackTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveCallbackTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.CallBackTemplateInfo) < 1 {
		return
	}

	callbackTemplate = response.Response.CallBackTemplateInfo[0]
	return
}

func (me *LiveService) DeleteLiveCallbackTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveCallbackTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveCallbackTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveBackupStreamByFilter(ctx context.Context, param map[string]interface{}) (backupStream []*live.BackupStreamGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeBackupStreamListRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeBackupStreamList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.StreamInfoList) < 1 {
			break
		}
		backupStream = append(backupStream, response.Response.StreamInfoList...)
		if len(response.Response.StreamInfoList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDomainById(ctx context.Context, domainName string) (domain *live.DomainInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveDomainRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainInfo) < 1 {
		return
	}

	domain = response.Response.DomainInfo[0]
	return
}

func (me *LiveService) DeleteLiveDomainById(ctx context.Context, domainName string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveDomainRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveDomainCertById(ctx context.Context, domainName string) (domainCert *live.DomainCertInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveDomainCertRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveDomainCert(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainCertInfo) < 1 {
		return
	}

	domainCert = response.Response.DomainCertInfo[0]
	return
}

func (me *LiveService) DescribeLiveDomainCertBindingsById(ctx context.Context, domainName string) (domainCertBindings *live.LiveDomainCertBindings, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveDomainCertBindingsRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveDomainCertBindings(request)
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

func (me *LiveService) DescribeLiveDomainRefererById(ctx context.Context, domainName string) (domainReferer *live.RefererAuthConfig, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveDomainRefererRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveDomainReferer(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RefererAuthConfig) < 1 {
		return
	}

	domainReferer = response.Response.RefererAuthConfig[0]
	return
}

func (me *LiveService) DescribeLivePlayAuthKeyById(ctx context.Context, domainName string) (playAuthKey *live.PlayAuthKeyInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLivePlayAuthKeyRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLivePlayAuthKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PlayAuthKeyInfo) < 1 {
		return
	}

	playAuthKey = response.Response.PlayAuthKeyInfo[0]
	return
}

func (me *LiveService) DescribeLivePushAuthKeyById(ctx context.Context, domainName string) (pushAuthKey *live.PushAuthKeyInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLivePushAuthKeyRequest()
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLivePushAuthKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PushAuthKeyInfo) < 1 {
		return
	}

	pushAuthKey = response.Response.PushAuthKeyInfo[0]
	return
}

func (me *LiveService) DescribeLiveRecordRuleById(ctx context.Context, templateId string) (recordRule *live.RuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveRecordRulesRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveRecordRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleInfo) < 1 {
		return
	}

	recordRule = response.Response.RuleInfo[0]
	return
}

func (me *LiveService) DeleteLiveRecordRuleById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveRecordRuleRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveRecordRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveRecordTemplateById(ctx context.Context, templateId string) (recordTemplate *live.RecordTemplateInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveRecordTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveRecordTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RecordTemplateInfo) < 1 {
		return
	}

	recordTemplate = response.Response.RecordTemplateInfo[0]
	return
}

func (me *LiveService) DeleteLiveRecordTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveRecordTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveRecordTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveSnapshotRuleById(ctx context.Context, templateId string) (snapshotRule *live.RuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveSnapshotRulesRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveSnapshotRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleInfo) < 1 {
		return
	}

	snapshotRule = response.Response.RuleInfo[0]
	return
}

func (me *LiveService) DeleteLiveSnapshotRuleById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveSnapshotRuleRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveSnapshotRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveSnapshotTemplateById(ctx context.Context, templateId string) (snapshotTemplate *live.SnapshotTemplateInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveSnapshotTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SnapshotTemplateInfo) < 1 {
		return
	}

	snapshotTemplate = response.Response.SnapshotTemplateInfo[0]
	return
}

func (me *LiveService) DeleteLiveSnapshotTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveSnapshotTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveStreamMonitorById(ctx context.Context, monitorId string) (streamMonitor *live.LiveStreamMonitorInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveStreamMonitorRequest()
	request.MonitorId = &monitorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveStreamMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LiveStreamMonitorInfo) < 1 {
		return
	}

	streamMonitor = response.Response.LiveStreamMonitorInfo[0]
	return
}

func (me *LiveService) DeleteLiveStreamMonitorById(ctx context.Context, monitorId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveStreamMonitorRequest()
	request.MonitorId = &monitorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveStreamMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLivePadRuleById(ctx context.Context, domainName string, appName string, streamName string) (padRule *live.RuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLivePadRulesRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLivePadRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleInfo) < 1 {
		return
	}

	padRule = response.Response.RuleInfo[0]
	return
}

func (me *LiveService) DeleteLivePadRuleById(ctx context.Context, domainName string, appName string, streamName string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLivePadRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLivePadRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLivePadTemplateById(ctx context.Context, templateId string) (padTemplate *live.PadTemplate, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLivePadTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLivePadTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PadTemplate) < 1 {
		return
	}

	padTemplate = response.Response.PadTemplate[0]
	return
}

func (me *LiveService) DeleteLivePadTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLivePadTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLivePadTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveBackupStreamById(ctx context.Context, pushDomainName string, appName string, streamName string) (backupStream *live.BackupStreamGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeBackupStreamListRequest()
	request.PushDomainName = &pushDomainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeBackupStreamList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.BackupStreamGroupInfo) < 1 {
		return
	}

	backupStream = response.Response.BackupStreamGroupInfo[0]
	return
}

func (me *LiveService) DescribeLiveTimeshiftRuleById(ctx context.Context, domainName string, appName string, streamName string) (timeshiftRule *live.RuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveTimeShiftRulesRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveTimeShiftRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleInfo) < 1 {
		return
	}

	timeshiftRule = response.Response.RuleInfo[0]
	return
}

func (me *LiveService) DeleteLiveTimeshiftRuleById(ctx context.Context, domainName string, appName string, streamName string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveTimeShiftRuleRequest()
	request.DomainName = &domainName
	request.AppName = &appName
	request.StreamName = &streamName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveTimeShiftRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveTimeshiftTemplateById(ctx context.Context, templateId string) (timeshiftTemplate *live.TimeShiftTemplate, errRet error) {
	logId := getLogId(ctx)

	request := live.NewDescribeLiveTimeShiftTemplatesRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DescribeLiveTimeShiftTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.TimeShiftTemplate) < 1 {
		return
	}

	timeshiftTemplate = response.Response.TimeShiftTemplate[0]
	return
}

func (me *LiveService) DeleteLiveTimeshiftTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := live.NewDeleteLiveTimeShiftTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLiveClient().DeleteLiveTimeShiftTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LiveService) DescribeLiveDescribeDeliverLogDownListByFilter(ctx context.Context, param map[string]interface{}) (DescribeDeliverLogDownList []*live.PushLogInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeDeliverLogDownListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeDeliverLogDownList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LogInfoList) < 1 {
			break
		}
		DescribeDeliverLogDownList = append(DescribeDeliverLogDownList, response.Response.LogInfoList...)
		if len(response.Response.LogInfoList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeLiveStreamMonitorListByFilter(ctx context.Context, param map[string]interface{}) (DescribeLiveStreamMonitorList []*live.DescribeLiveStreamMonitorListResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeLiveStreamMonitorListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Index" {
			request.Index = v.(*uint64)
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
		response, err := me.client.UseLiveClient().DescribeLiveStreamMonitorList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalNum) < 1 {
			break
		}
		DescribeLiveStreamMonitorList = append(DescribeLiveStreamMonitorList, response.Response.TotalNum...)
		if len(response.Response.TotalNum) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeLiveXP2PDetailInfoListByFilter(ctx context.Context, param map[string]interface{}) (DescribeLiveXP2PDetailInfoList []*live.XP2PDetailInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeLiveXP2PDetailInfoListRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeLiveXP2PDetailInfoList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DataInfoList) < 1 {
			break
		}
		DescribeLiveXP2PDetailInfoList = append(DescribeLiveXP2PDetailInfoList, response.Response.DataInfoList...)
		if len(response.Response.DataInfoList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeMonitorReportByFilter(ctx context.Context, param map[string]interface{}) (DescribeMonitorReport []*live.MPSResult, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeMonitorReportRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeMonitorReport(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MPSResult) < 1 {
			break
		}
		DescribeMonitorReport = append(DescribeMonitorReport, response.Response.MPSResult...)
		if len(response.Response.MPSResult) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeLivePadTemplatesByFilter(ctx context.Context, param map[string]interface{}) (DescribeLivePadTemplates []*live.PadTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeLivePadTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeLivePadTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Templates) < 1 {
			break
		}
		DescribeLivePadTemplates = append(DescribeLivePadTemplates, response.Response.Templates...)
		if len(response.Response.Templates) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeLivePullStreamTaskStatusByFilter(ctx context.Context, param map[string]interface{}) (DescribeLivePullStreamTaskStatus []*live.TaskStatusInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeLivePullStreamTaskStatusRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeLivePullStreamTaskStatus(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TaskStatusInfo) < 1 {
			break
		}
		DescribeLivePullStreamTaskStatus = append(DescribeLivePullStreamTaskStatus, response.Response.TaskStatusInfo...)
		if len(response.Response.TaskStatusInfo) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeTimeShiftRecordDetailByFilter(ctx context.Context, param map[string]interface{}) (DescribeTimeShiftRecordDetail []*live.TimeShiftRecord, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeTimeShiftRecordDetailRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeTimeShiftRecordDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordList) < 1 {
			break
		}
		DescribeTimeShiftRecordDetail = append(DescribeTimeShiftRecordDetail, response.Response.RecordList...)
		if len(response.Response.RecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LiveService) DescribeLiveDescribeTimeShiftStreamListByFilter(ctx context.Context, param map[string]interface{}) (DescribeTimeShiftStreamList []*live.DescribeTimeShiftStreamListResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = live.NewDescribeTimeShiftStreamListRequest()
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

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseLiveClient().DescribeTimeShiftStreamList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalSize) < 1 {
			break
		}
		DescribeTimeShiftStreamList = append(DescribeTimeShiftStreamList, response.Response.TotalSize...)
		if len(response.Response.TotalSize) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
