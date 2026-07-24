package ga2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

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

	request := ga2v20250115.NewDescribeEndpointGroupsRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.Filters = []*ga2v20250115.Filter{
		{
			Name:   helper.String("endpoint-group-id"),
			Values: []*string{helper.String(egId)},
		},
	}

	var (
		offset uint64 = 0
		limit  uint64 = 10
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

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
				return item, nil
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

// DescribeGa2GlobalAcceleratorById queries a global accelerator instance by its ID.
// Returns (nil, nil) when the instance does not exist.
func (me *Ga2Service) DescribeGa2GlobalAcceleratorById(ctx context.Context, gaId string) (*ga2v20250115.GlobalAcceleratorSet, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeGlobalAcceleratorsRequest()
	request.Filters = []*ga2v20250115.Filter{
		{
			Name:   helper.String("global-accelerator-id"),
			Values: []*string{helper.String(gaId)},
		},
	}

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeGlobalAcceleratorsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeGlobalAcceleratorsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 global accelerators failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 global accelerators failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.GlobalAcceleratorSet
		for i := range set {
			item := set[i]
			if item == nil || item.GlobalAcceleratorId == nil {
				continue
			}

			// Strict equality check: filter values are advisory and may not be enforced server-side.
			if *item.GlobalAcceleratorId == gaId {
				return item, nil
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

// DescribeGa2ListenerById queries a listener by its (gaId, listenerId) tuple.
// Returns (nil, nil) when the listener does not exist.
func (me *Ga2Service) DescribeGa2ListenerById(ctx context.Context, gaId, listenerId string) (*ga2v20250115.ListenerSet, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeListenersRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.Filters = []*ga2v20250115.Filter{
		{
			Name:   helper.String("listener-id"),
			Values: []*string{helper.String(listenerId)},
		},
	}

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeListenersResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeListenersWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 listeners failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 listeners failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.ListenerSet
		for i := range set {
			item := set[i]
			if item == nil || item.ListenerId == nil || item.GlobalAcceleratorId == nil {
				continue
			}

			// Strict equality check on both IDs: filter values are advisory and may not be enforced server-side.
			if *item.ListenerId == listenerId && *item.GlobalAcceleratorId == gaId {
				return item, nil
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

// DescribeGa2ForwardingRuleById queries a forwarding rule by its (gaId, listenerId, policyId, ruleId) tuple.
// Returns (nil, nil) when the rule does not exist.
//
// Note: DescribeForwardingRule is keyed by (GlobalAcceleratorId, ListenerId, ForwardingPolicyId)
// and lacks a per-rule filter slot. We paginate through every rule under the policy and match
// `ForwardingRuleId` strictly client-side.
func (me *Ga2Service) DescribeGa2ForwardingRuleById(ctx context.Context, gaId, listenerId, policyId, ruleId string) (*ga2v20250115.ForwardingRuleSet, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeForwardingRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(policyId)

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeForwardingRuleResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeForwardingRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 forwarding rule failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 forwarding rule failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.ForwardingRuleSet
		for i := range set {
			item := set[i]
			if item == nil || item.ForwardingRuleId == nil {
				continue
			}

			// Defensive: skip items whose parent IDs do not match the request, then strict-equal
			// on the rule ID. The parent IDs *should* match since they were sent on the request,
			// but we guard against API quirks all the same.
			if item.GlobalAcceleratorId != nil && *item.GlobalAcceleratorId != gaId {
				continue
			}
			if item.ListenerId != nil && *item.ListenerId != listenerId {
				continue
			}
			if item.ForwardingPolicyId != nil && *item.ForwardingPolicyId != policyId {
				continue
			}

			if *item.ForwardingRuleId == ruleId {
				return item, nil
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

// describeGa2AccelerateAreas paginates DescribeAccelerateAreas for the given accelerator and
// returns every acceleration region under it. DescribeAccelerateAreas has no per-area filter slot
// (only GlobalAcceleratorId + Offset/Limit), so callers must match the desired item client-side.
func (me *Ga2Service) describeGa2AccelerateAreas(ctx context.Context, gaId string) ([]*ga2v20250115.AcceleratorAreas, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeAccelerateAreasRequest()
	request.GlobalAcceleratorId = helper.String(gaId)

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 100
		areas []*ga2v20250115.AcceleratorAreas
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeAccelerateAreasResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeAccelerateAreasWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 accelerate areas failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 accelerate areas failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.AccelerateAreaSet
		areas = append(areas, set...)

		// Stop when the current page is the last page.
		if uint64(len(set)) < limit {
			break
		}

		offset += limit
	}

	return areas, nil
}

// DescribeGa2AccelerateAreaById queries an acceleration region by its (gaId, areaId) tuple.
// Returns (nil, nil) when the acceleration region does not exist.
func (me *Ga2Service) DescribeGa2AccelerateAreaById(ctx context.Context, gaId, areaId string) (*ga2v20250115.AcceleratorAreas, error) {
	areas, err := me.describeGa2AccelerateAreas(ctx, gaId)
	if err != nil {
		return nil, err
	}

	for i := range areas {
		item := areas[i]
		if item == nil || item.AcceleratorAreaId == nil {
			continue
		}

		// Strict equality check: the API has no per-area filter, so we match client-side.
		if *item.AcceleratorAreaId == areaId {
			return item, nil
		}
	}

	return nil, nil
}

// DescribeGa2AccelerateAreaByRegion queries an acceleration region by its (gaId, region) tuple.
// This is used on Create to resolve the server-generated AcceleratorAreaId, because
// CreateAccelerateAreas only returns a TaskId. Returns (nil, nil) when not found.
func (me *Ga2Service) DescribeGa2AccelerateAreaByRegion(ctx context.Context, gaId, region string) (*ga2v20250115.AcceleratorAreas, error) {
	areas, err := me.describeGa2AccelerateAreas(ctx, gaId)
	if err != nil {
		return nil, err
	}

	for i := range areas {
		item := areas[i]
		if item == nil || item.AccelerateRegion == nil {
			continue
		}

		// Strict equality check on the natural key: the API has no per-area filter, so we match client-side.
		if *item.AccelerateRegion == region {
			return item, nil
		}
	}

	return nil, nil
}

// DescribeGa2ForwardingPolicyById queries a forwarding policy by its (gaId, listenerId, policyId) tuple.
// Returns (nil, nil) when the policy does not exist.
//
// Note: DescribeForwardingPolicy is keyed by (GlobalAcceleratorId, ListenerId) and lacks a per-policy
// filter slot. We paginate through every policy under the listener and match `ForwardingPolicyId`
// strictly client-side.
func (me *Ga2Service) DescribeGa2ForwardingPolicyById(ctx context.Context, gaId, listenerId, policyId string) (*ga2v20250115.ForwardingPolicySet, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeForwardingPolicyRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeForwardingPolicyResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeForwardingPolicyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 forwarding policy failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 forwarding policy failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.ForwardingPolicySet
		for i := range set {
			item := set[i]
			if item == nil || item.ForwardingPolicyId == nil {
				continue
			}

			// Defensive: skip items whose parent IDs do not match the request, then strict-equal
			// on the policy ID. The parent IDs *should* match since they were sent on the request,
			// but we guard against API quirks all the same.
			if item.GlobalAcceleratorId != nil && *item.GlobalAcceleratorId != gaId {
				continue
			}
			if item.ListenerId != nil && *item.ListenerId != listenerId {
				continue
			}

			if *item.ForwardingPolicyId == policyId {
				return item, nil
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

// DescribeGa2GlobalAcceleratorAclRuleById queries an ACL rule by its (policyId, ruleId) tuple.
// Returns (nil, nil) when the rule does not exist.
//
// Note: DescribeGlobalAcceleratorAclRules is keyed by GlobalAcceleratorAclPolicyId only and
// lacks a per-rule filter slot. We paginate through every rule under the policy and match
// GlobalAcceleratorAclRuleId strictly client-side.
func (me *Ga2Service) DescribeGa2GlobalAcceleratorAclRuleById(ctx context.Context, policyId, ruleId string) (*ga2v20250115.GlobalAcceleratorAclRuleSet, error) {
	logId := tccommon.GetLogId(ctx)

	request := ga2v20250115.NewDescribeGlobalAcceleratorAclRulesRequest()
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)

	var (
		offset uint64 = 0
		// limit equals the API-documented maximum to minimize round-trips.
		limit uint64 = 200
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeGlobalAcceleratorAclRulesResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseGa2V20250115Client().DescribeGlobalAcceleratorAclRulesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe ga2 global accelerator acl rules failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s describe ga2 global accelerator acl rules failed, reason:%+v", logId, err)
			return nil, err
		}

		set := response.Response.GlobalAcceleratorAclRuleSet
		for i := range set {
			item := set[i]
			if item == nil || item.GlobalAcceleratorAclRuleId == nil {
				continue
			}

			// Strict equality check on the rule ID: the API has no per-rule filter, so we match client-side.
			if *item.GlobalAcceleratorAclRuleId == ruleId {
				return item, nil
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

// WaitForGa2TaskFinish polls DescribeTaskResult until the task reaches "SUCCESS" or the given timeout elapses.
// The timeout is supplied by the caller because different async operations (create/modify/delete on
// different resource types) may require very different waiting budgets.
func (me *Ga2Service) WaitForGa2TaskFinish(ctx context.Context, taskId string, timeout time.Duration) error {
	if taskId == "" {
		return fmt.Errorf("ga2 task id is empty, cannot poll task result")
	}

	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewDescribeTaskResultRequest()
	request.TaskId = helper.String(taskId)

	err := resource.Retry(timeout, func() *resource.RetryError {
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
