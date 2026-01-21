package cfw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CfwService struct {
	client *connectivity.TencentCloudClient
}

func (me *CfwService) DescribeCfwAddressTemplateById(ctx context.Context, uuid string) (addressTemplate *cfw.TemplateListInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewRemoveAclRuleRequest()
	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	request.RuleUuid = common.Int64Ptrs([]int64{uuidInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().RemoveAclRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw edgePolicy failed, reason:%+v", logId, err)
		return err
	}

	return
}

func (me *CfwService) DescribeCfwNatInstanceById(ctx context.Context, natinsId string) (natInstance *cfw.NatInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatFwInstancesInfoRequest()
	response := cfw.NewDescribeNatFwInstancesInfoResponse()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(10)
	request.Filter = []*cfw.NatFwFilter{
		{
			FilterType:    common.StringPtr("NatinsId"),
			FilterContent: common.StringPtr(natinsId),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = natinsId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient(iacExtInfo).DescribeNatFwInstancesInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.NatinsLst == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe nat fw instances info failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.NatinsLst) < 1 {
		return
	}

	natInstance = response.Response.NatinsLst[0]
	return
}

func (me *CfwService) DescribeCfwEipsById(ctx context.Context, instanceId string) (gwList []string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeCfwEipsRequest()
	response := cfw.NewDescribeCfwEipsResponse()
	request.Mode = common.Uint64Ptr(1)
	request.NatGatewayId = common.StringPtr("ALL")
	request.CfwInstance = common.StringPtr(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeCfwEips(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.NatFwEipList == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cfw eips failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.NatFwEipList) < 1 {
		return
	}

	for _, item := range response.Response.NatFwEipList {
		if item.NatGatewayId != nil {
			gwList = append(gwList, *item.NatGatewayId)
		}
	}

	return
}

func (me *CfwService) DeleteCfwNatInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDeleteNatFwInstanceRequest()
	request.CfwInstance = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DeleteNatFwInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete nat fw instance failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *CfwService) DescribeNatFwVpcDnsLstById(ctx context.Context, instanceId string) (vpcList []string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatFwVpcDnsLstRequest()
	response := cfw.NewDescribeNatFwVpcDnsLstResponse()
	request.NatFwInsId = &instanceId
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(10)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeNatFwVpcDnsLst(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.VpcDnsSwitchLst == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe nat fw vpc dns list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.VpcDnsSwitchLst) < 1 {
		return
	}

	for _, item := range response.Response.VpcDnsSwitchLst {
		if item.VpcId != nil {
			vpcList = append(vpcList, *item.VpcId)
		}
	}

	return
}

func (me *CfwService) DescribeCfwNatPolicyById(ctx context.Context, uuid string) (natPolicy *cfw.DescAcItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatAcRuleRequest()
	response := cfw.NewDescribeNatAcRuleResponse()
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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe nat ac rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Data) < 1 {
		return
	}

	natPolicy = response.Response.Data[0]
	return
}

func (me *CfwService) DeleteCfwNatPolicyById(ctx context.Context, uuid string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewRemoveNatAcRuleRequest()
	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	request.RuleUuid = common.Int64Ptrs([]int64{uuidInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().RemoveNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *CfwService) DescribeCfwVpcInstanceById(ctx context.Context, fwGroupId string) (vpcInstance *cfw.VpcFwGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeFwGroupInstanceInfoRequest()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("FwGroupId"),
			Values:       common.StringPtrs([]string{fwGroupId}),
			OperatorType: common.Int64Ptr(1),
		},
	}
	request.Limit = common.Int64Ptr(10)
	request.Offset = common.Int64Ptr(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = fwGroupId
	response, err := me.client.UseCfwClient(iacExtInfo).DescribeFwGroupInstanceInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.VpcFwGroupLst) < 1 {
		return
	}

	vpcInstance = response.Response.VpcFwGroupLst[0]
	return
}

func (me *CfwService) DeleteCfwVpcInstanceById(ctx context.Context, fwGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDeleteVpcFwGroupRequest()
	request.FwGroupId = &fwGroupId
	request.DeleteFwGroup = common.Int64Ptr(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DeleteVpcFwGroup(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfwService) DescribeFwGroupInstanceInfoById(ctx context.Context, fwGroupId string) (vpcFwGroupInfo *cfw.VpcFwGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeFwGroupInstanceInfoRequest()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(10)
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("FwGroupId"),
			Values:       common.StringPtrs([]string{fwGroupId}),
			OperatorType: common.Int64Ptr(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeFwGroupInstanceInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.VpcFwGroupLst) < 1 {
		return
	}

	vpcFwGroupInfo = response.Response.VpcFwGroupLst[0]
	return
}

func (me *CfwService) DescribeCfwVpcPolicyById(ctx context.Context, uuid string) (vpcPolicy *cfw.VpcRuleItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeVpcAcRuleRequest()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("Id"),
			Values:       common.StringPtrs([]string{uuid}),
			OperatorType: common.Int64Ptr(1),
		},
	}
	request.Limit = common.Uint64Ptr(20)
	request.Offset = common.Uint64Ptr(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeVpcAcRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	vpcPolicy = response.Response.Data[0]
	return
}

func (me *CfwService) DeleteCfwVpcPolicyById(ctx context.Context, uuid string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewRemoveVpcAcRuleRequest()
	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	request.RuleUuids = common.Int64Ptrs([]int64{uuidInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().RemoveVpcAcRule(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfwService) DescribeCfwNatFirewallFwSwitchById(ctx context.Context, natInsId, subnetId string) (natFirewallSwitch *cfw.NatSwitchListData, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatFwSwitchRequest()
	response := cfw.NewDescribeNatFwSwitchResponse()
	request.Offset = common.Int64Ptr(0)
	request.Limit = common.Int64Ptr(20)
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("NatInsId"),
			OperatorType: common.Int64Ptr(1),
			Values:       common.StringPtrs([]string{natInsId}),
		},
		{
			Name:         common.StringPtr("SubnetId"),
			OperatorType: common.Int64Ptr(1),
			Values:       common.StringPtrs([]string{subnetId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeNatFwSwitch(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe nat firewall switch failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Data) < 1 {
		return
	}

	natFirewallSwitch = response.Response.Data[0]
	return
}

func (me *CfwService) DescribeCfwNatFwSwitchesByFilter(ctx context.Context, param map[string]interface{}) (natFwSwitches []*cfw.NatSwitchListData, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cfw.NewDescribeNatFwSwitchRequest()
		response = cfw.NewDescribeNatFwSwitchResponse()
		filters  []*cfw.CommonFilter
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "NatInsId" {
			filters = append(filters, &cfw.CommonFilter{
				Name:         common.StringPtr("NatInsId"),
				OperatorType: common.Int64Ptr(1),
				Values:       common.StringPtrs([]string{v.(string)}),
			})
		}

		if k == "Status" {
			filters = append(filters, &cfw.CommonFilter{
				Name:         common.StringPtr("Status"),
				OperatorType: common.Int64Ptr(1),
				Values:       common.StringPtrs([]string{helper.IntToStr(v.(int))}),
			})
		}

		if k == "Enable" {
			filters = append(filters, &cfw.CommonFilter{
				Name:         common.StringPtr("Enable"),
				OperatorType: common.Int64Ptr(1),
				Values:       common.StringPtrs([]string{helper.IntToStr(v.(int))}),
			})
		}
	}

	if len(filters) > 0 {
		request.Filters = filters
	}

	var (
		offset int64 = 0
		limit  int64 = 20
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCfwClient().DescribeNatFwSwitch(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe nat firewall switch failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Data) < 1 {
			break
		}

		natFwSwitches = append(natFwSwitches, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CfwService) DescribeCfwVpcFirewallSwitchById(ctx context.Context, vpcInsId, switchId string) (vpcFirewallSwitch *cfw.FwGroupSwitchShow, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeVpcFwGroupSwitchRequest()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("SwitchId"),
			Values:       common.StringPtrs([]string{switchId}),
			OperatorType: common.Int64Ptr(1),
		},
		{
			Name:         common.StringPtr("FwGroupId"),
			Values:       common.StringPtrs([]string{vpcInsId}),
			OperatorType: common.Int64Ptr(1),
		},
	}
	request.Limit = common.Uint64Ptr(20)
	request.Offset = common.Uint64Ptr(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeVpcFwGroupSwitch(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SwitchList) < 1 {
		return
	}

	vpcFirewallSwitch = response.Response.SwitchList[0]
	return
}

func (me *CfwService) DescribeCfwVpcFwSwitchesByFilter(ctx context.Context, vpcInsId string) (vpcFirewallSwitch []*cfw.FwGroupSwitchShow, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeVpcFwGroupSwitchRequest()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("FwGroupId"),
			Values:       common.StringPtrs([]string{vpcInsId}),
			OperatorType: common.Int64Ptr(1),
		},
	}

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
		response, err := me.client.UseCfwClient().DescribeVpcFwGroupSwitch(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SwitchList) < 1 {
			break
		}

		vpcFirewallSwitch = append(vpcFirewallSwitch, response.Response.SwitchList...)
		if len(response.Response.SwitchList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CfwService) DescribeCfwEdgeFwSwitchesByFilter(ctx context.Context) (edgeFwSwitches []*cfw.EdgeIpInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cfw.NewDescribeFwEdgeIpsRequest()
	)

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
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCfwClient().DescribeFwEdgeIps(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}

		edgeFwSwitches = append(edgeFwSwitches, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CfwService) DescribeCfwEdgeFirewallSwitchById(ctx context.Context, publicIp string) (edgeFirewallSwitch *cfw.EdgeIpInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeFwEdgeIpsRequest()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("PublicIp"),
			Values:       common.StringPtrs([]string{publicIp}),
			OperatorType: common.Int64Ptr(1),
		},
	}
	request.Limit = helper.Int64(10)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwClient().DescribeFwEdgeIps(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Data) < 1 {
		return
	}

	edgeFirewallSwitch = response.Response.Data[0]
	return
}

func NewCfwService(client *connectivity.TencentCloudClient) CfwService {
	return CfwService{client: client}
}

func (me *CfwService) DescribeSgRuleById(ctx context.Context, ruleUuid string) (ret *cfw.DescribeEnterpriseSecurityGroupRuleListResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeEnterpriseSecurityGroupRuleListRequest()
	commonFilter := &cfw.CommonFilter{
		Name:         helper.String("RuleUuid"),
		OperatorType: helper.Int64(9),
		Values:       []*string{helper.String(ruleUuid)},
	}
	request.SearchFilters = append(request.SearchFilters, commonFilter)

	if err := resourceTencentCloudSgRuleReadPostFillRequest0(ctx, request); err != nil {
		return nil, err
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfwV20190904Client().DescribeEnterpriseSecurityGroupRuleList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *CfwService) DescribeCfwNatPolicyOrderConfigById(ctx context.Context, uuid string) (natPolicy *cfw.DescAcItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatAcRuleRequest()
	response := cfw.NewDescribeNatAcRuleResponse()
	request.Limit = common.Uint64Ptr(100)
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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe nat ac rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Data) < 1 {
		return
	}

	natPolicy = response.Response.Data[0]
	return
}

func (me *CfwService) DescribeCfwNatPolicyOrderConfigs(ctx context.Context) (ret []*cfw.DescAcItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeNatAcRuleRequest()
	response := cfw.NewDescribeNatAcRuleResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCfwClient().DescribeNatAcRule(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe nat ac rule failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Data) < 1 {
			break
		}

		ret = append(ret, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CfwService) DescribeCfwCcnInstanceRegionStatusByFilter(ctx context.Context, param map[string]interface{}) (ret []*cfw.RegionFwStatus, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cfw.NewDescribeCcnInstanceRegionStatusRequest()
		response = cfw.NewDescribeCcnInstanceRegionStatusResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CcnId" {
			request.CcnId = v.(*string)
		}
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "RoutingMode" {
			request.RoutingMode = v.(*uint64)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeCcnInstanceRegionStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ccn instance region status failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.RegionFwStatus
	return
}

func (me *CfwService) DescribeCfwCcnAssociatedInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*cfw.CcnAssociatedInstance, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cfw.NewDescribeCcnAssociatedInstancesRequest()
		response = cfw.NewDescribeCcnAssociatedInstancesResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CcnId" {
			request.CcnId = v.(*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeCcnAssociatedInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ccn associated instances failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.CcnAssociatedInstances
	return
}

func (me *CfwService) DescribeCfwCcnVpcFwSwitchByFilter(ctx context.Context, param map[string]interface{}) (ret []*cfw.InterconnectPair, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cfw.NewDescribeCcnVpcFwSwitchRequest()
		response = cfw.NewDescribeCcnVpcFwSwitchResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CcnId" {
			request.CcnId = v.(*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeCcnVpcFwSwitch(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ccn vpc fw switch failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.InterconnectPairs
	return
}

func (me *CfwService) DescribeCfwClusterVpcFwSwitchById(ctx context.Context, ccnId string) (ret []*cfw.InterconnectPair, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeCcnVpcFwSwitchRequest()
	response := cfw.NewDescribeCcnVpcFwSwitchResponse()
	request.CcnId = &ccnId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeCcnVpcFwSwitch(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ccn vpc fw switch failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.InterconnectPairs
	return
}

func (me *CfwService) DescribeCfwClusterVpcFwSwitchsById(ctx context.Context, ccnId string) (ret *cfw.ClusterSwitchDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeClusterVpcFwSwitchsRequest()
	response := cfw.NewDescribeClusterVpcFwSwitchsResponse()
	request.Filters = []*cfw.CommonFilter{
		{
			Name:         common.StringPtr("InsObj"),
			OperatorType: common.Int64Ptr(1),
			Values:       common.StringPtrs([]string{ccnId}),
		},
	}
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(20)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeClusterVpcFwSwitchs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster vpc fw switchs failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Data) == 0 {
		return
	}

	ret = response.Response.Data[0]
	return
}

func (me *CfwService) DescribeCfwVpcSwitchErrorById(ctx context.Context, ccnId, errKey string) (ret *cfw.SwitchError, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cfw.NewDescribeSwitchErrorRequest()
	response := cfw.NewDescribeSwitchErrorResponse()
	request.FwType = common.StringPtr("VPC_FW")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeSwitchError(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ccn vpc fw switch error failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	for _, item := range response.Response.Data {
		if item != nil && item.ErrIns != nil && *item.ErrIns == ccnId {
			if item.ErrKey != nil && *item.ErrKey == errKey {
				ret = item
				break
			}
		}
	}

	return
}

func (me *CfwService) DescribeCfwIpsModeSwitchById(ctx context.Context) (ret *cfw.DescribeIpsModeSwitchResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cfw.NewDescribeIpsModeSwitchRequest()
	response := cfw.NewDescribeIpsModeSwitchResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCfwClient().DescribeIpsModeSwitch(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ips mode switch failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}
