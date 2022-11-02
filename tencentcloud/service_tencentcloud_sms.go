package tencentcloud

import (
	"context"
	"log"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SmsService struct {
	client *connectivity.TencentCloudClient
}

func (me *SmsService) DescribeSmsSign(ctx context.Context, signId string) (sign *sms.DescribeSignListStatus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sms.NewDescribeSmsSignListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.SignIdSet = []*uint64{helper.Uint64(helper.StrToUInt64(signId))}

	response, err := me.client.UseSmsClient().DescribeSmsSignList(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.DescribeSignListStatusSet) < 1 {
		return
	}
	sign = response.Response.DescribeSignListStatusSet[0]
	return
}

func (me *SmsService) DeleteSmsSignById(ctx context.Context, signId string) (errRet error) {
	logId := getLogId(ctx)

	request := sms.NewDeleteSmsSignRequest()
	request.SignId = helper.Uint64(helper.StrToUInt64(signId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSmsClient().DeleteSmsSign(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SmsService) DescribeSmsTemplate(ctx context.Context, templateId string) (template *sms.DescribeTemplateListStatus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sms.NewDescribeSmsTemplateListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.TemplateIdSet = []*uint64{helper.Uint64(helper.StrToUInt64(templateId))}

	response, err := me.client.UseSmsClient().DescribeSmsTemplateList(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.DescribeTemplateStatusSet) < 1 {
		return
	}
	template = response.Response.DescribeTemplateStatusSet[0]
	return
}

func (me *SmsService) DeleteSmsTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := sms.NewDeleteSmsTemplateRequest()

	request.TemplateId = helper.Uint64(helper.StrToUInt64(templateId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSmsClient().DeleteSmsTemplate(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
