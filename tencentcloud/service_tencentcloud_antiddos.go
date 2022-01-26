package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AntiddosService struct {
	client *connectivity.TencentCloudClient
}

func (me *AntiddosService) DescribeListBGPIPInstances(ctx context.Context, instanceId string, status []string, offset int, limit int) (result []*antiddos.BGPIPInstance, err error) {
	logId := getLogId(ctx)
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
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseAntiddosClient().DescribeListBGPIPInstances(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
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
	logId := getLogId(ctx)
	request := antiddos.NewAssociateDDoSEipAddressRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)
	request.CvmInstanceID = common.StringPtr(cvmInstanceID)
	request.CvmRegion = common.StringPtr(cvmRegion)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().AssociateDDoSEipAddress(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
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
	logId := getLogId(ctx)
	request := antiddos.NewAssociateDDoSEipLoadBalancerRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)
	request.LoadBalancerID = common.StringPtr(loadBalancerID)
	request.LoadBalancerRegion = common.StringPtr(loadBalancerRegion)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().AssociateDDoSEipLoadBalancer(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
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
	logId := getLogId(ctx)
	request := antiddos.NewDisassociateDDoSEipAddressRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.Eip = common.StringPtr(eip)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseAntiddosClient().DisassociateDDoSEipAddress(request)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
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
