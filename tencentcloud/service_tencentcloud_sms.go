package tencentcloud

import (
	"context"
	"log"
	"strconv"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SmsService struct {
	client *connectivity.TencentCloudClient
}

// sms sign
func (me *SmsService) DescribeSmsSign(ctx context.Context, signId_string string, international *uint64) (sign *sms.DescribeSignListStatus, errRet error) {
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
	// 类型转换
	signId, _ := strconv.ParseUint(signId_string, 10, 64)
	request.SignIdSet = []*uint64{&signId}
	request.International = international

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

func (me *SmsService) DeleteSmsSignById(ctx context.Context, signId_string string) (errRet error) {
	logId := getLogId(ctx)

	request := sms.NewDeleteSmsSignRequest()

	// 类型转换
	signId, _ := strconv.ParseUint(signId_string, 10, 64)
	request.SignId = &signId

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

func (me *SmsService) DescribeSmsTemplate(ctx context.Context, templateId_string string, international *uint64) (template *sms.DescribeTemplateListStatus, errRet error) {
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
	templateId, _ := strconv.ParseUint(templateId_string, 10, 64) //类型转换
	request.TemplateIdSet = []*uint64{&templateId}  // id数组，需进行类型转换
	request.International = international  //添加新的传参 international

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

func (me *SmsService) DeleteSmsTemplateById(ctx context.Context, templateId_string string) (errRet error) {
	logId := getLogId(ctx)

	request := sms.NewDeleteSmsTemplateRequest()

	templateId, _ := strconv.ParseUint(templateId_string, 10, 64) //类型转换
	request.TemplateId = &templateId

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

