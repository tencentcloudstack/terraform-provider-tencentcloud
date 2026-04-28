package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaForLoadBalancer struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForLoadBalancer) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForLoadBalancer{}

func newMockMetaForLoadBalancer() *mockMetaForLoadBalancer {
	return &mockMetaForLoadBalancer{client: &connectivity.TencentCloudClient{}}
}

func ptrStringLB(s string) *string {
	return &s
}

func ptrUint64LB(u uint64) *uint64 {
	return &u
}

// go test ./tencentcloud/services/teo/ -run "TestTeoLoadBalancer" -v -count=1 -gcflags="all=-l"

// TestTeoLoadBalancer_Create_Success tests Create calls API and sets ID
func TestTeoLoadBalancer_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.CreateLoadBalancerRequest) (*teov20220901.CreateLoadBalancerResponse, error) {
		resp := teov20220901.NewCreateLoadBalancerResponse()
		resp.Response = &teov20220901.CreateLoadBalancerResponseParams{
			InstanceId: ptrStringLB("lb-12345678"),
			RequestId:  ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeLoadBalancerList", func(request *teov20220901.DescribeLoadBalancerListRequest) (*teov20220901.DescribeLoadBalancerListResponse, error) {
		resp := teov20220901.NewDescribeLoadBalancerListResponse()
		resp.Response = &teov20220901.DescribeLoadBalancerListResponseParams{
			TotalCount: ptrUint64LB(1),
			LoadBalancerList: []*teov20220901.LoadBalancer{
				{
					InstanceId:     ptrStringLB("lb-12345678"),
					Name:           ptrStringLB("test-lb"),
					Type:           ptrStringLB("HTTP"),
					SteeringPolicy: ptrStringLB("Pritory"),
					FailoverPolicy: ptrStringLB("OtherRecordInOriginGroup"),
					Status:         ptrStringLB("Running"),
					HealthChecker: &teov20220901.HealthChecker{
						Type:     ptrStringLB("HTTP"),
						Port:     ptrUint64LB(80),
						Interval: ptrUint64LB(30),
						Timeout:  ptrUint64LB(5),
						Path:     ptrStringLB("/health"),
					},
					OriginGroupHealthStatus: []*teov20220901.OriginGroupHealthStatus{
						{
							OriginGroupID:   ptrStringLB("og-aaa"),
							OriginGroupName: ptrStringLB("origin-group-a"),
							OriginType:      ptrStringLB("GENERAL"),
							Priority:        ptrStringLB("priority_1"),
						},
					},
				},
			},
			RequestId: ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
		"steering_policy": "Pritory",
		"failover_policy": "OtherRecordInOriginGroup",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#lb-12345678", d.Id())
}

// TestTeoLoadBalancer_Create_APIError tests Create handles API error
func TestTeoLoadBalancer_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.CreateLoadBalancerRequest) (*teov20220901.CreateLoadBalancerResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoLoadBalancer_Create_EmptyInstanceId tests Create handles empty InstanceId
func TestTeoLoadBalancer_Create_EmptyInstanceId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.CreateLoadBalancerRequest) (*teov20220901.CreateLoadBalancerResponse, error) {
		resp := teov20220901.NewCreateLoadBalancerResponse()
		resp.Response = &teov20220901.CreateLoadBalancerResponseParams{
			InstanceId: ptrStringLB(""),
			RequestId:  ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InstanceId is empty")
}

// TestTeoLoadBalancer_Read_Success tests Read retrieves load balancer data
func TestTeoLoadBalancer_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeLoadBalancerList", func(request *teov20220901.DescribeLoadBalancerListRequest) (*teov20220901.DescribeLoadBalancerListResponse, error) {
		resp := teov20220901.NewDescribeLoadBalancerListResponse()
		resp.Response = &teov20220901.DescribeLoadBalancerListResponseParams{
			TotalCount: ptrUint64LB(1),
			LoadBalancerList: []*teov20220901.LoadBalancer{
				{
					InstanceId:     ptrStringLB("lb-12345678"),
					Name:           ptrStringLB("test-lb"),
					Type:           ptrStringLB("HTTP"),
					SteeringPolicy: ptrStringLB("Pritory"),
					FailoverPolicy: ptrStringLB("OtherRecordInOriginGroup"),
					Status:         ptrStringLB("Running"),
					HealthChecker: &teov20220901.HealthChecker{
						Type:     ptrStringLB("HTTP"),
						Port:     ptrUint64LB(80),
						Interval: ptrUint64LB(30),
						Timeout:  ptrUint64LB(5),
						Path:     ptrStringLB("/health"),
					},
					OriginGroupHealthStatus: []*teov20220901.OriginGroupHealthStatus{
						{
							OriginGroupID:   ptrStringLB("og-aaa"),
							OriginGroupName: ptrStringLB("origin-group-a"),
							OriginType:      ptrStringLB("GENERAL"),
							Priority:        ptrStringLB("priority_1"),
						},
					},
					L4UsedList: []*string{ptrStringLB("proxy-1")},
					L7UsedList: []*string{ptrStringLB("example.com")},
					References: []*teov20220901.OriginGroupReference{
						{
							InstanceType: ptrStringLB("acceleration-domain"),
							InstanceId:   ptrStringLB("domain-1"),
							InstanceName: ptrStringLB("example.com"),
						},
					},
				},
			},
			RequestId: ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-lb", d.Get("name"))
	assert.Equal(t, "HTTP", d.Get("type"))
	assert.Equal(t, "Pritory", d.Get("steering_policy"))
	assert.Equal(t, "OtherRecordInOriginGroup", d.Get("failover_policy"))
	assert.Equal(t, "Running", d.Get("status"))
	assert.Equal(t, "lb-12345678", d.Get("instance_id"))
}

// TestTeoLoadBalancer_Read_NotFound tests Read handles load balancer not found
func TestTeoLoadBalancer_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeLoadBalancerList", func(request *teov20220901.DescribeLoadBalancerListRequest) (*teov20220901.DescribeLoadBalancerListResponse, error) {
		resp := teov20220901.NewDescribeLoadBalancerListResponse()
		resp.Response = &teov20220901.DescribeLoadBalancerListResponseParams{
			TotalCount:       ptrUint64LB(0),
			LoadBalancerList: []*teov20220901.LoadBalancer{},
			RequestId:        ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoLoadBalancer_Update_Success tests Update calls ModifyLoadBalancer API
func TestTeoLoadBalancer_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.ModifyLoadBalancerRequest) (*teov20220901.ModifyLoadBalancerResponse, error) {
		resp := teov20220901.NewModifyLoadBalancerResponse()
		resp.Response = &teov20220901.ModifyLoadBalancerResponseParams{
			RequestId: ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeLoadBalancerList", func(request *teov20220901.DescribeLoadBalancerListRequest) (*teov20220901.DescribeLoadBalancerListResponse, error) {
		resp := teov20220901.NewDescribeLoadBalancerListResponse()
		resp.Response = &teov20220901.DescribeLoadBalancerListResponseParams{
			TotalCount: ptrUint64LB(1),
			LoadBalancerList: []*teov20220901.LoadBalancer{
				{
					InstanceId:     ptrStringLB("lb-12345678"),
					Name:           ptrStringLB("test-lb-updated"),
					Type:           ptrStringLB("HTTP"),
					SteeringPolicy: ptrStringLB("Pritory"),
					FailoverPolicy: ptrStringLB("OtherOriginGroup"),
					Status:         ptrStringLB("Running"),
				},
			},
			RequestId: ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb-updated",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
		"steering_policy": "Pritory",
		"failover_policy": "OtherOriginGroup",
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoLoadBalancer_Update_APIError tests Update handles API error
func TestTeoLoadBalancer_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.ModifyLoadBalancerRequest) (*teov20220901.ModifyLoadBalancerResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid InstanceId")
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoLoadBalancer_Delete_Success tests Delete removes load balancer
func TestTeoLoadBalancer_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.DeleteLoadBalancerRequest) (*teov20220901.DeleteLoadBalancerResponse, error) {
		resp := teov20220901.NewDeleteLoadBalancerResponse()
		resp.Response = &teov20220901.DeleteLoadBalancerResponseParams{
			RequestId: ptrStringLB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoLoadBalancer_Delete_APIError tests Delete handles API error
func TestTeoLoadBalancer_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLoadBalancer().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteLoadBalancerWithContext", func(ctx interface{}, request *teov20220901.DeleteLoadBalancerRequest) (*teov20220901.DeleteLoadBalancerResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Load balancer not found")
	})

	meta := newMockMetaForLoadBalancer()
	res := teo.ResourceTencentCloudTeoLoadBalancer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-lb",
		"type":    "HTTP",
		"origin_groups": []interface{}{
			map[string]interface{}{
				"priority":        "priority_1",
				"origin_group_id": "og-aaa",
			},
		},
	})
	d.SetId("zone-1234567890#lb-12345678")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoLoadBalancer_Schema validates schema definition
func TestTeoLoadBalancer_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoLoadBalancer()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "name")
	nameField := res.Schema["name"]
	assert.Equal(t, schema.TypeString, nameField.Type)
	assert.True(t, nameField.Required)

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)
	assert.True(t, typeField.ForceNew)

	assert.Contains(t, res.Schema, "origin_groups")
	originGroups := res.Schema["origin_groups"]
	assert.Equal(t, schema.TypeList, originGroups.Type)
	assert.True(t, originGroups.Required)

	// Check optional fields
	assert.Contains(t, res.Schema, "health_checker")
	healthChecker := res.Schema["health_checker"]
	assert.Equal(t, schema.TypeList, healthChecker.Type)
	assert.True(t, healthChecker.Optional)
	assert.Equal(t, 1, healthChecker.MaxItems)

	assert.Contains(t, res.Schema, "steering_policy")
	steeringPolicy := res.Schema["steering_policy"]
	assert.Equal(t, schema.TypeString, steeringPolicy.Type)
	assert.True(t, steeringPolicy.Optional)

	assert.Contains(t, res.Schema, "failover_policy")
	failoverPolicy := res.Schema["failover_policy"]
	assert.Equal(t, schema.TypeString, failoverPolicy.Type)
	assert.True(t, failoverPolicy.Optional)

	// Check computed fields
	assert.Contains(t, res.Schema, "instance_id")
	instanceId := res.Schema["instance_id"]
	assert.Equal(t, schema.TypeString, instanceId.Type)
	assert.True(t, instanceId.Computed)

	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	assert.Contains(t, res.Schema, "origin_group_health_status")
	ogHealthStatus := res.Schema["origin_group_health_status"]
	assert.Equal(t, schema.TypeList, ogHealthStatus.Type)
	assert.True(t, ogHealthStatus.Computed)

	assert.Contains(t, res.Schema, "l4_used_list")
	l4UsedList := res.Schema["l4_used_list"]
	assert.Equal(t, schema.TypeList, l4UsedList.Type)
	assert.True(t, l4UsedList.Computed)

	assert.Contains(t, res.Schema, "l7_used_list")
	l7UsedList := res.Schema["l7_used_list"]
	assert.Equal(t, schema.TypeList, l7UsedList.Type)
	assert.True(t, l7UsedList.Computed)

	assert.Contains(t, res.Schema, "references")
	references := res.Schema["references"]
	assert.Equal(t, schema.TypeList, references.Type)
	assert.True(t, references.Computed)
}
