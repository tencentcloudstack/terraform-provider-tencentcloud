package alb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// NewAlbService builds the ALB service wrapper used by the framework
// references of the alb product.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps every TencentCloud ALB API call used by the framework
// references of the alb product.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateListener calls the synchronous CreateListener API and returns the
// created listener id.
func (me *AlbService) CreateListener(ctx context.Context, request *albv20251030.CreateListenerRequest) (listenerId string, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.CreateListenerResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create listener failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create listener failed, reason:%+v", logId, err)
		return
	}

	if response.Response.ListenerId == nil {
		err = fmt.Errorf("api %s returns empty listener id", request.GetAction())
		return
	}

	listenerId = *response.Response.ListenerId

	return
}

// DescribeListeners calls DescribeListeners and returns the full listener
// list. The API is token-paginated (MaxResults + NextToken), so the paging
// loop is built in: a unique-id filter terminates after the first round, a
// full-list query keeps paging until the cloud stops handing out tokens.
func (me *AlbService) DescribeListeners(ctx context.Context, request *albv20251030.DescribeListenersRequest) (listeners []*albv20251030.ListenerOutput, err error) {
	var (
		logId     = tccommon.GetLogId(ctx)
		nextToken string
		limit     = uint64(100) // the maximum MaxResults documented by the API
	)

	listeners = make([]*albv20251030.ListenerOutput, 0)

	for {
		request.MaxResults = &limit
		request.NextToken = &nextToken
		var response *albv20251030.DescribeListenersResponse
		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeListeners(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe listeners failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe listeners failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		if len(response.Response.Listeners) < 1 {
			break
		}

		listeners = append(listeners, response.Response.Listeners...)

		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}

		nextToken = *response.Response.NextToken
	}

	return
}

// ModifyListenerAttributes calls the synchronous ModifyListenerAttributes API.
func (me *AlbService) ModifyListenerAttributes(ctx context.Context, request *albv20251030.ModifyListenerAttributesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyListenerAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify listener attributes failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify listener attributes failed, reason:%+v", logId, err)
	}

	return
}

// DeleteListener calls the synchronous DeleteListener API.
func (me *AlbService) DeleteListener(ctx context.Context, request *albv20251030.DeleteListenerRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete listener failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete listener failed, reason:%+v", logId, err)
	}

	return
}
