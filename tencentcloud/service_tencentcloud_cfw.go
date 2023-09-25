package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CfwService struct {
	client *connectivity.TencentCloudClient
}

func (me *CfwService) DescribeCfwAddressTemplateById(ctx context.Context, uuid string) (addressTemplate *cfw.TemplateListInfo, errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewDescribeAddressTemplateListRequest()
	request.Uuid = &uuid

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeAddressTemplateList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	addressTemplate = response.Response.Data[0]
	return
}

func (me *CfwService) DeleteCfwAddressTemplateById(ctx context.Context, uuid string) (errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewDeleteAddressTemplateRequest()
	request.Uuid = &uuid

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DeleteAddressTemplate(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfwService) DescribeCfwBlockIgnoreListById(ctx context.Context, iP, domain, direction, ruleType string) (blockIgnoreRule *cfw.BlockIgnoreRule, errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewDescribeBlockIgnoreListRequest()
	var searchStr string
	if iP != "" {
		searchStr = fmt.Sprintf(`{"domain":"%s"}`, iP)
	} else {
		searchStr = fmt.Sprintf(`{"domain":"%s"}`, domain)
	}

	request.Limit = common.Int64Ptr(20)
	request.Offset = common.Int64Ptr(0)
	request.SearchValue = &searchStr
	request.Direction = &direction
	ruleTypeInt, _ := strconv.ParseUint(ruleType, 10, 64)
	request.RuleType = &ruleTypeInt
	request.By = common.StringPtr("EndTime")
	request.Order = common.StringPtr("desc")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeBlockIgnoreList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	blockIgnoreRule = response.Response.Data[0]
	return
}

func (me *CfwService) DeleteCfwBlockIgnoreListById(ctx context.Context, iP, domain, direction, ruleType string) (errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewDeleteBlockIgnoreRuleListRequest()
	directionInt, _ := strconv.ParseInt(direction, 10, 64)
	if iP != "" {
		request.Rules = []*cfw.IocListData{
			{
				IP:        common.StringPtr(iP),
				Direction: common.Int64Ptr(directionInt),
			},
		}
	} else {
		request.Rules = []*cfw.IocListData{
			{
				Domain:    common.StringPtr(domain),
				Direction: common.Int64Ptr(directionInt),
			},
		}
	}

	ruleTypeInt, _ := strconv.ParseInt(ruleType, 10, 64)
	request.RuleType = common.Int64Ptr(ruleTypeInt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DeleteBlockIgnoreRuleList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfwService) DescribeCfwEdgePolicyById(ctx context.Context, uuid string) (edgePolicy *cfw.DescAcItem, errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewDescribeAclRuleRequest()
	request.Limit = common.Uint64Ptr(20)
	request.Offset = common.Uint64Ptr(0)
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("Id"),
			Values:       common.StringPtrs([]string{uuid}),
			OperatorType: common.Int64Ptr(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeAclRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	edgePolicy = response.Response.Data[0]
	return
}

func (me *CfwService) DeleteCfwEdgePolicyById(ctx context.Context, uuid string) (errRet error) {
	logId := getLogId(ctx)

	request := cfw.NewRemoveAclRuleRequest()
	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	request.RuleUuid = common.Int64Ptrs([]int64{uuidInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().RemoveAclRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
