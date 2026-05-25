package ga2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// describeEndpointGroupsLimit is the API documented maximum page size.
const describeEndpointGroupsLimit = 100

// taskStatusSuccess is the terminal success status returned by DescribeTaskResult.
const taskStatusSuccess = "SUCCESS"

// Ga2Service wraps the GA2 v20250115 SDK client for the provider.
type Ga2Service struct {
	client *connectivity.TencentCloudClient
}

// NewGa2Service constructs a Ga2Service.
func NewGa2Service(client *connectivity.TencentCloudClient) Ga2Service {
	return Ga2Service{client: client}
}

// DescribeGa2EndpointGroupById queries an endpoint group by its three identifying IDs.
// Returns (nil, nil) when the endpoint group does not exist.
func (me *Ga2Service) DescribeGa2EndpointGroupById(ctx context.Context, gaId, listenerId, egId string) (*ga2v20250115.EndpointGroupConfigurationSet, error) {
	logId := tccommon.GetLogId(ctx)

	var (
		offset uint64 = 0
		limit  uint64 = describeEndpointGroupsLimit
		match  *ga2v20250115.EndpointGroupConfigurationSet
	)

	for {
		request := ga2v20250115.NewDescribeEndpointGroupsRequest()
		request.GlobalAcceleratorId = helper.String(gaId)
		request.Offset = &offset
		request.Limit = &limit
		request.Filters = []*ga2v20250115.Filter{
			{
				Name:   helper.String("endpoint-group-id"),
				Values: []*string{helper.String(egId)},
			},
		}

		var response *ga2v20250115.DescribeEndpointGroupsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeEndpointGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 endpoint groups failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 endpoint groups failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.EndpointGroupConfigurationSet
		for i := range set {
			item := set[i]
			if item == nil {
				continue
			}

			// Strict match against the three composite-id components since the API
			// filter values are advisory and may not be enforced server-side.
			if item.EndpointGroupId == nil || item.ListenerId == nil {
				continue
			}

			if *item.EndpointGroupId == egId && *item.ListenerId == listenerId {
				match = item
				return match, nil
			}
		}

		// Stop when the current page is the last page.
		if uint64(len(set)) < limit {
			break
		}

		offset += limit
	}

	return nil, nil
}

// WaitForGa2TaskFinish polls DescribeTaskResult until the task reaches "SUCCESS" or the timeout elapses.
func (me *Ga2Service) WaitForGa2TaskFinish(ctx context.Context, taskId string) error {
	if taskId == "" {
		return fmt.Errorf("ga2 task id is empty, cannot poll task result")
	}

	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
		request := ga2v20250115.NewDescribeTaskResultRequest()
		request.TaskId = helper.String(taskId)

		result, e := me.client.UseGa2V20250115Client().DescribeTaskResultWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ga2 task result failed, Response is nil."))
		}

		status := *result.Response.Status
		if status == taskStatusSuccess {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("ga2 task [%s] is not ready, current status: %s", taskId, status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for ga2 task [%s] failed, reason:%+v", logId, taskId, err)
		return err
	}

	return nil
}
