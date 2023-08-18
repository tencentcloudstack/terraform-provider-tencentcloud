package tencentcloud

import (
	"context"
	"log"

	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WedataService struct {
	client *connectivity.TencentCloudClient
}

func (me *WedataService) DescribeWedataRuleTemplateById(ctx context.Context, ruleTemplateId string) (ruleTemplate *wedata.RuleTemplate, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeRuleTemplateRequest()
	request.TemplateId = helper.StrToUint64Point(ruleTemplateId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Data != nil {

		ruleTemplate = response.Response.Data
	}

	return
}

func (me *WedataService) DeleteWedataRuleTemplateById(ctx context.Context, ruleTemplateId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleTemplateRequest()
	request.Ids = []*uint64{helper.StrToUint64Point(ruleTemplateId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
