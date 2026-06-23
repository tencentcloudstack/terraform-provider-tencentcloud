package ga2_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcga2 "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

// mockMetaGa2 implements tccommon.ProviderMeta
type mockMetaGa2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2{}

func newMockMetaGa2() *mockMetaGa2 {
	return &mockMetaGa2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGa2(s string) *string {
	return &s
}

func ptrUint64Ga2(u uint64) *uint64 {
	return &u
}

func ptrBoolGa2(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2EndpointGroup_" -v -count=1 -gcflags="all=-l"

// TestGa2EndpointGroup_Create tests that endpoint_group_id and composite ID are set correctly after Create
func TestGa2EndpointGroup_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock CreateEndpointGroupWithContext
	patches.ApplyMethodFunc(ga2Client, "CreateEndpointGroupWithContext",
		func(_ context.Context, req *ga2v20250115.CreateEndpointGroupRequest) (*ga2v20250115.CreateEndpointGroupResponse, error) {
			resp := ga2v20250115.NewCreateEndpointGroupResponse()
			resp.Response = &ga2v20250115.CreateEndpointGroupResponseParams{
				EndpointGroupId: ptrStringGa2("eg-12345678"),
				TaskId:          ptrStringGa2("task-create-001"),
				RequestId:       ptrStringGa2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock WaitForGa2TaskFinish on the service layer
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish",
		func(_ context.Context, taskId string, _ interface{}) error {
			assert.Equal(t, "task-create-001", taskId)
			return nil
		})

	// Mock DescribeGa2EndpointGroupById for the Read call after Create
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2EndpointGroupById",
		func(_ context.Context, gaId, listenerId, egId string) (*ga2v20250115.EndpointGroupConfigurationSet, error) {
			return &ga2v20250115.EndpointGroupConfigurationSet{
				GlobalAcceleratorId: ptrStringGa2("ga2-xxxxxxxx"),
				ListenerId:          ptrStringGa2("lis-xxxxxxxx"),
				EndpointGroupId:     ptrStringGa2("eg-12345678"),
				EndpointGroupType:   ptrStringGa2("DEFAULT"),
				Name:                ptrStringGa2("tf-example"),
				EndpointGroupRegion: ptrStringGa2("ap-guangzhou"),
				Description:         ptrStringGa2("tf example"),
				EnableHealthCheck:   ptrBoolGa2(true),
				CheckType:           ptrStringGa2("HTTP"),
				CheckPort:           ptrUint64Ga2(80),
				CheckPath:           ptrStringGa2("/"),
				CheckMethod:         ptrStringGa2("GET"),
				ConnectTimeout:      ptrUint64Ga2(5000),
				HealthCheckInterval: ptrUint64Ga2(30),
				HealthyThreshold:    ptrUint64Ga2(3),
				UnhealthyThreshold:  ptrUint64Ga2(3),
				ForwardProtocol:     ptrStringGa2("HTTP"),
				EndpointConfigurations: []*ga2v20250115.EndpointConfigurations{
					{
						EndpointType:    ptrStringGa2("PublicIp"),
						EndpointService: ptrStringGa2("1.1.1.1"),
						Weight:          ptrUint64Ga2(10),
					},
				},
			}, nil
		})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-xxxxxxxx",
		"listener_id":           "lis-xxxxxxxx",
		"endpoint_group_type":   "DEFAULT",
		"endpoint_group_configuration": []interface{}{
			map[string]interface{}{
				"name":                  "tf-example",
				"endpoint_group_region": "ap-guangzhou",
				"description":           "tf example",
				"enable_health_check":   true,
				"check_type":            "HTTP",
				"check_port":            "80",
				"check_path":            "/",
				"check_method":          "GET",
				"connect_timeout":       5000,
				"health_check_interval": 30,
				"healthy_threshold":     3,
				"unhealthy_threshold":   3,
				"forward_protocol":      "HTTP",
				"endpoint_configurations": []interface{}{
					map[string]interface{}{
						"endpoint_type":    "PublicIp",
						"endpoint_service": "1.1.1.1",
						"weight":           10,
					},
				},
				"port_overrides": []interface{}{},
				"status_mask":    []interface{}{},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify endpoint_group_id is set correctly
	egId := d.Get("endpoint_group_id").(string)
	assert.Equal(t, "eg-12345678", egId)

	// Verify composite ID format: ga2-xxxxxxxx#lis-xxxxxxxx#eg-12345678
	assert.Equal(t, "ga2-xxxxxxxx#lis-xxxxxxxx#eg-12345678", d.Id())
}

// TestGa2EndpointGroup_Read tests that fields are set correctly after Read
func TestGa2EndpointGroup_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeGa2EndpointGroupById on the service layer
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2EndpointGroupById",
		func(_ context.Context, gaId, listenerId, egId string) (*ga2v20250115.EndpointGroupConfigurationSet, error) {
			assert.Equal(t, "ga2-xxxxxxxx", gaId)
			assert.Equal(t, "lis-xxxxxxxx", listenerId)
			assert.Equal(t, "eg-87654321", egId)
			return &ga2v20250115.EndpointGroupConfigurationSet{
				GlobalAcceleratorId: ptrStringGa2("ga2-xxxxxxxx"),
				ListenerId:          ptrStringGa2("lis-xxxxxxxx"),
				EndpointGroupId:     ptrStringGa2("eg-87654321"),
				EndpointGroupType:   ptrStringGa2("DEFAULT"),
				Name:                ptrStringGa2("tf-read-test"),
				EndpointGroupRegion: ptrStringGa2("ap-beijing"),
				Description:         ptrStringGa2("read test"),
				EnableHealthCheck:   ptrBoolGa2(false),
				CheckType:           ptrStringGa2("TCP"),
				CheckPort:           ptrUint64Ga2(443),
				ConnectTimeout:      ptrUint64Ga2(3000),
				HealthCheckInterval: ptrUint64Ga2(10),
				HealthyThreshold:    ptrUint64Ga2(2),
				UnhealthyThreshold:  ptrUint64Ga2(2),
				ForwardProtocol:     ptrStringGa2("HTTPS"),
				EndpointConfigurations: []*ga2v20250115.EndpointConfigurations{
					{
						EndpointType:      ptrStringGa2("Domain"),
						EndpointService:   ptrStringGa2("example.com"),
						Weight:            ptrUint64Ga2(20),
						HealthCheckStatus: ptrStringGa2("HEALTH"),
					},
				},
				PortOverrides: []*ga2v20250115.PortOverride{
					{
						ListenerPort: ptrUint64Ga2(8080),
						EndpointPort: ptrUint64Ga2(80),
					},
				},
			}, nil
		})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-xxxxxxxx",
		"listener_id":           "lis-xxxxxxxx",
		"endpoint_group_type":   "DEFAULT",
		"endpoint_group_configuration": []interface{}{
			map[string]interface{}{
				"endpoint_configurations": []interface{}{},
				"port_overrides":          []interface{}{},
				"status_mask":             []interface{}{},
			},
		},
	})
	d.SetId("ga2-xxxxxxxx#lis-xxxxxxxx#eg-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify fields are set correctly
	assert.Equal(t, "ga2-xxxxxxxx", d.Get("global_accelerator_id").(string))
	assert.Equal(t, "lis-xxxxxxxx", d.Get("listener_id").(string))
	assert.Equal(t, "eg-87654321", d.Get("endpoint_group_id").(string))
	assert.Equal(t, "DEFAULT", d.Get("endpoint_group_type").(string))
}

// TestGa2EndpointGroup_Read_NotFound tests read when resource is not found
func TestGa2EndpointGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeGa2EndpointGroupById to return nil (not found)
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2EndpointGroupById",
		func(_ context.Context, gaId, listenerId, egId string) (*ga2v20250115.EndpointGroupConfigurationSet, error) {
			return nil, nil
		})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-xxxxxxxx",
		"listener_id":           "lis-xxxxxxxx",
		"endpoint_group_type":   "DEFAULT",
		"endpoint_group_configuration": []interface{}{
			map[string]interface{}{
				"endpoint_configurations": []interface{}{},
				"port_overrides":          []interface{}{},
				"status_mask":             []interface{}{},
			},
		},
	})
	d.SetId("ga2-xxxxxxxx#lis-xxxxxxxx#eg-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestGa2EndpointGroup_Update tests that update calls ModifyEndpointGroup correctly
func TestGa2EndpointGroup_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock ModifyEndpointGroupWithContext
	var capturedRequest *ga2v20250115.ModifyEndpointGroupRequest
	patches.ApplyMethodFunc(ga2Client, "ModifyEndpointGroupWithContext",
		func(_ context.Context, req *ga2v20250115.ModifyEndpointGroupRequest) (*ga2v20250115.ModifyEndpointGroupResponse, error) {
			capturedRequest = req
			resp := ga2v20250115.NewModifyEndpointGroupResponse()
			resp.Response = &ga2v20250115.ModifyEndpointGroupResponseParams{
				TaskId:    ptrStringGa2("task-modify-001"),
				RequestId: ptrStringGa2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock WaitForGa2TaskFinish
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish",
		func(_ context.Context, taskId string, _ interface{}) error {
			assert.Equal(t, "task-modify-001", taskId)
			return nil
		})

	// Mock DescribeGa2EndpointGroupById for the Read call after Update
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "DescribeGa2EndpointGroupById",
		func(_ context.Context, gaId, listenerId, egId string) (*ga2v20250115.EndpointGroupConfigurationSet, error) {
			return &ga2v20250115.EndpointGroupConfigurationSet{
				GlobalAcceleratorId: ptrStringGa2("ga2-xxxxxxxx"),
				ListenerId:          ptrStringGa2("lis-xxxxxxxx"),
				EndpointGroupId:     ptrStringGa2("eg-12345678"),
				EndpointGroupType:   ptrStringGa2("DEFAULT"),
				Name:                ptrStringGa2("tf-example-update"),
				EndpointGroupRegion: ptrStringGa2("ap-guangzhou"),
				Description:         ptrStringGa2("updated description"),
				EnableHealthCheck:   ptrBoolGa2(true),
				CheckType:           ptrStringGa2("HTTP"),
				CheckPort:           ptrUint64Ga2(80),
				ConnectTimeout:      ptrUint64Ga2(5000),
				HealthCheckInterval: ptrUint64Ga2(30),
				HealthyThreshold:    ptrUint64Ga2(3),
				UnhealthyThreshold:  ptrUint64Ga2(3),
				ForwardProtocol:     ptrStringGa2("HTTP"),
				EndpointConfigurations: []*ga2v20250115.EndpointConfigurations{
					{
						EndpointType:    ptrStringGa2("PublicIp"),
						EndpointService: ptrStringGa2("1.1.1.1"),
						Weight:          ptrUint64Ga2(10),
					},
				},
			}, nil
		})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()
	d := res.TestResourceData()
	d.SetId("ga2-xxxxxxxx#lis-xxxxxxxx#eg-12345678")

	// Set old state
	_ = d.Set("global_accelerator_id", "ga2-xxxxxxxx")
	_ = d.Set("listener_id", "lis-xxxxxxxx")
	_ = d.Set("endpoint_group_type", "DEFAULT")
	_ = d.Set("endpoint_group_configuration", []interface{}{
		map[string]interface{}{
			"name":                  "tf-example-update",
			"endpoint_group_region": "ap-guangzhou",
			"description":           "updated description",
			"enable_health_check":   true,
			"check_type":            "HTTP",
			"check_port":            "80",
			"connect_timeout":       5000,
			"health_check_interval": 30,
			"healthy_threshold":     3,
			"unhealthy_threshold":   3,
			"forward_protocol":      "HTTP",
			"endpoint_configurations": []interface{}{
				map[string]interface{}{
					"endpoint_type":    "PublicIp",
					"endpoint_service": "1.1.1.1",
					"weight":           10,
				},
			},
			"port_overrides": []interface{}{},
			"status_mask":    []interface{}{},
		},
	})

	// Mark endpoint_group_configuration as changed to trigger update
	patches.ApplyMethodReturn(d, "HasChange", true)

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the modify request was sent with correct IDs
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "ga2-xxxxxxxx", *capturedRequest.GlobalAcceleratorId)
	assert.Equal(t, "lis-xxxxxxxx", *capturedRequest.ListenerId)
	assert.Equal(t, "eg-12345678", *capturedRequest.EndpointGroupId)
	assert.Equal(t, "tf-example-update", *capturedRequest.Name)
	assert.Equal(t, "updated description", *capturedRequest.Description)
}

// TestGa2EndpointGroup_Delete tests that delete calls DeleteEndpointGroups correctly
func TestGa2EndpointGroup_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DeleteEndpointGroupsWithContext
	var capturedRequest *ga2v20250115.DeleteEndpointGroupsRequest
	patches.ApplyMethodFunc(ga2Client, "DeleteEndpointGroupsWithContext",
		func(_ context.Context, req *ga2v20250115.DeleteEndpointGroupsRequest) (*ga2v20250115.DeleteEndpointGroupsResponse, error) {
			capturedRequest = req
			resp := ga2v20250115.NewDeleteEndpointGroupsResponse()
			resp.Response = &ga2v20250115.DeleteEndpointGroupsResponseParams{
				TaskId:    ptrStringGa2("task-delete-001"),
				RequestId: ptrStringGa2("fake-request-id"),
			}
			return resp, nil
		})

	// Mock WaitForGa2TaskFinish
	patches.ApplyMethodFunc(&svcga2.Ga2Service{}, "WaitForGa2TaskFinish",
		func(_ context.Context, taskId string, _ interface{}) error {
			assert.Equal(t, "task-delete-001", taskId)
			return nil
		})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-xxxxxxxx",
		"listener_id":           "lis-xxxxxxxx",
		"endpoint_group_type":   "DEFAULT",
		"endpoint_group_configuration": []interface{}{
			map[string]interface{}{
				"endpoint_configurations": []interface{}{},
				"port_overrides":          []interface{}{},
				"status_mask":             []interface{}{},
			},
		},
	})
	d.SetId("ga2-xxxxxxxx#lis-xxxxxxxx#eg-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	// Verify the delete request was sent with correct IDs
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "ga2-xxxxxxxx", *capturedRequest.GlobalAcceleratorId)
	assert.Equal(t, "lis-xxxxxxxx", *capturedRequest.ListenerId)
	assert.Equal(t, 1, len(capturedRequest.EndpointGroupIds))
	assert.Equal(t, "eg-12345678", *capturedRequest.EndpointGroupIds[0])
}

// TestGa2EndpointGroup_Schema tests that endpoint_group_id is defined as computed attribute
func TestGa2EndpointGroup_Schema(t *testing.T) {
	res := svcga2.ResourceTencentCloudGa2EndpointGroup()

	assert.NotNil(t, res)

	// Verify endpoint_group_id is computed
	assert.Contains(t, res.Schema, "endpoint_group_id")
	egIdSchema := res.Schema["endpoint_group_id"]
	assert.Equal(t, schema.TypeString, egIdSchema.Type)
	assert.True(t, egIdSchema.Computed)
	assert.False(t, egIdSchema.Optional)
	assert.False(t, egIdSchema.Required)

	// Verify ForceNew fields
	assert.True(t, res.Schema["global_accelerator_id"].ForceNew)
	assert.True(t, res.Schema["listener_id"].ForceNew)
	assert.True(t, res.Schema["endpoint_group_type"].ForceNew)

	// Verify endpoint_group_configuration is required
	assert.True(t, res.Schema["endpoint_group_configuration"].Required)
}
