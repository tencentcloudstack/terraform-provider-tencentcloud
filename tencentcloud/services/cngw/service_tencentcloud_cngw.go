package cngw

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cngwv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cngw/v20230418"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// NewCngwService returns the service wrapper for the CNGW product.
func NewCngwService(client *connectivity.TencentCloudClient) CngwService {
	return CngwService{client: client}
}

// CngwService wraps the CNGW SDK client for the framework reference resources.
type CngwService struct {
	client *connectivity.TencentCloudClient
}

// CreateCloudNativeAPIGatewayAIServiceSource creates an AI service source and
// returns the created source id extracted from the create response.
func (me *CngwService) CreateCloudNativeAPIGatewayAIServiceSource(ctx context.Context, request *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceRequest) (sourceId string, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCngwV20230418Client().CreateCloudNativeAPIGatewayAIServiceSource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create cloud native API gateway AI service source failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cloud native API gateway AI service source failed, reason:%+v", logId, err)
		return
	}

	if response.Response.Result == nil || response.Response.Result.ID == nil {
		err = fmt.Errorf("api %s returns empty source id", request.GetAction())
		return
	}

	sourceId = *response.Response.Result.ID

	return
}

// DescribeCloudNativeAPIGatewayAIServiceSourceList lists the AI service sources
// of a gateway instance with built-in pagination; the caller picks its own item.
func (me *CngwService) DescribeCloudNativeAPIGatewayAIServiceSourceList(ctx context.Context, request *cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListRequest) (sources []*cngwv20230418.CNAPIGwAIServiceSource, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListResponse
		offset   = uint64(0)
		limit    = uint64(100)
	)

	sources = make([]*cngwv20230418.CNAPIGwAIServiceSource, 0)

	for {
		request.Offset = &offset
		request.Limit = &limit
		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCngwV20230418Client().DescribeCloudNativeAPIGatewayAIServiceSourceList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe cloud native API gateway AI service source list failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe cloud native API gateway AI service source list failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		if response.Response.Result == nil || len(response.Response.Result.DataList) < 1 {
			break
		}

		sources = append(sources, response.Response.Result.DataList...)

		if len(response.Response.Result.DataList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

// ModifyCloudNativeAPIGatewayAIServiceSource modifies an AI service source.
func (me *CngwService) ModifyCloudNativeAPIGatewayAIServiceSource(ctx context.Context, request *cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCngwV20230418Client().ModifyCloudNativeAPIGatewayAIServiceSource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify cloud native API gateway AI service source failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify cloud native API gateway AI service source failed, reason:%+v", logId, err)
	}

	return
}

// DeleteCloudNativeAPIGatewayAIServiceSource deletes an AI service source.
func (me *CngwService) DeleteCloudNativeAPIGatewayAIServiceSource(ctx context.Context, request *cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCngwV20230418Client().DeleteCloudNativeAPIGatewayAIServiceSource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete cloud native API gateway AI service source failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cloud native API gateway AI service source failed, reason:%+v", logId, err)
	}

	return
}
