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

func (me *WedataService) DescribeWedataRuleTemplateById(ctx context.Context, projectId string, ruleTemplateId string) (ruleTemplate *wedata.RuleTemplate, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
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

func (me *WedataService) DeleteWedataRuleTemplateById(ctx context.Context, projectId, ruleTemplateId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
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

func (me *WedataService) DescribeWedataRuleTemplatesByFilter(ctx context.Context, param map[string]interface{}) (ruleTemplates []*wedata.RuleTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeRuleTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*uint64)
		}
		if k == "SourceObjectType" {
			request.SourceObjectType = v.(*uint64)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}
		if k == "SourceEngineTypes" {
			request.SourceEngineTypes = v.([]*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ruleTemplates = response.Response.Data

	return
}
