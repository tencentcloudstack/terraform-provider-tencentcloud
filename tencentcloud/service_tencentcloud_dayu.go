package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DayuService struct {
	client *connectivity.TencentCloudClient
}

func (me *DayuService) DescribeCCSelfdefinePolicies(ctx context.Context, resourceType string, resourceId string, policyName string, policyId string) (infos []*dayu.CCPolicy, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeCCSelfDefinePolicyRequest()

	infos = make([]*dayu.CCPolicy, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit uint64 = 0, 20

	request.Business = &resourceType
	if resourceId != "" {
		request.Id = &resourceId
	}

	request.Offset = &offset
	request.Limit = &limit
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().DescribeCCSelfDefinePolicy(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					//this error case is what sdk returns when the dayu service is overdue
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		if policyName == "" && policyId == "" {
			infos = append(infos, response.Response.Policys...)
		} else {
			for _, policy := range response.Response.Policys {
				if policyName != "" && *policy.Name != policyName {
					continue
				}
				if policyId != "" && *policy.SetId != policyId {
					continue
				}
				infos = append(infos, policy)
			}
		}
		if len(response.Response.Policys) < int(limit) {
			if len(infos) > 0 {
				has = true
			}
			return
		}
		offset += limit
	}
}

func (me *DayuService) DescribeCCSelfdefinePolicy(ctx context.Context, resourceType string, resourceId string, policyId string) (infos *dayu.CCPolicy, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeCCSelfDefinePolicyRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	policies, _, err := me.DescribeCCSelfdefinePolicies(ctx, resourceType, resourceId, "", policyId)
	if err != nil {
		errRet = err
		return
	}

	length := len(policies)
	if length == 0 {
		return
	}
	if length > 1 {
		errRet = fmt.Errorf("Create CC self-define policy returns %d policies", length)
		return
	}

	infos = policies[0]
	has = true
	return
}

func (me *DayuService) CreateCCSelfdefinePolicy(ctx context.Context, resourceType string, resourceId string, ccPolicy dayu.CCPolicy) (policyId string, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateCCSelfDefinePolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Policy = &ccPolicy
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateCCSelfDefinePolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	//describe CC self-define policies to get policy ID
	policies, has, dErr := me.DescribeCCSelfdefinePolicies(ctx, resourceType, resourceId, *ccPolicy.Name, "")
	if dErr != nil {
		errRet = dErr
		return
	}
	if !has {
		errRet = fmt.Errorf("Create CC self-define policy failed")
		return
	}
	if len(policies) != 1 {
		errRet = fmt.Errorf("Create CC self-define policy returns %d policies", len(policies))
		return
	}
	policyId = *policies[0].SetId
	return
}

func (me *DayuService) ModifyCCSelfdefinePolicy(ctx context.Context, resourceType string, resourceId string, policyId string, ccPolicy dayu.CCPolicy) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyCCSelfDefinePolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.SetId = &policyId
	request.Policy = &ccPolicy
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyCCSelfDefinePolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DeleteCCSelfdefinePolicy(ctx context.Context, resourceType string, resourceId string, policyId string) (errRet error) {

	logId := getLogId(ctx)
	request := dayu.NewDeleteCCSelfDefinePolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.SetId = &policyId
	request.Business = &resourceType
	request.Id = &resourceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DeleteCCSelfDefinePolicy(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}
	return
}

func (me *DayuService) CreateDdosPolicy(ctx context.Context, resourceType string, name string, ddosPolicyDropOption []*dayu.DDoSPolicyDropOption, ddosPolicyPortLimit []*dayu.DDoSPolicyPortLimit, ipBlackWhite []*dayu.IpBlackWhite, ddosPacketFilter []*dayu.DDoSPolicyPacketFilter, waterPrintPolicy []*dayu.WaterPrintPolicy) (policyId string, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateDDoSPolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Name = &name
	request.Business = &resourceType
	request.DropOptions = ddosPolicyDropOption
	request.PortLimits = ddosPolicyPortLimit
	request.IpAllowDenys = ipBlackWhite
	request.PacketFilters = ddosPacketFilter
	request.WaterPrint = waterPrintPolicy
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateDDoSPolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if response.Response.PolicyId == nil || *response.Response.PolicyId == "" {
		errRet = errors.New("TencentCloud SDK  return empty DDoS policy Id")
		return
	}
	policyId = *response.Response.PolicyId
	return
}

func (me *DayuService) DescribeDdosPolicies(ctx context.Context, resourceType string, policyId string) (ddosPolicies []*dayu.DDosPolicy, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeDDoSPolicyRequest()

	ddosPolicies = make([]*dayu.DDosPolicy, 0, 100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = &resourceType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribeDDoSPolicy(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "InvalidParameterValue" {
				//this is when resource is not exist
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if policyId != "" {
		for i, p := range response.Response.DDosPolicyList {
			if *p.PolicyId == policyId {
				ddosPolicies = append(ddosPolicies, response.Response.DDosPolicyList[i])
				return
			}
		}
	} else {
		ddosPolicies = append(ddosPolicies, response.Response.DDosPolicyList...)
		return
	}
	return
}

func (me *DayuService) DescribeDdosPolicy(ctx context.Context, resourceType string, policyId string) (ddosPolicy dayu.DDosPolicy, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeDDoSPolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = &resourceType
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribeDDoSPolicy(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "InvalidParameterValue" {
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.DDosPolicyList) == 0 {
		return
	}

	for i, p := range response.Response.DDosPolicyList {
		if *p.PolicyId == policyId {
			has = true
			ddosPolicy = *response.Response.DDosPolicyList[i]
			return
		}
	}

	return
}

func flattenDdosDropOptionList(list []*dayu.DDoSPolicyDropOption) (mapping []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"drop_tcp":           *v.DropTcp > 0,
			"drop_udp":           *v.DropUdp > 0,
			"drop_icmp":          *v.DropIcmp > 0,
			"drop_other":         *v.DropOther > 0,
			"drop_abroad":        *v.DropAbroad > 0,
			"check_sync_conn":    *v.CheckSyncConn > 0,
			"s_new_limit":        int(*v.SdNewLimit),
			"d_new_limit":        int(*v.DstNewLimit),
			"s_conn_limit":       int(*v.SdConnLimit),
			"d_conn_limit":       int(*v.DstConnLimit),
			"bad_conn_threshold": int(*v.BadConnThreshold),
			"null_conn_enable":   *v.NullConnEnable > 0,
			"conn_timeout":       int(*v.ConnTimeout),
			"syn_rate":           int(*v.SynRate),
			"syn_limit":          int(*v.SynLimit),
			"tcp_mbps_limit":     int(*v.DTcpMbpsLimit),
			"udp_mbps_limit":     int(*v.DUdpMbpsLimit),
			"icmp_mbps_limit":    int(*v.DIcmpMbpsLimit),
			"other_mbps_limit":   int(*v.DOtherMbpsLimit),
		}
		result = append(result, mapping)
	}
	return result
}

func flattenCCRuleList(list []*dayu.CCRule) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"skey":     v.Skey,
			"operator": v.Operator,
			"value":    v.Value,
		}
		result = append(result, mapping)
	}
	return result
}

func flattenDdosPortLimitList(list []*dayu.DDoSPolicyPortLimit) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"protocol": v.Protocol,
			"action":   v.Action,
			"kind":     int(*v.Kind),
		}
		if int(*v.Kind) == 0 || int(*v.Kind) == 2 {
			mapping["start_port"] = int(*v.DPortStart)
			mapping["end_port"] = int(*v.DPortEnd)
		} else {
			mapping["start_port"] = int(*v.SPortStart)
			mapping["end_port"] = int(*v.SPortEnd)
		}
		result = append(result, mapping)
	}
	return result
}

func flattenDdosPacketFilterList(list []*dayu.DDoSPolicyPacketFilter) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"protocol":       v.Protocol,
			"action":         v.Action,
			"d_start_port":   int(*v.DportStart),
			"d_end_port":     int(*v.DportEnd),
			"s_start_port":   int(*v.SportStart),
			"s_end_port":     int(*v.SportEnd),
			"pkt_length_max": int(*v.PktlenMax),
			"pkt_length_min": int(*v.PktlenMin),
			"match_begin":    v.MatchBegin,
			"match_type":     v.MatchType,
			"match_str":      v.Str,
			"is_include":     *v.IsNot > 0,
			"depth":          int(*v.Depth),
			"offset":         int(*v.Offset),
		}
		result = append(result, mapping)
	}
	return result
}

func flattenIpBlackWhiteList(list []*dayu.IpBlackWhite) (reB []string, reW []string) {
	reB = make([]string, 0)
	reW = make([]string, 0)
	for _, v := range list {
		if *v.Type == DAYU_IP_TYPE_BLACK {
			reB = append(reB, *v.Ip)
		}
		if *v.Type == DAYU_IP_TYPE_WHITE {
			reW = append(reW, *v.Ip)
		}

	}
	return
}

func flattenWaterPrintPolicyList(list []*dayu.WaterPrintPolicy) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"tcp_port_list": helper.StringsInterfaces(v.TcpPortList),
			"udp_port_list": helper.StringsInterfaces(v.UdpPortList),
			"offset":        int(*v.Offset),
			"auto_remove":   *v.RemoveSwitch > 0,
			"open_switch":   *v.OpenStatus > 0,
		}
		result = append(result, mapping)
	}
	return result
}

func flattenWaterPrintKeyList(list []*dayu.WaterPrintKey) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"id":          *v.KeyId,
			"content":     *v.KeyContent,
			"create_time": *v.CreateTime,
			"open_switch": *v.OpenStatus > 0,
		}
		result = append(result, mapping)
	}
	return result
}

func setDdosPolicyDropOption(mapping []interface{}) (result []*dayu.DDoSPolicyDropOption, err error) {
	result = make([]*dayu.DDoSPolicyDropOption, 0, len(mapping))
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		var r dayu.DDoSPolicyDropOption
		r.DropTcp = helper.BoolToInt64Pointer(v["drop_tcp"].(bool))
		r.DropUdp = helper.BoolToInt64Pointer(v["drop_udp"].(bool))
		r.DropIcmp = helper.BoolToInt64Pointer(v["drop_icmp"].(bool))
		r.DropOther = helper.BoolToInt64Pointer(v["drop_other"].(bool))
		r.DropAbroad = helper.BoolToInt64Pointer(v["drop_abroad"].(bool))
		r.CheckSyncConn = helper.BoolToInt64Pointer(v["check_sync_conn"].(bool))
		r.SdNewLimit = helper.IntUint64(v["s_new_limit"].(int))
		r.DstNewLimit = helper.IntUint64(v["d_new_limit"].(int))
		r.SdConnLimit = helper.IntUint64(v["s_conn_limit"].(int))
		r.DstConnLimit = helper.IntUint64(v["d_conn_limit"].(int))
		r.BadConnThreshold = helper.IntUint64(v["bad_conn_threshold"].(int))
		r.NullConnEnable = helper.BoolToInt64Pointer(v["null_conn_enable"].(bool))
		r.ConnTimeout = helper.IntUint64(v["conn_timeout"].(int))
		r.SynRate = helper.IntUint64(v["syn_rate"].(int))
		r.SynLimit = helper.IntUint64((v["syn_limit"]).(int))
		r.DTcpMbpsLimit = helper.IntUint64(v["tcp_mbps_limit"].(int))
		r.DUdpMbpsLimit = helper.IntUint64(v["udp_mbps_limit"].(int))
		r.DIcmpMbpsLimit = helper.IntUint64(v["icmp_mbps_limit"].(int))
		r.DOtherMbpsLimit = helper.IntUint64(v["other_mbps_limit"].(int))
		result = append(result, &r)
	}
	return
}

func setDdosPolicyPortLimit(mapping []interface{}) (result []*dayu.DDoSPolicyPortLimit, err error) {
	result = make([]*dayu.DDoSPolicyPortLimit, 0, len(mapping))
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		var r dayu.DDoSPolicyPortLimit
		startPort := v["start_port"].(int)
		endPort := v["end_port"].(int)
		kind := v["kind"].(int)
		if startPort > endPort {
			err = fmt.Errorf("The `start_port` should not be greater than `end_port`.")
			return
		}
		if kind == 0 || kind == 2 {
			r.DPortStart = helper.IntUint64(startPort)
			r.DPortEnd = helper.IntUint64(endPort)
		} else if kind == 1 {
			r.SPortStart = helper.IntUint64(startPort)
			r.SPortEnd = helper.IntUint64(endPort)
		}

		r.Protocol = helper.String(v["protocol"].(string))
		r.Action = helper.String(v["action"].(string))
		r.Kind = helper.IntUint64(kind)
		result = append(result, &r)
	}
	return
}

func setIpBlackWhite(blackIps []interface{}, whiteIps []interface{}) (result []*dayu.IpBlackWhite, err error) {
	result = make([]*dayu.IpBlackWhite, 0, len(blackIps)+len(whiteIps))
	for _, vv := range blackIps {
		var r dayu.IpBlackWhite
		r.Ip = helper.String(vv.(string))
		r.Type = helper.String(DAYU_IP_TYPE_BLACK)
		result = append(result, &r)
	}
	for _, vv := range whiteIps {
		var r dayu.IpBlackWhite
		r.Ip = helper.String(vv.(string))
		r.Type = helper.String(DAYU_IP_TYPE_WHITE)
		result = append(result, &r)
	}
	return
}

func setDdosPolicyPacketFilter(mapping []interface{}) (result []*dayu.DDoSPolicyPacketFilter, err error) {
	result = make([]*dayu.DDoSPolicyPacketFilter, 0, len(mapping))
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		var r dayu.DDoSPolicyPacketFilter
		dStartPort := v["d_start_port"].(int)
		dEndPort := v["d_end_port"].(int)
		sStartPort := v["s_start_port"].(int)
		sEndPort := v["s_end_port"].(int)
		if dStartPort > dEndPort {
			err = fmt.Errorf("The `d_start_port` should not be greater than `d_end_port`.")
			return
		}
		if sStartPort > sEndPort {
			err = fmt.Errorf("The `s_start_port` should not be greater than `s_end_port`.")
			return
		}
		pktLenMax := v["pkt_length_max"].(int)
		pktLenMin := v["pkt_length_min"].(int)
		if pktLenMax < pktLenMin {
			err = fmt.Errorf("The `pkt_length_min` should not be greater than `pkt_length_max`.")
			return
		}
		r.Protocol = helper.String(v["protocol"].(string))
		r.DportStart = helper.IntUint64(dStartPort)
		r.DportEnd = helper.IntUint64(dEndPort)
		r.SportStart = helper.IntUint64(sStartPort)
		r.SportEnd = helper.IntUint64(sEndPort)
		r.Action = helper.String(v["action"].(string))
		r.IsNot = helper.BoolToInt64Pointer(v["is_include"].(bool))
		r.PktlenMax = helper.IntUint64(pktLenMax)
		r.PktlenMin = helper.IntUint64(pktLenMin)
		r.MatchBegin = helper.String(v["match_begin"].(string))
		r.MatchType = helper.String(v["match_type"].(string))
		r.Str = helper.String(v["match_str"].(string))
		r.Depth = helper.IntUint64(v["depth"].(int))
		r.Offset = helper.IntUint64(v["offset"].(int))
		result = append(result, &r)
	}
	return
}

func setWaterPrintPolicy(mapping []interface{}) (result []*dayu.WaterPrintPolicy, err error) {
	result = make([]*dayu.WaterPrintPolicy, 0, len(mapping))
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		var r dayu.WaterPrintPolicy
		tcpPortList := v["tcp_port_list"].([]interface{})
		r.TcpPortList = make([]*string, 0, len(tcpPortList))
		for _, tcpPort := range tcpPortList {
			r.TcpPortList = append(r.TcpPortList, helper.String(tcpPort.(string)))
		}
		udpPortList := v["udp_port_list"].([]interface{})
		r.UdpPortList = make([]*string, 0, len(udpPortList))
		for _, udpPort := range udpPortList {
			r.UdpPortList = append(r.UdpPortList, helper.String(udpPort.(string)))
		}
		r.RemoveSwitch = helper.BoolToInt64Pointer(v["auto_remove"].(bool))
		r.OpenStatus = helper.BoolToInt64Pointer(v["open_switch"].(bool))
		r.Offset = helper.IntUint64(v["offset"].(int))
		result = append(result, &r)
	}
	return
}

func (me *DayuService) DeleteDdosPolicy(ctx context.Context, resourceType string, policyId string) (errRet error) {

	logId := getLogId(ctx)
	request := dayu.NewDeleteDDoSPolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.PolicyId = &policyId
	request.Business = &resourceType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DeleteDDoSPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}
	return
}

func (me *DayuService) ModifyDdosPolicy(ctx context.Context, resourceType string, policyId string, ddosPolicyDropOption []*dayu.DDoSPolicyDropOption, ddosPolicyPortLimit []*dayu.DDoSPolicyPortLimit, ipBlackWhite []*dayu.IpBlackWhite, ddosPacketFilter []*dayu.DDoSPolicyPacketFilter, waterPrintPolicy []*dayu.WaterPrintPolicy) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyDDoSPolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = &resourceType
	request.PolicyId = &policyId
	request.DropOptions = ddosPolicyDropOption
	request.PortLimits = ddosPolicyPortLimit
	request.IpAllowDenys = ipBlackWhite
	request.PacketFilters = ddosPacketFilter
	request.WaterPrint = waterPrintPolicy
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyDDoSPolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) ModifyDdosPolicyName(ctx context.Context, resourceType string, policyId string, name string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyDDoSPolicyNameRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = &resourceType
	request.PolicyId = &policyId
	request.Name = helper.String(url.QueryEscape(name))

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyDDoSPolicyName(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}
	return
}

func (me *DayuService) CreateDdosPolicyCase(ctx context.Context, request *dayu.CreateDDoSPolicyCaseRequest) (sceneId string, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateDDoSPolicyCase(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if response.Response.SceneId == nil || *response.Response.SceneId == "" {
		errRet = errors.New("TencentCloud SDK  return empty DDoS policy case Id")
		return
	}
	sceneId = *response.Response.SceneId
	return
}

func (me *DayuService) DescribeDdosPolicyCase(ctx context.Context, resourceType string, sceneId string) (ddosPolicyCase dayu.KeyValueRecord, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribePolicyCaseRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.SceneId = &sceneId
	request.Business = &resourceType
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribePolicyCase(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "InvalidParameterValue" {
				//this is when resource is not exist
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.CaseList) == 0 {
		return
	}
	if len(response.Response.CaseList) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d appInfo with one applicationId %s",
			len(response.Response.CaseList), sceneId)
		return
	}
	ddosPolicyCase = *response.Response.CaseList[0]
	has = true
	return
}

func (me *DayuService) DeleteDdosPolicyCase(ctx context.Context, resourceType string, sceneId string) (errRet error) {

	logId := getLogId(ctx)
	request := dayu.NewDeleteDDoSPolicyCaseRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.SceneId = &sceneId
	request.Business = &resourceType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DeleteDDoSPolicyCase(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		//describe the scene
		_, has, err := me.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		if err != nil {
			errRet = err
			return
		}
		if !has {
			return
		}
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}
	return
}

func (me *DayuService) ModifyDdosPolicyCase(ctx context.Context, request *dayu.ModifyDDoSPolicyCaseRequest) (errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyDDoSPolicyCase(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DescribeDdosPolicyAttachments(ctx context.Context, resourceId string, resourceType string, policyId string) (attachments []map[string]interface{}, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeDDoSPolicyRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if resourceId != "" {
		request.Id = &resourceId
	}
	request.Business = &resourceType
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribeDDoSPolicy(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "InvalidParameterValue" {
				//this is when resource is not exist
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if len(response.Response.DDosPolicyList) == 0 {
		return
	}

	for _, policy := range response.Response.DDosPolicyList {
		if policyId != "" && *policy.PolicyId != policyId {
			continue
		}
		for _, resource := range policy.BoundResources {
			attachments = append(attachments, map[string]interface{}{"resource_id": *resource, "policy_id": *policy.PolicyId, "resource_type": resourceType})
		}
	}
	if len(attachments) > 0 {
		has = true
	}
	return
}

func (me *DayuService) BindDdosPolicy(ctx context.Context, resourceId string, resourceType string, policyId string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyResBindDDoSPolicyRequest()
	request.PolicyId = &policyId
	request.Business = &resourceType
	request.Id = &resourceId
	request.Method = helper.String("bind")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyResBindDDoSPolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) UnbindDdosPolicy(ctx context.Context, resourceId string, resourceType string, policyId string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyResBindDDoSPolicyRequest()
	request.PolicyId = &policyId
	request.Business = &resourceType
	request.Id = &resourceId
	request.Method = helper.String("unbind")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyResBindDDoSPolicy(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		//when resource is not exist, return
		if *response.Response.Success.Code == "InvalidParameterValue" && *response.Response.Success.Message == "resource not exist" {
			return
		}
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DescribeL7Rules(ctx context.Context, resourceType string, resourceId string, ruleDomain string, ruleId string, protocol string) (infos []*dayu.L7RuleEntry, healths []*dayu.L7RuleHealth, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribleL7RulesRequest()

	infos = make([]*dayu.L7RuleEntry, 0, 100)
	healths = make([]*dayu.L7RuleHealth, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit uint64 = 0, 20

	request.Business = &resourceType
	request.Id = &resourceId
	if protocol != "" {
		request.ProtocolList = []*string{&protocol}
	}

	if ruleDomain != "" {
		request.Domain = &ruleDomain
	}
	request.Offset = &offset
	request.Limit = &limit
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().DescribleL7Rules(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		if ruleId == "" {
			infos = append(infos, response.Response.Rules...)
			healths = append(healths, response.Response.Healths...)
		} else {
			for _, rule := range response.Response.Rules {
				if *rule.RuleId != ruleId {
					continue
				}
				infos = append(infos, rule)
				//get right health, the SDK returns with no order
				var theHealth dayu.L7RuleHealth
				for _, health := range response.Response.Healths {
					if *health.RuleId != *rule.RuleId {
						continue
					}
					theHealth = *health
				}
				healths = append(healths, &theHealth)
			}
		}
		if len(response.Response.Rules) < int(limit) {
			if len(infos) > 0 {
				has = true
			}
			return
		}
		offset += limit
	}
}

func (me *DayuService) DescribeL7Rule(ctx context.Context, resourceType string, resourceId string, ruleId string) (infos *dayu.L7RuleEntry, health *dayu.L7RuleHealth, has bool, errRet error) {
	policies, healths, _, err := me.DescribeL7Rules(ctx, resourceType, resourceId, "", ruleId, "")
	if err != nil {
		errRet = err
		return
	}

	length := len(policies)
	if length == 0 {
		return
	}
	if length > 1 {
		errRet = fmt.Errorf("Create l7 rule returns %d rules", length)
		return
	}

	infos = policies[0]
	if len(healths) > 0 {
		health = healths[0]
	}
	has = true
	return
}

func (me *DayuService) SetL7Health(ctx context.Context, resourceType string, resourceId string, healthCheck dayu.L7HealthConfig) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateL7HealthConfigRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.HealthConfig = []*dayu.L7HealthConfig{&healthCheck}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateL7HealthConfig(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DescribeL7Health(ctx context.Context, resourceType string, resourceId string, ruleId string) (healthCheck *dayu.L7HealthConfig, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeL7HealthConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Id = &resourceId
	request.Business = &resourceType
	request.RuleIdList = []*string{&ruleId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribeL7HealthConfig(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "InvalidParameterValue" {
				//this is when resource is not exist
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	healthChecks := response.Response.HealthConfig
	length := len(healthChecks)
	if length == 0 {
		return
	}
	if length > 1 {
		errRet = fmt.Errorf("Get L7 health check returns %d healthchecks", length)
	}

	healthCheck = healthChecks[0]
	has = true
	return
}

func (me *DayuService) CreateL7Rule(ctx context.Context, resourceType string, resourceId string, rule dayu.L7RuleEntry) (ruleId string, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateL7RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Rules = []*dayu.L7RuleEntry{&rule}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateL7Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	//describe rules to get rule ID
	rules, _, has, dErr := me.DescribeL7Rules(ctx, resourceType, resourceId, *rule.Domain, "", "")
	if dErr != nil {
		errRet = dErr
		return
	}
	if !has {
		errRet = fmt.Errorf("Create L7 rule failed")
		return
	}
	if len(rules) != 1 {
		errRet = fmt.Errorf("Create L7 rule returns %d rules", len(rules))
	}
	ruleId = *rules[0].RuleId
	return
}

func (me *DayuService) ModifyL7Rule(ctx context.Context, resourceType string, resourceId string, rule dayu.L7RuleEntry) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyL7RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Rule = &rule
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyL7Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func flattenSourceList(list []*dayu.L4RuleSource) (re []map[string]interface{}) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		mapping := map[string]interface{}{
			"weight": v.Weight,
			"source": v.Source,
		}
		result = append(result, mapping)
	}

	return result
}

func (me *DayuService) DeleteL7Rule(ctx context.Context, resourceType string, resourceId string, ruleId string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDeleteL7RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.RuleIdList = []*string{&ruleId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DeleteL7Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) SetRuleSwitch(ctx context.Context, resourceType string, resourceId string, ruleId string, switchFlag bool, protocol string) (errRet error) {
	logId := getLogId(ctx)
	if protocol == DAYU_L7_RULE_PROTOCOL_HTTP {
		request := dayu.NewModifyCCHostProtectionRequest()
		defer func() {
			if errRet != nil {
				log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
			}
		}()
		request.Id = &resourceId
		request.Business = &resourceType
		request.RuleId = &ruleId
		if switchFlag {
			request.Method = helper.String("open")
		} else {
			request.Method = helper.String("close")
		}
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().ModifyCCHostProtection(request)
		if err != nil {
			errRet = err
			return
		}

		if response == nil || response.Response == nil || response.Response.Success == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}

		if *response.Response.Success.Code != "Success" {
			errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
			return
		}
	} else {
		request := dayu.NewModifyCCThresholdRequest()
		defer func() {
			if errRet != nil {
				log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
			}
		}()
		request.Id = &resourceId
		request.Business = &resourceType
		request.RuleId = &ruleId
		request.Protocol = &protocol

		if !switchFlag {
			request.Threshold = helper.IntUint64(DAYU_L7_HTTPS_SWITCH_OFF)
		} else {
			//this default value can be a request value that asks the whole threshold value if needed
			request.Threshold = helper.IntUint64(DAYU_L7_HTTPS_SWITCH_ON_DEFAULT)
		}
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().ModifyCCThreshold(request)
		if err != nil {
			errRet = err
			return
		}

		if response == nil || response.Response == nil || response.Response.Success == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}

		if *response.Response.Success.Code != "Success" {
			errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
			return
		}
	}
	return
}
func (me *DayuService) DescribeNewL4Rules(ctx context.Context, business string, extendParams map[string]interface{}) (infos []*dayu.NewL4RuleEntry, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeNewL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = common.StringPtr(business)

	var offset, limit uint64 = 0, 20
	request.Offset = common.Uint64Ptr(offset)
	request.Limit = common.Uint64Ptr(limit)

	if ip, ok := extendParams["ip"]; ok {
		request.Ip = common.StringPtr(ip.(string))
	}
	if virtualPort, ok := extendParams["virtual_port"]; ok {
		request.VirtualPort = common.Uint64Ptr(uint64(virtualPort.(int)))
	}

	response, errRet := me.client.UseDayuClient().DescribeNewL4Rules(request)
	if errRet != nil {
		return
	}
	infos = response.Response.Rules
	return

}

func (me *DayuService) DeleteNewL4Rules(ctx context.Context, business string, id string, ip string, ruleIds []string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDeleteNewL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = common.StringPtr(business)
	request.Rule = []*dayu.L4DelRule{
		{
			Id:         common.StringPtr(id),
			Ip:         common.StringPtr(ip),
			RuleIdList: common.StringPtrs(ruleIds),
		},
	}
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().DeleteNewL4Rules(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if *response.Response.Success.Code == "Success" {
			return
		}
	}
}

func (me *DayuService) ModifyNewL4Rule(ctx context.Context, business string, id string, singleRule interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyNewL4RuleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	rule := singleRule.(map[string]interface{})
	tmpRule := dayu.L4RuleEntry{}
	tmpRule.Protocol = common.StringPtr(rule["protocol"].(string))
	tmpRule.SourcePort = common.Uint64Ptr((uint64(rule["source_port"].(int))))
	tmpRule.VirtualPort = common.Uint64Ptr((uint64(rule["virtual_port"].(int))))
	tmpRule.KeepTime = common.Uint64Ptr((uint64(rule["keeptime"].(int))))
	tmpRule.RuleId = common.StringPtr((rule["rule_id"].(string)))
	tmpRule.LbType = common.Uint64Ptr((uint64(rule["lb_type"].(int))))
	if rule["keep_enable"].(bool) {
		tmpRule.KeepEnable = common.Uint64Ptr((uint64(1)))
	} else {
		tmpRule.KeepEnable = common.Uint64Ptr((uint64(0)))
	}
	tmpRule.SourceType = common.Uint64Ptr((uint64(rule["source_type"].(int))))
	if rule["remove_switch"].(bool) {
		tmpRule.RemoveSwitch = common.Uint64Ptr((uint64(1)))
	} else {
		tmpRule.RemoveSwitch = common.Uint64Ptr((uint64(0)))
	}
	sourceList := rule["source_list"].([]interface{})
	tmpRule.SourceList = make([]*dayu.L4RuleSource, 0)
	for _, singleSource := range sourceList {
		source := singleSource.(map[string]interface{})
		tSource := source["source"].(string)
		tWeight := uint64(source["weight"].(int))
		tmpRule.SourceList = append(tmpRule.SourceList, &dayu.L4RuleSource{Source: &tSource, Weight: &tWeight})
	}

	request.Business = common.StringPtr(business)
	request.Id = common.StringPtr(id)
	request.Rule = &tmpRule

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().ModifyNewL4Rule(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if *response.Response.Success.Code == "Success" {
			return
		}
	}
}

func (me *DayuService) CreateNewL4Rules(ctx context.Context, business string, id string, vip string, ruleList []interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateNewL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	rules := make([]*dayu.L4RuleEntry, 0)
	for _, singleRule := range ruleList {
		rule := singleRule.(map[string]interface{})
		tmpRule := dayu.L4RuleEntry{}
		tmpRule.Protocol = common.StringPtr(rule["protocol"].(string))
		tmpRule.SourcePort = common.Uint64Ptr((uint64(rule["source_port"].(int))))
		tmpRule.VirtualPort = common.Uint64Ptr((uint64(rule["virtual_port"].(int))))
		tmpRule.KeepTime = common.Uint64Ptr((uint64(rule["keeptime"].(int))))
		tmpRule.RuleName = common.StringPtr((rule["rule_name"].(string)))
		tmpRule.LbType = common.Uint64Ptr((uint64(rule["lb_type"].(int))))
		if rule["keep_enable"].(bool) {
			tmpRule.KeepEnable = common.Uint64Ptr((uint64(1)))
		} else {
			tmpRule.KeepEnable = common.Uint64Ptr((uint64(0)))
		}
		tmpRule.SourceType = common.Uint64Ptr((uint64(rule["source_type"].(int))))
		if rule["remove_switch"].(bool) {
			tmpRule.RemoveSwitch = common.Uint64Ptr((uint64(1)))
		} else {
			tmpRule.RemoveSwitch = common.Uint64Ptr((uint64(0)))
		}
		sourceList := rule["source_list"].([]interface{})
		tmpRule.SourceList = make([]*dayu.L4RuleSource, 0)
		for _, singleSource := range sourceList {
			source := singleSource.(map[string]interface{})
			tSource := source["source"].(string)
			tWeight := uint64(source["weight"].(int))
			tmpRule.SourceList = append(tmpRule.SourceList, &dayu.L4RuleSource{Source: &tSource, Weight: &tWeight})
		}
		rules = append(rules, &tmpRule)
	}
	request.Business = common.StringPtr(business)
	request.IdList = common.StringPtrs([]string{id})
	request.VipList = common.StringPtrs([]string{vip})
	request.Rules = rules

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().CreateNewL4Rules(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if *response.Response.Success.Code == "Success" {
			return
		}
	}
}

func (me *DayuService) DescribeL4Rules(ctx context.Context, resourceType string, resourceId string, ruleName string, ruleId string) (infos []*dayu.L4RuleEntry, healths []*dayu.L4RuleHealth, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribleL4RulesRequest()

	infos = make([]*dayu.L4RuleEntry, 0, 100)
	healths = make([]*dayu.L4RuleHealth, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit uint64 = 0, 20

	request.Business = &resourceType
	request.Id = &resourceId

	request.Offset = &offset
	request.Limit = &limit
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDayuClient().DescribleL4Rules(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		if ruleId == "" && ruleName == "" {
			infos = append(infos, response.Response.Rules...)
			healths = append(healths, response.Response.Healths...)
		} else {
			for _, rule := range response.Response.Rules {
				if ruleId != "" && *rule.RuleId != ruleId {
					continue
				}
				if ruleName != "" && *rule.RuleName != ruleName {
					continue
				}
				infos = append(infos, rule)

				//get right health, the SDK returns with no order
				var theHealth dayu.L4RuleHealth
				for _, health := range response.Response.Healths {
					if *health.RuleId != *rule.RuleId {
						continue
					}
					theHealth = *health
				}
				healths = append(healths, &theHealth)
			}
		}
		if len(response.Response.Rules) < int(limit) {
			if len(infos) > 0 {
				has = true
			}
			return
		}
		offset += limit
	}
}

func (me *DayuService) DescribeL4Rule(ctx context.Context, resourceType string, resourceId string, ruleId string) (infos *dayu.L4RuleEntry, health *dayu.L4RuleHealth, has bool, errRet error) {
	policies, healths, _, err := me.DescribeL4Rules(ctx, resourceType, resourceId, "", ruleId)
	if err != nil {
		errRet = err
		return
	}

	length := len(policies)
	if length == 0 {
		return
	}
	if length > 1 {
		errRet = fmt.Errorf("Create L4 rule returns %d rules", length)
		return
	}

	infos = policies[0]
	if len(healths) > 0 {
		health = healths[0]
	}
	has = true
	return
}

func (me *DayuService) SetL4Health(ctx context.Context, resourceType string, resourceId string, healthCheck dayu.L4HealthConfig) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateL4HealthConfigRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.HealthConfig = []*dayu.L4HealthConfig{&healthCheck}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateL4HealthConfig(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DescribeL4Health(ctx context.Context, resourceType string, resourceId string, ruleId string) (healthCheck *dayu.L4HealthConfig, has bool, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribeL4HealthConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Id = &resourceId
	request.Business = &resourceType
	request.RuleIdList = []*string{&ruleId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DescribeL4HealthConfig(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	healthChecks := response.Response.HealthConfig
	length := len(healthChecks)
	if length == 0 {
		return
	}
	if length > 1 {
		errRet = fmt.Errorf("Get L4 health check returns %d healthchecks", length)
		return
	}
	healthCheck = healthChecks[0]
	has = true
	return
}

func (me *DayuService) CreateL4Rule(ctx context.Context, resourceType string, resourceId string, rule dayu.L4RuleEntry) (ruleId string, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Rules = []*dayu.L4RuleEntry{&rule}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().CreateL4Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	//describe rules to get rule ID
	rules, _, has, dErr := me.DescribeL4Rules(ctx, resourceType, resourceId, *rule.RuleName, "")
	if dErr != nil {
		errRet = dErr
		return
	}
	if !has {
		errRet = fmt.Errorf("Create L4 rule failed")
		return
	}
	if len(rules) != 1 {
		errRet = fmt.Errorf("Create L4 rule returns %d rules", len(rules))
		return
	}
	ruleId = *rules[0].RuleId
	return
}

func (me *DayuService) ModifyL4Rule(ctx context.Context, resourceType string, resourceId string, rule dayu.L4RuleEntry) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Rule = &rule
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyL4Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DeleteL4Rule(ctx context.Context, resourceType string, resourceId string, ruleId string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDeleteL4RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.RuleIdList = []*string{&ruleId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().DeleteL4Rules(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) SetSession(ctx context.Context, resourceType string, resourceId string, ruleId string, switchFlag bool, sessionTime int) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyL4KeepTimeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.RuleId = &ruleId
	request.KeepEnable = helper.BoolToInt64Pointer(switchFlag)
	request.KeepTime = helper.IntUint64(sessionTime)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDayuClient().ModifyL4KeepTime(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Success == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if *response.Response.Success.Code != "Success" {
		errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
		return
	}

	return
}

func (me *DayuService) DescribeL7RulesV2(ctx context.Context, business string, offset int, limit int, extendParams map[string]interface{}) (rules []*dayu.NewL7RuleEntry, healths []*dayu.L7RuleHealth, errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDescribleNewL7RulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.Business = &business
	if v, ok := extendParams["protocol"]; ok {
		protocol := v.(string)
		if protocol != "" {
			request.ProtocolList = []*string{&protocol}
		}
	}
	if v, ok := extendParams["domain"]; ok {
		domain := v.(string)
		if domain != "" {
			request.Domain = &domain
		}
	}
	offsetUint64 := uint64(offset)
	request.Offset = &offsetUint64
	limitUint64 := uint64(limit)
	request.Limit = &limitUint64
	var response *dayu.DescribleNewL7RulesResponse
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseDayuClient().DescribleNewL7Rules(request)

		if e, ok := errRet.(*sdkError.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if errRet != nil {
			return resource.RetryableError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}
	rules = response.Response.Rules
	healths = response.Response.Healths
	return

}

func (me *DayuService) CreateL7RuleV2(ctx context.Context, business string, resourceId string, resourceIp string, ruleList []interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewCreateNewL7RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	rule := ruleList[0].(map[string]interface{})
	keeptime := uint64(rule["keeptime"].(int))
	domain := rule["domain"].(string)
	protocol := rule["protocol"].(string)
	sourceType := uint64(rule["source_type"].(int))
	lbType := uint64(rule["lb_type"].(int))
	keepEnable := uint64(rule["keep_enable"].(int))
	certType := uint64(rule["cert_type"].(int))
	sslId := rule["ssl_id"].(string)
	ccEnable := uint64(rule["cc_enable"].(int))
	httpsToHttpEnable := uint64(rule["https_to_http_enable"].(int))

	sourceList := rule["source_list"].([]interface{})
	sources := make([]*dayu.L4RuleSource, 0)
	for _, source := range sourceList {
		sourceItem := source.(map[string]interface{})
		weight := uint64(sourceItem["weight"].(int))
		subSource := sourceItem["source"].(string)
		tmpSource := dayu.L4RuleSource{
			Source: &subSource,
			Weight: &weight,
		}
		sources = append(sources, &tmpSource)
	}

	ruleEntry := dayu.L7RuleEntry{
		KeepTime:          &keeptime,
		Domain:            &domain,
		Protocol:          &protocol,
		SourceType:        &sourceType,
		LbType:            &lbType,
		KeepEnable:        &keepEnable,
		SourceList:        sources,
		CertType:          &certType,
		SSLId:             &sslId,
		CCEnable:          &ccEnable,
		HttpsToHttpEnable: &httpsToHttpEnable,
	}

	request.IdList = []*string{&resourceId}
	request.VipList = []*string{&resourceIp}
	request.Business = &business
	request.Rules = []*dayu.L7RuleEntry{&ruleEntry}
	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseDayuClient().CreateNewL7Rules(request)

		if e, ok := errRet.(*sdkError.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if errRet != nil {
			return resource.RetryableError(errRet)
		}
		if *response.Response.Success.Code != "Success" {
			errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
			return resource.RetryableError(errRet)
		}
		return nil
	})
	return
}

func (me *DayuService) ModifyL7RuleV2(ctx context.Context, resourceType string, resourceId string, rule *dayu.NewL7RuleEntry) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewModifyNewDomainRulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Id = &resourceId
	request.Business = &resourceType
	request.Rule = rule

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseDayuClient().ModifyNewDomainRules(request)

		if e, ok := errRet.(*sdkError.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if errRet != nil {
			return resource.RetryableError(errRet)
		}
		if *response.Response.Success.Code != "Success" {
			errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
			return resource.RetryableError(errRet)
		}
		return nil
	})

	return
}

func (me *DayuService) DeleteL7RulesV2(ctx context.Context, resourceType string, resourceId string, resourceIp string, ruleId string) (errRet error) {
	logId := getLogId(ctx)
	request := dayu.NewDeleteNewL7RulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Business = &resourceType
	request.Rule = []*dayu.L4DelRule{
		{
			Id:         &resourceId,
			Ip:         &resourceIp,
			RuleIdList: []*string{&ruleId},
		},
	}
	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseDayuClient().DeleteNewL7Rules(request)

		if e, ok := errRet.(*sdkError.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if errRet != nil {
			return resource.RetryableError(errRet)
		}
		if *response.Response.Success.Code != "Success" {
			errRet = fmt.Errorf("TencentCloud SDK return %s response,%s", *response.Response.Success.Code, *response.Response.Success.Message)
			return resource.RetryableError(errRet)
		}
		return nil
	})
	return
}
