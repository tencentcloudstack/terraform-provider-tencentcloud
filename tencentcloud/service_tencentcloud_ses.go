package tencentcloud

import (
	"context"
	"log"

	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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

func (me *SesService) DescribeSesDomain(ctx context.Context, emailIdentity string) (domain *ses.DNSAttributes, errRet error) {
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
	domain = response.Response.Attributes[0]
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
