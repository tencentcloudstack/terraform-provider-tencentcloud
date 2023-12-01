package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SesService struct {
	client *connectivity.TencentCloudClient
}

func (me *SesService) DescribeSesTemplateMetadata(ctx context.Context, templateId uint64) (templatesMetadata *ses.TemplatesMetadata, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListEmailTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	var offset uint64 = 0
	var pageSize uint64 = 100
	templates := make([]*ses.TemplatesMetadata, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSesClient().ListEmailTemplates(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TemplatesMetadata) < 1 {
			break
		}
		templates = append(templates, response.Response.TemplatesMetadata...)
		if len(response.Response.TemplatesMetadata) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(templates) < 1 {
		return
	}

	for _, v := range templates {
		if *v.TemplateID == templateId {
			templatesMetadata = v
			break
		}
	}

	return
}

func (me *SesService) DescribeSesTemplate(ctx context.Context, templateId uint64) (templateResponse *ses.GetEmailTemplateResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewGetEmailTemplateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.TemplateID = &templateId

	response, err := me.client.UseSesClient().GetEmailTemplate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	templateResponse = response.Response
	return
}

func (me *SesService) DeleteSesTemplateById(ctx context.Context, templateID uint64) (errRet error) {
	logId := getLogId(ctx)

	request := ses.NewDeleteEmailTemplateRequest()

	request.TemplateID = &templateID

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSesClient().DeleteEmailTemplate(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SesService) DescribeSesEmailAddress(ctx context.Context, emailAddress string) (emailSender *ses.EmailSender, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListEmailAddressRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseSesClient().ListEmailAddress(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.EmailSenders) < 1 {
		return
	}

	for _, v := range response.Response.EmailSenders {
		if *v.EmailAddress == emailAddress {
			emailSender = v
			break
		}
	}

	return
}

func (me *SesService) DeleteSesEmail_addressById(ctx context.Context, emailAddress string) (errRet error) {
	logId := getLogId(ctx)

	request := ses.NewDeleteEmailAddressRequest()

	request.EmailAddress = &emailAddress

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSesClient().DeleteEmailAddress(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SesService) DescribeSesDomain(ctx context.Context, emailIdentity string) (attributes []*ses.DNSAttributes, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewGetEmailIdentityRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.EmailIdentity = &emailIdentity

	response, err := me.client.UseSesClient().GetEmailIdentity(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Attributes) < 1 {
		return
	}
	attributes = response.Response.Attributes
	return
}

func (me *SesService) DeleteSesDomainById(ctx context.Context, emailIdentity string) (errRet error) {
	logId := getLogId(ctx)

	request := ses.NewDeleteEmailIdentityRequest()

	request.EmailIdentity = &emailIdentity

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSesClient().DeleteEmailIdentity(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SesService) DescribeSesReceiversByFilter(ctx context.Context, param map[string]interface{}) (receivers []*ses.ReceiverData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListReceiversRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Status" {
			request.Status = v.(*uint64)
		}
		if k == "KeyWord" {
			request.KeyWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().ListReceivers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		receivers = append(receivers, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DescribeSesSendTasksByFilter(ctx context.Context, param map[string]interface{}) (sendTasks []*ses.SendTaskData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListSendTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Status" {
			request.Status = v.(*uint64)
		}
		if k == "ReceiverId" {
			request.ReceiverId = v.(*uint64)
		}
		if k == "TaskType" {
			request.TaskType = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().ListSendTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		sendTasks = append(sendTasks, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DescribeSesEmailIdentitiesByFilter(ctx context.Context) (emailIdentities *ses.ListEmailIdentitiesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListEmailIdentitiesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSesClient().ListEmailIdentities(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return nil, nil
	}

	emailIdentities = response.Response

	return
}

func (me *SesService) DescribeSesBlackEmailAddressByFilter(ctx context.Context, param map[string]interface{}) (blackEmailAddress []*ses.BlackEmailAddress, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewListBlackEmailAddressRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "EmailAddress" {
			request.EmailAddress = v.(*string)
		}
		if k == "TaskID" {
			request.TaskID = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().ListBlackEmailAddress(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.BlackList) < 1 {
			break
		}
		blackEmailAddress = append(blackEmailAddress, response.Response.BlackList...)
		if len(response.Response.BlackList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DescribeSesStatisticsReportByFilter(ctx context.Context, param map[string]interface{}) (statisticsReport *ses.GetStatisticsReportResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewGetStatisticsReportRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "ReceivingMailboxType" {
			request.ReceivingMailboxType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSesClient().GetStatisticsReport(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return nil, nil
	}

	statisticsReport = response.Response

	return
}

func (me *SesService) DescribeSesSendEmailStatusByFilter(ctx context.Context, param map[string]interface{}) (sendEmailStatus []*ses.SendEmailStatus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = ses.NewGetSendEmailStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "RequestDate" {
			request.RequestDate = v.(*string)
		}
		if k == "MessageId" {
			request.MessageId = v.(*string)
		}
		if k == "ToEmailAddress" {
			request.ToEmailAddress = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().GetSendEmailStatus(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EmailStatusList) < 1 {
			break
		}

		sendEmailStatus = append(sendEmailStatus, response.Response.EmailStatusList...)
		if len(response.Response.EmailStatusList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DescribeSesReceiverById(ctx context.Context, receiverId string) (Receiver *ses.ReceiverData, errRet error) {
	logId := getLogId(ctx)

	id, err := strconv.Atoi(receiverId)
	if err != nil {
		errRet = fmt.Errorf("[ERROR]%s id data type error: %v", logId, receiverId)
		return
	}

	request := ses.NewListReceiversRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().ListReceivers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		for _, v := range response.Response.Data {
			if *v.ReceiverId == uint64(id) {
				Receiver = v
				return
			}
		}

		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DescribeSesReceiverDetailById(ctx context.Context, receiverId string) (receiverDetail []*ses.ReceiverDetail, errRet error) {
	logId := getLogId(ctx)

	id, err := strconv.Atoi(receiverId)
	if err != nil {
		errRet = fmt.Errorf("[ERROR]%s id data type error: %v", logId, receiverId)
		return
	}

	request := ses.NewListReceiverDetailsRequest()
	request.ReceiverId = helper.IntUint64(id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSesClient().ListReceiverDetails(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		receiverDetail = append(receiverDetail, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SesService) DeleteSesReceiverById(ctx context.Context, receiverId string) (errRet error) {
	logId := getLogId(ctx)

	id, _ := strconv.Atoi(receiverId)

	request := ses.NewDeleteReceiverRequest()
	request.ReceiverId = helper.IntUint64(id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSesClient().DeleteReceiver(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SesService) DescribeSesVerifyDomainById(ctx context.Context, emailIdentity string) (verifyDomain *ses.GetEmailIdentityResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := ses.NewGetEmailIdentityRequest()
	request.EmailIdentity = &emailIdentity

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSesClient().GetEmailIdentity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	verifyDomain = response.Response
	return
}

func (me *SesService) CheckEmailIdentityById(ctx context.Context, emailIdentity string) (errRet error) {
	logId := getLogId(ctx)

	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		verifyDomain, e := me.DescribeSesVerifyDomainById(ctx, emailIdentity)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if verifyDomain == nil {
			return resource.NonRetryableError(fmt.Errorf("emailIdentity %s not exists", emailIdentity))
		}

		if !*verifyDomain.VerifiedForSendingStatus {
			return resource.RetryableError(fmt.Errorf("check emailIdentity status is %v,start retrying ...", *verifyDomain.VerifiedForSendingStatus))
		}
		if *verifyDomain.VerifiedForSendingStatus {
			return nil
		}

		return resource.NonRetryableError(fmt.Errorf("emailIdentity status is %v,we won't wait for it finish", *verifyDomain.VerifiedForSendingStatus))
	})

	if err != nil {
		log.Printf("[CRITAL]%s verifyDomain emailIdentity fail, reason:%s\n ", logId, err.Error())
		errRet = err
		return
	}

	return
}
