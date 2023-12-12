package antiddos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AntiddosService struct {
	client *connectivity.TencentCloudClient
}

func (me *AntiddosService) DescribeListBGPIPInstances(ctx context.Context, instanceId string, status []string, offset int, limit int) (result []*antiddos.BGPIPInstance, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeListBGPIPInstancesRequest()
	request.FilterInstanceId = &instanceId
	eipType := int64(1)
	request.FilterEipType = &eipType
	filterEipEipAddressStatus := make([]*string, 0)
	for _, singleStatus := range status {
		status := singleStatus
		filterEipEipAddressStatus = append(filterEipEipAddressStatus, &status)
	}
	request.FilterEipEipAddressStatus = filterEipEipAddressStatus
	offsetInt64 := uint64(offset)
	request.Offset = &offsetInt64
	limitInt64 := uint64(limit)
	request.Limit = &limitInt64
	ratelimit.Check(request.GetAction())
	var response *antiddos.DescribeListBGPIPInstancesResponse
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseAntiddosClient().DescribeListBGPIPInstances(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "tccommon.InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	result = response.Response.InstanceList
	return
}

func (me *AntiddosService) AssociateDDoSEipAddress(ctx context.Context, instanceId string, eip string, cvmInstanceID string, cvmRegion string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewAssociateDDoSEipAddressRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)
	request.CvmInstanceID = common.StringPtr(cvmInstanceID)
	request.CvmRegion = common.StringPtr(cvmRegion)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().AssociateDDoSEipAddress(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "tccommon.InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return

}

func (me *AntiddosService) AssociateDDoSEipLoadBalancer(ctx context.Context, instanceId string, eip string, loadBalancerID string, loadBalancerRegion string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewAssociateDDoSEipLoadBalancerRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)
	request.LoadBalancerID = common.StringPtr(loadBalancerID)
	request.LoadBalancerRegion = common.StringPtr(loadBalancerRegion)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().AssociateDDoSEipLoadBalancer(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "tccommon.InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return

}

func (me *AntiddosService) DisassociateDDoSEipAddress(ctx context.Context, instanceId string, eip string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDisassociateDDoSEipAddressRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DisassociateDDoSEipAddress(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "tccommon.InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return

}

func (me *AntiddosService) DescribeListProtectThresholdConfig(ctx context.Context, instanceId string) (result antiddos.ProtectThresholdRelation, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeListProtectThresholdConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	request.Limit = helper.IntUint64(1)
	request.Offset = helper.Int64Uint64(0)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseAntiddosClient().DescribeListProtectThresholdConfig(request)
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = *configList[0]
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeListBlackWhiteIpList(ctx context.Context, instanceId string) (result []*antiddos.BlackWhiteIpRelation, err error) {
	request := antiddos.NewDescribeListBlackWhiteIpListRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	var limit int64 = 10
	var offset int64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListBlackWhiteIpList(request)
		if e != nil {
			err = e
			return
		}
		ipList := response.Response.IpList
		if len(ipList) > 0 {
			result = append(result, ipList...)
		}
		if len(ipList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DescribeListPortAclList(ctx context.Context, instanceId string) (result []*antiddos.AclConfigRelation, err error) {
	request := antiddos.NewDescribeListPortAclListRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	var limit uint64 = 10
	var offset uint64 = 0

	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListPortAclList(request)
		if e != nil {
			err = e
			return
		}
		aclList := response.Response.AclList
		if len(aclList) > 0 {
			result = append(result, aclList...)
		}
		if len(aclList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DescribeListProtocolBlockConfig(ctx context.Context, instanceId string) (result antiddos.ProtocolBlockRelation, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeListProtocolBlockConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	request.Limit = helper.IntInt64(1)
	request.Offset = helper.IntInt64(0)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseAntiddosClient().DescribeListProtocolBlockConfig(request)
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = *configList[0]
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeDDoSConnectLimitList(ctx context.Context, instanceId string) (result antiddos.ConnectLimitConfig, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeDDoSConnectLimitListRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	request.Limit = helper.IntUint64(1)
	request.Offset = helper.IntUint64(0)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseAntiddosClient().DescribeDDoSConnectLimitList(request)
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = *configList[0].ConnectLimitConfig
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeListDDoSAI(ctx context.Context, instanceId string) (result antiddos.DDoSAIRelation, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeListDDoSAIRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	request.Limit = helper.IntInt64(1)
	request.Offset = helper.IntInt64(0)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseAntiddosClient().DescribeListDDoSAI(request)
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = *configList[0]
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeListDDoSGeoIPBlockConfig(ctx context.Context, instanceId string) (result []*antiddos.DDoSGeoIPBlockConfigRelation, err error) {
	request := antiddos.NewDescribeListDDoSGeoIPBlockConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListDDoSGeoIPBlockConfig(request)
		if e != nil {
			err = e
			return
		}
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = append(result, configList...)
		}
		if len(configList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DescribeListDDoSSpeedLimitConfig(ctx context.Context, instanceId string) (result []*antiddos.DDoSSpeedLimitConfigRelation, err error) {
	request := antiddos.NewDescribeListDDoSSpeedLimitConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListDDoSSpeedLimitConfig(request)
		if e != nil {
			err = e
			return
		}
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = append(result, configList...)
		}
		if len(configList) < int(limit) {
			return
		}
		offset += limit
	}

}

func (me *AntiddosService) DescribeListPacketFilterConfig(ctx context.Context, instanceId string) (result []*antiddos.PacketFilterRelation, err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDescribeListPacketFilterConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	request.Limit = helper.IntInt64(1)
	request.Offset = helper.IntInt64(0)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err := me.client.UseAntiddosClient().DescribeListPacketFilterConfig(request)
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = configList
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

// func (me *AntiddosService) DescribeListWaterPrintConfig(ctx context.Context, instanceId string) (result antiddos.WaterPrintRelation, err error) {
// 	logId := tccommon.GetLogId(ctx)
// 	request := antiddos.NewDescribeListWaterPrintConfigRequest()
// 	request.FilterInstanceId = common.StringPtr(instanceId)
// 	request.Limit = helper.IntInt64(1)
// 	request.Offset = helper.IntInt64(0)

// 	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		response, err := me.client.UseAntiddosClient().DescribeListWaterPrintConfig(request)
// 		configList := response.Response.ConfigList
// 		if len(configList) > 0 {
// 			result = *configList[0]
// 		}
// 		if err != nil {
// 			return resource.RetryableError(err)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
// 			logId, request.GetAction(), request.ToJsonString(), err.Error())
// 		return
// 	}
// 	return
// }

func (me *AntiddosService) CreateDDoSBlackWhiteIpList(ctx context.Context, instanceId string, ipList []string, ipType string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSBlackWhiteIpListRequest()
	request.InstanceId = common.StringPtr(instanceId)
	requestIpList := make([]*antiddos.IpSegment, 0)
	ip32mask := uint64(0)
	for _, ip := range ipList {
		requestIpList = append(requestIpList, &antiddos.IpSegment{Ip: &ip, Mask: &ip32mask})
	}
	request.IpList = requestIpList
	request.Type = common.StringPtr(ipType)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSBlackWhiteIpList(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) ModifyDDoSThreshold(ctx context.Context, business, instanceId string, threshold int) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyDDoSThresholdRequest()
	request.Business = common.StringPtr(business)
	request.Id = common.StringPtr(instanceId)
	request.Threshold = helper.IntUint64(threshold)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyDDoSThreshold(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceUnavailable" {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) ModifyDDoSLevel(ctx context.Context, business, instanceId, ddosLevel string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyDDoSLevelRequest()
	request.Business = common.StringPtr(business)
	request.Id = common.StringPtr(instanceId)
	request.Method = common.StringPtr("set")
	request.DDoSLevel = common.StringPtr(ddosLevel)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyDDoSLevel(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreatePortAclConfig(ctx context.Context, instanceId string, aclConfig antiddos.AclConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreatePortAclConfigRequest()
	request.AclConfig = &aclConfig
	request.InstanceId = &instanceId

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreatePortAclConfig(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceInUse" {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateProtocolBlockConfig(ctx context.Context, instanceId string, protocolBlockConfig antiddos.ProtocolBlockConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateProtocolBlockConfigRequest()
	request.InstanceId = &instanceId
	request.ProtocolBlockConfig = &protocolBlockConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateProtocolBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateWaterPrintConfig(ctx context.Context, instanceId string, waterPrintConfig antiddos.WaterPrintConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateWaterPrintConfigRequest()
	request.InstanceId = &instanceId
	request.WaterPrintConfig = &waterPrintConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateWaterPrintConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeListWaterPrintConfig(ctx context.Context, instanceId string) (result []*antiddos.WaterPrintRelation, err error) {
	request := antiddos.NewDescribeListWaterPrintConfigRequest()
	request.FilterInstanceId = common.StringPtr(instanceId)
	var limit int64 = 10
	var offset int64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListWaterPrintConfig(request)
		if e != nil {
			err = e
			return
		}
		configList := response.Response.ConfigList
		if len(configList) > 0 {
			result = append(result, configList...)
		}
		if len(configList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DeleteWaterPrintConfig(ctx context.Context, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteWaterPrintConfigRequest()
	request.InstanceId = &instanceId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteWaterPrintConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) SwitchWaterPrintConfig(ctx context.Context, instanceId string, openStatus int) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewSwitchWaterPrintConfigRequest()
	request.InstanceId = &instanceId
	request.OpenStatus = helper.IntInt64(openStatus)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().SwitchWaterPrintConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateDDoSConnectLimit(ctx context.Context, instanceId string, connectLimitConfig antiddos.ConnectLimitConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSConnectLimitRequest()
	request.InstanceId = &instanceId
	request.ConnectLimitConfig = &connectLimitConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSConnectLimit(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateDDoSAI(ctx context.Context, instanceId, ddosAI string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSAIRequest()
	request.DDoSAI = &ddosAI
	request.InstanceIdList = []*string{&instanceId}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSAI(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateDDoSGeoIPBlockConfig(ctx context.Context, instanceId string, ddosGeoIPBlockConfig antiddos.DDoSGeoIPBlockConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSGeoIPBlockConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSGeoIPBlockConfig = &ddosGeoIPBlockConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSGeoIPBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreateDDoSSpeedLimitConfig(ctx context.Context, instanceId string, ddosSpeedLimitConfig antiddos.DDoSSpeedLimitConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSSpeedLimitConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSSpeedLimitConfig = &ddosSpeedLimitConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSSpeedLimitConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) CreatePacketFilterConfig(ctx context.Context, instanceId string, packetFilterConfig antiddos.PacketFilterConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreatePacketFilterConfigRequest()
	request.InstanceId = &instanceId
	request.PacketFilterConfig = &packetFilterConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreatePacketFilterConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSBlackWhiteIpList(ctx context.Context, instanceId string, ips []string, ipType string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteDDoSBlackWhiteIpListRequest()
	request.InstanceId = &instanceId
	request.Type = &ipType
	ipList := make([]*antiddos.IpSegment, 0)
	for _, ip := range ips {
		ipList = append(ipList, &antiddos.IpSegment{
			Ip:   &ip,
			Mask: helper.IntUint64(0),
		})
	}
	request.IpList = ipList
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteDDoSBlackWhiteIpList(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeletePortAclConfig(ctx context.Context, instanceId string, aclConfig antiddos.AclConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeletePortAclConfigRequest()
	request.InstanceId = &instanceId
	request.AclConfig = &aclConfig
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeletePortAclConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteProtocolBlockConfig(ctx context.Context, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateProtocolBlockConfigRequest()
	request.InstanceId = &instanceId
	protocolBlockConfig := antiddos.ProtocolBlockConfig{
		DropIcmp:               helper.IntInt64(0),
		DropTcp:                helper.IntInt64(0),
		DropUdp:                helper.IntInt64(0),
		DropOther:              helper.IntInt64(0),
		CheckExceptNullConnect: helper.IntInt64(1),
	}
	request.ProtocolBlockConfig = &protocolBlockConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateProtocolBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSConnectLimit(ctx context.Context, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSConnectLimitRequest()
	request.InstanceId = &instanceId
	connectLimitConfig := antiddos.ConnectLimitConfig{
		SdNewLimit:       helper.IntUint64(0),
		DstNewLimit:      helper.IntUint64(0),
		SdConnLimit:      helper.IntUint64(0),
		DstConnLimit:     helper.IntUint64(0),
		BadConnThreshold: helper.IntUint64(0),
		NullConnEnable:   helper.IntUint64(0),
		ConnTimeout:      helper.IntUint64(0),
		SynRate:          helper.IntUint64(0),
		SynLimit:         helper.IntUint64(0),
	}
	request.ConnectLimitConfig = &connectLimitConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSConnectLimit(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSAI(ctx context.Context, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateDDoSAIRequest()
	request.DDoSAI = common.StringPtr("off")
	request.InstanceIdList = []*string{&instanceId}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateDDoSAI(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSGeoIPBlockConfig(ctx context.Context, instanceId string, ddosGeoIPBlockConfig antiddos.DDoSGeoIPBlockConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteDDoSGeoIPBlockConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSGeoIPBlockConfig = &ddosGeoIPBlockConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteDDoSGeoIPBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSSpeedLimitConfig(ctx context.Context, instanceId string, ddosSpeedLimitConfig antiddos.DDoSSpeedLimitConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteDDoSSpeedLimitConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSSpeedLimitConfig = &ddosSpeedLimitConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteDDoSSpeedLimitConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeletePacketFilterConfig(ctx context.Context, instanceId string, packetFilterConfig antiddos.PacketFilterConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeletePacketFilterConfigRequest()
	request.InstanceId = &instanceId
	request.PacketFilterConfig = &packetFilterConfig

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeletePacketFilterConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteDDoSThreshold(ctx context.Context, business, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyDDoSThresholdRequest()
	request.Business = common.StringPtr(business)
	request.Id = common.StringPtr(instanceId)
	request.Threshold = helper.IntUint64(0)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyDDoSThreshold(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceUnavailable" {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteyDDoSLevel(ctx context.Context, business, instanceId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyDDoSLevelRequest()
	request.Business = common.StringPtr(business)
	request.Id = common.StringPtr(instanceId)
	request.Method = common.StringPtr("set")
	request.DDoSLevel = common.StringPtr("middle")

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyDDoSLevel(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCCThresholdList(ctx context.Context, business, instanceId string) (result []*antiddos.CCThresholdPolicy, err error) {
	request := antiddos.NewDescribeCCThresholdListRequest()
	request.Business = &business
	request.InstanceId = &instanceId
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCCThresholdList(request)
		if e != nil {
			err = e
			return
		}
		thresholdList := response.Response.ThresholdList
		if len(thresholdList) > 0 {
			result = append(result, thresholdList...)
		}
		if len(thresholdList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) ModifyCCThresholdPolicy(ctx context.Context, instanceId, protocol, ip, domain string, threshold int) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyCCThresholdPolicyRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Protocol = &protocol
	request.Threshold = helper.IntInt64(threshold)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyCCThresholdPolicy(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceUnavailable" {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCcGeoIPBlockConfigList(ctx context.Context, business, instanceId string) (result []*antiddos.CcGeoIpPolicyNew, err error) {
	request := antiddos.NewDescribeCcGeoIPBlockConfigListRequest()
	request.Business = &business
	request.InstanceId = &instanceId
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCcGeoIPBlockConfigList(request)
		if e != nil {
			err = e
			return
		}
		ccGeoIpPolicyNew := response.Response.CcGeoIpPolicyList
		if len(ccGeoIpPolicyNew) > 0 {
			result = append(result, ccGeoIpPolicyNew...)
		}
		if len(ccGeoIpPolicyNew) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) CreateCcGeoIPBlockConfig(ctx context.Context, instanceId, protocol, ip, domain string, ccGeoIPBlockConfig antiddos.CcGeoIPBlockConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateCcGeoIPBlockConfigRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.IP = &ip
	request.Protocol = &protocol
	request.CcGeoIPBlockConfig = &ccGeoIPBlockConfig
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateCcGeoIPBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteCcGeoIPBlockConfig(ctx context.Context, instanceId string, ccGeoIPBlockConfig antiddos.CcGeoIPBlockConfig) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCcGeoIPBlockConfigRequest()
	request.InstanceId = &instanceId
	request.CcGeoIPBlockConfig = &ccGeoIPBlockConfig
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCcGeoIPBlockConfig(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCcBlackWhiteIpList(ctx context.Context, business, instanceId string) (result []*antiddos.CcBlackWhiteIpPolicy, err error) {
	request := antiddos.NewDescribeCcBlackWhiteIpListRequest()
	request.Business = &business
	request.InstanceId = &instanceId
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCcBlackWhiteIpList(request)
		if e != nil {
			err = e
			return
		}
		ccBlackWhiteIpList := response.Response.CcBlackWhiteIpList
		if len(ccBlackWhiteIpList) > 0 {
			result = append(result, ccBlackWhiteIpList...)
		}
		if len(ccBlackWhiteIpList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) CreateCcBlackWhiteIpList(ctx context.Context, instanceId, protocol, ip, domain, ipType string, posIps []string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateCcBlackWhiteIpListRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Protocol = &protocol
	request.Type = &ipType
	ipLists := make([]*antiddos.IpSegment, 0)
	for _, posIp := range posIps {
		ipLists = append(ipLists, &antiddos.IpSegment{
			Ip:   &posIp,
			Mask: helper.IntUint64(0),
		})
	}
	request.IpList = ipLists
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateCcBlackWhiteIpList(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteCcBlackWhiteIpList(ctx context.Context, instanceId, policyId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCcBlackWhiteIpListRequest()
	request.InstanceId = &instanceId
	request.PolicyId = &policyId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCcBlackWhiteIpList(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCCPrecisionPlyList(ctx context.Context, business, instanceId string) (result []*antiddos.CCPrecisionPolicy, err error) {
	request := antiddos.NewDescribeCCPrecisionPlyListRequest()
	request.Business = &business
	request.InstanceId = &instanceId
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCCPrecisionPlyList(request)
		if e != nil {
			err = e
			return
		}
		precisionPolicyList := response.Response.PrecisionPolicyList
		if len(precisionPolicyList) > 0 {
			result = append(result, precisionPolicyList...)
		}
		if len(precisionPolicyList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) CreateCCPrecisionPolicy(ctx context.Context, instanceId, protocol, ip, domain, policyAction string, policyList []*antiddos.CCPrecisionPlyRecord) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateCCPrecisionPolicyRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Protocol = &protocol
	request.PolicyAction = &policyAction
	request.PolicyList = policyList
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateCCPrecisionPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteCCPrecisionPolicy(ctx context.Context, instanceId, policyId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCCPrecisionPolicyRequest()
	request.InstanceId = &instanceId
	request.PolicyId = &policyId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCCPrecisionPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) ModifyCCLevelPolicy(ctx context.Context, instanceId, ip, domain, protocol, level string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewModifyCCLevelPolicyRequest()
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Domain = &domain
	request.Protocol = &protocol
	request.Level = &level
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().ModifyCCLevelPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCCReqLimitPolicyList(ctx context.Context, business, instanceId string) (result []*antiddos.CCReqLimitPolicy, err error) {
	request := antiddos.NewDescribeCCReqLimitPolicyListRequest()
	request.Business = &business
	request.InstanceId = &instanceId
	var limit uint64 = 10
	var offset uint64 = 0
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCCReqLimitPolicyList(request)
		if e != nil {
			err = e
			return
		}
		requestLimitPolicyList := response.Response.RequestLimitPolicyList
		if len(requestLimitPolicyList) > 0 {
			result = append(result, requestLimitPolicyList...)
		}
		if len(requestLimitPolicyList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) CreateCCReqLimitPolicy(ctx context.Context, instanceId, protocol, ip, domain string, ccReqLimitPolicyRecord antiddos.CCReqLimitPolicyRecord) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewCreateCCReqLimitPolicyRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Protocol = &protocol
	request.Policy = &ccReqLimitPolicyRecord
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().CreateCCReqLimitPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteCCRequestLimitPolicy(ctx context.Context, instanceId, policyId string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCCRequestLimitPolicyRequest()
	request.InstanceId = &instanceId
	request.PolicyId = &policyId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCCRequestLimitPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeCCLevelPolicy(ctx context.Context, domain, instanceId, ip, protocol string) (level string, err error) {
	request := antiddos.NewDescribeCCLevelPolicyRequest()
	request.Domain = &domain
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Protocol = &protocol

	ratelimit.Check(request.GetAction())
	response, e := me.client.UseAntiddosClient().DescribeCCLevelPolicy(request)
	if e != nil {
		err = e
		return
	}
	level = *response.Response.Level
	return
}

func (me *AntiddosService) DescribeListBGPIPInstanceById(ctx context.Context, business, instanceId string) (result []*antiddos.BGPIPInstance, err error) {

	var limit uint64 = 10
	var offset uint64 = 0
	request := antiddos.NewDescribeListBGPIPInstancesRequest()
	request.Limit = &limit
	request.Offset = &offset
	request.FilterInstanceId = &instanceId

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListBGPIPInstances(request)
		if e != nil {
			err = e
			return
		}
		bgpipInstanceList := response.Response.InstanceList
		if len(bgpipInstanceList) > 0 {
			result = append(result, bgpipInstanceList...)
		}
		if len(bgpipInstanceList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DescribeListBGPInstanceById(ctx context.Context, business, instanceId string) (result []*antiddos.BGPInstance, err error) {

	var limit uint64 = 10
	var offset uint64 = 0
	request := antiddos.NewDescribeListBGPInstancesRequest()
	request.Limit = &limit
	request.Offset = &offset
	request.FilterInstanceId = &instanceId

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeListBGPInstances(request)
		if e != nil {
			err = e
			return
		}
		bgpipInstanceList := response.Response.InstanceList
		if len(bgpipInstanceList) > 0 {
			result = append(result, bgpipInstanceList...)
		}
		if len(bgpipInstanceList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DescribeCCLevelList(ctx context.Context, business, instanceId string) (result []*antiddos.CCLevelPolicy, err error) {

	var limit uint64 = 10
	var offset uint64 = 0
	request := antiddos.NewDescribeCCLevelListRequest()
	request.Limit = &limit
	request.Offset = &offset
	request.InstanceId = &instanceId
	request.Business = &business

	for {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAntiddosClient().DescribeCCLevelList(request)
		if e != nil {
			err = e
			return
		}
		bgpipInstanceList := response.Response.LevelList
		if len(bgpipInstanceList) > 0 {
			result = append(result, bgpipInstanceList...)
		}
		if len(bgpipInstanceList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *AntiddosService) DeleteCCLevelPolicy(ctx context.Context, instanceId, ip, domain string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCCLevelPolicyRequest()
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Domain = &domain
	request.Protocol = common.StringPtr("http")
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCCLevelPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DeleteCCThresholdPolicy(ctx context.Context, instanceId, ip, domain string) (err error) {
	logId := tccommon.GetLogId(ctx)
	request := antiddos.NewDeleteCCThresholdPolicyRequest()
	request.InstanceId = &instanceId
	request.Ip = &ip
	request.Domain = &domain
	request.Protocol = common.StringPtr("http")
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DeleteCCThresholdPolicy(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeAntiddosBoundipById(ctx context.Context, id string) (boundip *antiddos.BGPInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListBGPInstancesRequest()
	request.FilterInstanceId = &id
	request.Limit = helper.IntUint64(10)
	request.Offset = helper.IntUint64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeListBGPInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceList) < 1 {
		return
	}

	boundip = response.Response.InstanceList[0]
	return
}

func (me *AntiddosService) DescribeAntiddosPendingRiskInfoByFilter(ctx context.Context) (pendingRiskInfoResponseParams *antiddos.DescribePendingRiskInfoResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribePendingRiskInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribePendingRiskInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		pendingRiskInfoResponseParams = response.Response
	}

	return
}

func (me *AntiddosService) DescribeAntiddosOverviewIndexByFilter(ctx context.Context, param map[string]interface{}) (describeOverviewIndexResponseParams *antiddos.DescribeOverviewIndexResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeOverviewIndexRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeOverviewIndex(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	describeOverviewIndexResponseParams = response.Response
	return
}

func (me *AntiddosService) DescribeAntiddosOverviewDdosTrendByFilter(ctx context.Context, param map[string]interface{}) (describeOverviewDDoSTrendResponseParams *antiddos.DescribeOverviewDDoSTrendResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeOverviewDDoSTrendRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Period" {
			request.Period = v.(*int64)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "MetricName" {
			request.MetricName = v.(*string)
		}
		if k == "Business" {
			request.Business = v.(*string)
		}
		if k == "IpList" {
			request.IpList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeOverviewDDoSTrend(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		describeOverviewDDoSTrendResponseParams = response.Response
	}

	return
}

func (me *AntiddosService) DescribeAntiddosOverviewDdosEventListByFilter(ctx context.Context, param map[string]interface{}) (overviewDdosEventList []*antiddos.OverviewDDoSEvent, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeOverviewDDoSEventListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "AttackStatus" {
			request.AttackStatus = v.(*string)
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
		response, err := me.client.UseAntiddosClient().DescribeOverviewDDoSEventList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EventList) < 1 {
			break
		}
		overviewDdosEventList = append(overviewDdosEventList, response.Response.EventList...)
		if len(response.Response.EventList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DescribeAntiddosOverviewCcTrendByFilter(ctx context.Context, param map[string]interface{}) (overviewCCTrendResponseParams *antiddos.DescribeOverviewCCTrendResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeOverviewCCTrendRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Period" {
			request.Period = v.(*int64)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "MetricName" {
			request.MetricName = v.(*string)
		}
		if k == "Business" {
			request.Business = v.(*string)
		}
		if k == "IpList" {
			request.IpList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeOverviewCCTrend(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		overviewCCTrendResponseParams = response.Response
	}

	return
}

func (me *AntiddosService) DescribeAntiddosDdosBlackWhiteIpListById(ctx context.Context, instanceId string) (ddosBlackWhiteIpListResponseParams *antiddos.DescribeDDoSBlackWhiteIpListResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeDDoSBlackWhiteIpListRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeDDoSBlackWhiteIpList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ddosBlackWhiteIpListResponseParams = response.Response

	return
}

func (me *AntiddosService) DeleteAntiddosDdosBlackWhiteIpListById(ctx context.Context, params map[string]interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeleteDDoSBlackWhiteIpListRequest()
	if v, ok := params["instanceId"]; ok {
		request.InstanceId = helper.String(v.(string))
	}
	if v, ok := params["ipType"]; ok {
		request.Type = helper.String(v.(string))
	}
	ipSegment := antiddos.IpSegment{}
	if v, ok := params["ip"]; ok {
		ipSegment.Ip = helper.String(v.(string))
	}
	if v, ok := params["ipMask"]; ok {
		ipSegment.Mask = helper.IntUint64(v.(int))
	}
	request.IpList = []*antiddos.IpSegment{&ipSegment}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeleteDDoSBlackWhiteIpList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosBasicDeviceStatusByFilter(ctx context.Context, param map[string]interface{}) (basicDeviceStatus *antiddos.DescribeBasicDeviceStatusResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeBasicDeviceStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "IpList" {
			request.IpList = v.([]*string)
		}
		if k == "IdList" {
			request.IdList = v.([]*string)
		}
		if k == "FilterRegion" {
			request.FilterRegion = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeBasicDeviceStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	basicDeviceStatus = response.Response
	return
}

func (me *AntiddosService) DescribeAntiddosBgpBizTrendByFilter(ctx context.Context, param map[string]interface{}) (bgpBizTrend *antiddos.DescribeBgpBizTrendResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeBgpBizTrendRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Business" {
			request.Business = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "MetricName" {
			request.MetricName = v.(*string)
		}
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Flag" {
			request.Flag = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeBgpBizTrend(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	bgpBizTrend = response.Response
	return
}

func (me *AntiddosService) DescribeAntiddosListListenerByFilter(ctx context.Context) (listListener *antiddos.DescribeListListenerResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeListListenerRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeListListener(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	listListener = response.Response
	return
}

func (me *AntiddosService) DescribeAntiddosOverviewAttackTrendByFilter(ctx context.Context, param map[string]interface{}) (overviewAttackTrend *antiddos.DescribeOverviewAttackTrendResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = antiddos.NewDescribeOverviewAttackTrendRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Dimension" {
			request.Dimension = v.(*string)
		}
		if k == "Period" {
			request.Period = v.(*uint64)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeOverviewAttackTrend(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	overviewAttackTrend = response.Response
	return
}

func (me *AntiddosService) DescribeAntiddosDdosGeoIpBlockConfigById(ctx context.Context, instanceId string) (configList []*antiddos.DDoSGeoIPBlockConfigRelation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListDDoSGeoIPBlockConfigRequest()
	request.FilterInstanceId = &instanceId

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
		response, err := me.client.UseAntiddosClient().DescribeListDDoSGeoIPBlockConfig(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ConfigList) < 1 {
			break
		}
		configList = append(configList, response.Response.ConfigList...)
		if len(response.Response.ConfigList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DeleteAntiddosDdosGeoIpBlockConfigById(ctx context.Context, instanceId string, config *antiddos.DDoSGeoIPBlockConfig) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeleteDDoSGeoIPBlockConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSGeoIPBlockConfig = config

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeleteDDoSGeoIPBlockConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosDdosSpeedLimitConfigById(ctx context.Context, instanceId string) (configList []*antiddos.DDoSSpeedLimitConfigRelation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListDDoSSpeedLimitConfigRequest()
	request.FilterInstanceId = &instanceId

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
		response, err := me.client.UseAntiddosClient().DescribeListDDoSSpeedLimitConfig(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ConfigList) < 1 {
			break
		}
		configList = append(configList, response.Response.ConfigList...)
		if len(response.Response.ConfigList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DeleteAntiddosDdosSpeedLimitConfigById(ctx context.Context, instanceId string, config *antiddos.DDoSSpeedLimitConfig) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeleteDDoSSpeedLimitConfigRequest()
	request.InstanceId = &instanceId
	request.DDoSSpeedLimitConfig = config

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeleteDDoSSpeedLimitConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosDefaultAlarmThresholdById(ctx context.Context, instanceType string, filterAlarmType int64) (defaultAlarmThreshold *antiddos.DefaultAlarmThreshold, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeDefaultAlarmThresholdRequest()
	request.InstanceType = &instanceType
	request.FilterAlarmType = &filterAlarmType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeDefaultAlarmThreshold(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DefaultAlarmConfigList) < 1 {
		return
	}

	defaultAlarmThreshold = response.Response.DefaultAlarmConfigList[0]
	return
}

func (me *AntiddosService) DescribeAntiddosSchedulingDomainUserNameById(ctx context.Context, domainName string) (schedulingDomainUserName *antiddos.SchedulingDomainInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListSchedulingDomainRequest()
	request.FilterDomain = &domainName
	request.Offset = helper.Uint64(0)
	request.Limit = helper.Uint64(10)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeListSchedulingDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainList) < 1 {
		return
	}

	schedulingDomainUserName = response.Response.DomainList[0]
	return
}

func (me *AntiddosService) DescribeAntiddosIpAlarmThresholdConfigById(ctx context.Context, instanceId, eip string, alarmType int) (ipAlarmThresholdConfig *antiddos.IPAlarmThresholdRelation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListIPAlarmConfigRequest()
	request.FilterInstanceId = &instanceId
	request.FilterIp = &eip
	request.FilterAlarmType = helper.IntInt64(alarmType)
	request.Limit = helper.Int64(10)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DescribeListIPAlarmConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ConfigList) < 1 {
		return
	}

	ipAlarmThresholdConfig = response.Response.ConfigList[0]
	return
}

func (me *AntiddosService) DescribeAntiddosPacketFilterConfigById(ctx context.Context, instanceId string) (configList []*antiddos.PacketFilterRelation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListPacketFilterConfigRequest()
	request.FilterInstanceId = &instanceId

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
		response, err := me.client.UseAntiddosClient().DescribeListPacketFilterConfig(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ConfigList) < 1 {
			break
		}
		configList = append(configList, response.Response.ConfigList...)
		if len(response.Response.ConfigList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DeleteAntiddosPacketFilterConfigById(ctx context.Context, instanceId string, config *antiddos.PacketFilterConfig) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeletePacketFilterConfigRequest()
	request.InstanceId = &instanceId
	request.PacketFilterConfig = config

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeletePacketFilterConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosPortAclConfigById(ctx context.Context, instanceId string) (portAclConfig []*antiddos.AclConfigRelation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeListPortAclListRequest()
	request.FilterInstanceId = &instanceId

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
		response, err := me.client.UseAntiddosClient().DescribeListPortAclList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AclList) < 1 {
			break
		}
		portAclConfig = append(portAclConfig, response.Response.AclList...)
		if len(response.Response.AclList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DeleteAntiddosPortAclConfigById(ctx context.Context, instanceId string, config *antiddos.AclConfig) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeletePortAclConfigRequest()
	request.InstanceId = &instanceId
	request.AclConfig = config

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeletePortAclConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosCcBlackWhiteIpById(ctx context.Context, business, instanceId, ip, domain, protocol string) (ccBlackWhiteIps []*antiddos.CcBlackWhiteIpPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeCcBlackWhiteIpListRequest()
	request.InstanceId = &instanceId
	request.Business = &business
	request.Ip = &ip
	request.Domain = &domain
	request.Protocol = &protocol

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
		response, err := me.client.UseAntiddosClient().DescribeCcBlackWhiteIpList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CcBlackWhiteIpList) < 1 {
			break
		}
		ccBlackWhiteIps = append(ccBlackWhiteIps, response.Response.CcBlackWhiteIpList...)
		if len(response.Response.CcBlackWhiteIpList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *AntiddosService) DeleteAntiddosCcBlackWhiteIpById(ctx context.Context, instanceId, policyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeleteCcBlackWhiteIpListRequest()
	request.InstanceId = &instanceId
	request.PolicyId = &policyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeleteCcBlackWhiteIpList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *AntiddosService) DescribeAntiddosCcPrecisionPolicyById(ctx context.Context, business, instanceId, ip, domain, protocol string) (ccPrecisionPolicys []*antiddos.CCPrecisionPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDescribeCCPrecisionPlyListRequest()
	request.InstanceId = &instanceId
	request.Business = &business
	request.Ip = &ip
	request.Domain = &domain
	request.Protocol = &protocol

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
		response, err := me.client.UseAntiddosClient().DescribeCCPrecisionPlyList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.PrecisionPolicyList) < 1 {
			break
		}
		ccPrecisionPolicys = append(ccPrecisionPolicys, response.Response.PrecisionPolicyList...)
		if len(response.Response.PrecisionPolicyList) < int(limit) {
			break
		}

		offset += limit
	}
	return
}

func (me *AntiddosService) DeleteAntiddosCcPrecisionPolicyById(ctx context.Context, instanceId, policyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := antiddos.NewDeleteCCPrecisionPolicyRequest()
	request.InstanceId = &instanceId
	request.PolicyId = &policyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseAntiddosClient().DeleteCCPrecisionPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
