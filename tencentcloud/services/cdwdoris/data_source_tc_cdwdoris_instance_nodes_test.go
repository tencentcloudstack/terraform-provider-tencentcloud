package cdwdoris_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cdwdorisv20211228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwdoris"
)

type mockMetaCdwdorisInstanceNodesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCdwdorisInstanceNodesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCdwdorisInstanceNodesDS{}

func newMockMetaCdwdorisInstanceNodesDS() *mockMetaCdwdorisInstanceNodesDS {
	return &mockMetaCdwdorisInstanceNodesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrInstanceNodes(s string) *string {
	return &s
}

func ptrInt64InstanceNodes(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/cdwdoris/ -run "TestCdwdorisInstanceNodesDS" -v -count=1 -gcflags="all=-l"

func TestCdwdorisInstanceNodesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdwdorisClient := &cdwdorisv20211228.Client{}
	patches.ApplyMethodReturn(newMockMetaCdwdorisInstanceNodesDS().client, "UseCdwdorisV20211228Client", cdwdorisClient)

	patches.ApplyMethodFunc(cdwdorisClient, "DescribeInstanceNodes", func(request *cdwdorisv20211228.DescribeInstanceNodesRequest) (*cdwdorisv20211228.DescribeInstanceNodesResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "cdwdoris-rhbflamd", *request.InstanceId)
		resp := cdwdorisv20211228.NewDescribeInstanceNodesResponse()
		resp.Response = &cdwdorisv20211228.DescribeInstanceNodesResponseParams{
			TotalCount: ptrInt64InstanceNodes(1),
			InstanceNodesList: []*cdwdorisv20211228.InstanceNode{
				{
					Ip:             ptrStrInstanceNodes("10.0.0.1"),
					Spec:           ptrStrInstanceNodes("S1"),
					Core:           ptrInt64InstanceNodes(8),
					Memory:         ptrInt64InstanceNodes(16),
					DiskType:       ptrStrInstanceNodes("CLOUD_SSD"),
					DiskSize:       ptrInt64InstanceNodes(200),
					Role:           ptrStrInstanceNodes("default_cluster"),
					Status:         ptrStrInstanceNodes("Running"),
					Rip:            ptrStrInstanceNodes("10.0.0.1"),
					FeRole:         ptrStrInstanceNodes("Follower"),
					UUID:           ptrStrInstanceNodes("uuid-abc"),
					Zone:           ptrStrInstanceNodes("ap-guangzhou-3"),
					VirtualZone:    ptrStrInstanceNodes("ap-guangzhou-3"),
					CreateTime:     ptrStrInstanceNodes("2024-01-01 00:00:00"),
					ComputeGroupId: ptrStrInstanceNodes("1"),
				},
			},
			NodeRoles: []*string{
				ptrStrInstanceNodes("data"),
				ptrStrInstanceNodes("master"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaCdwdorisInstanceNodesDS()
	res := cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdwdoris-rhbflamd",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	instanceNodesList := d.Get("instance_nodes_list").([]interface{})
	assert.Len(t, instanceNodesList, 1)

	node0 := instanceNodesList[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.1", node0["ip"].(string))
	assert.Equal(t, "S1", node0["spec"].(string))
	assert.Equal(t, 8, node0["core"].(int))
	assert.Equal(t, 16, node0["memory"].(int))
	assert.Equal(t, "CLOUD_SSD", node0["disk_type"].(string))
	assert.Equal(t, 200, node0["disk_size"].(int))
	assert.Equal(t, "default_cluster", node0["role"].(string))
	assert.Equal(t, "Running", node0["status"].(string))
	assert.Equal(t, "10.0.0.1", node0["rip"].(string))
	assert.Equal(t, "Follower", node0["fe_role"].(string))
	assert.Equal(t, "uuid-abc", node0["uuid"].(string))
	assert.Equal(t, "ap-guangzhou-3", node0["zone"].(string))
	assert.Equal(t, "ap-guangzhou-3", node0["virtual_zone"].(string))
	assert.Equal(t, "2024-01-01 00:00:00", node0["create_time"].(string))
	assert.Equal(t, "1", node0["compute_group_id"].(string))

	nodeRoles := d.Get("node_roles").([]interface{})
	assert.Len(t, nodeRoles, 2)
	assert.Equal(t, "data", nodeRoles[0].(string))
	assert.Equal(t, "master", nodeRoles[1].(string))
}

func TestCdwdorisInstanceNodesDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdwdorisClient := &cdwdorisv20211228.Client{}
	patches.ApplyMethodReturn(newMockMetaCdwdorisInstanceNodesDS().client, "UseCdwdorisV20211228Client", cdwdorisClient)

	var capturedRole *string
	var capturedPolicy *string
	patches.ApplyMethodFunc(cdwdorisClient, "DescribeInstanceNodes", func(request *cdwdorisv20211228.DescribeInstanceNodesRequest) (*cdwdorisv20211228.DescribeInstanceNodesResponse, error) {
		capturedRole = request.NodeRole
		capturedPolicy = request.DisplayPolicy
		resp := cdwdorisv20211228.NewDescribeInstanceNodesResponse()
		resp.Response = &cdwdorisv20211228.DescribeInstanceNodesResponseParams{
			InstanceNodesList: []*cdwdorisv20211228.InstanceNode{
				{
					Ip:   ptrStrInstanceNodes("10.0.0.2"),
					Role: ptrStrInstanceNodes("data"),
				},
			},
			NodeRoles: []*string{ptrStrInstanceNodes("data")},
		}
		return resp, nil
	})

	meta := newMockMetaCdwdorisInstanceNodesDS()
	res := cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "cdwdoris-rhbflamd",
		"node_role":      "data",
		"display_policy": "All",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRole)
	assert.Equal(t, "data", *capturedRole)
	assert.NotNil(t, capturedPolicy)
	assert.Equal(t, "All", *capturedPolicy)

	instanceNodesList := d.Get("instance_nodes_list").([]interface{})
	assert.Len(t, instanceNodesList, 1)
	node0 := instanceNodesList[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.2", node0["ip"].(string))
	assert.Equal(t, "data", node0["role"].(string))
}

func TestCdwdorisInstanceNodesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdwdorisClient := &cdwdorisv20211228.Client{}
	patches.ApplyMethodReturn(newMockMetaCdwdorisInstanceNodesDS().client, "UseCdwdorisV20211228Client", cdwdorisClient)

	patches.ApplyMethodFunc(cdwdorisClient, "DescribeInstanceNodes", func(request *cdwdorisv20211228.DescribeInstanceNodesRequest) (*cdwdorisv20211228.DescribeInstanceNodesResponse, error) {
		resp := cdwdorisv20211228.NewDescribeInstanceNodesResponse()
		resp.Response = &cdwdorisv20211228.DescribeInstanceNodesResponseParams{
			TotalCount:        ptrInt64InstanceNodes(0),
			InstanceNodesList: []*cdwdorisv20211228.InstanceNode{},
		}
		return resp, nil
	})

	meta := newMockMetaCdwdorisInstanceNodesDS()
	res := cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdwdoris-rhbflamd",
	})

	err := res.Read(d, meta)
	// When response InstanceNodesList is empty, the Read function retry block returns NonRetryableError
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

func TestCdwdorisInstanceNodesDS_Schema(t *testing.T) {
	res := cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "instance_id")
	assert.Contains(t, res.Schema, "node_role")
	assert.Contains(t, res.Schema, "display_policy")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "instance_nodes_list")
	assert.Contains(t, res.Schema, "node_roles")

	instanceIdSchema := res.Schema["instance_id"]
	assert.Equal(t, schema.TypeString, instanceIdSchema.Type)
	assert.True(t, instanceIdSchema.Required)

	nodeRoleSchema := res.Schema["node_role"]
	assert.Equal(t, schema.TypeString, nodeRoleSchema.Type)
	assert.True(t, nodeRoleSchema.Optional)

	displayPolicySchema := res.Schema["display_policy"]
	assert.Equal(t, schema.TypeString, displayPolicySchema.Type)
	assert.True(t, displayPolicySchema.Optional)

	instanceNodesListSchema := res.Schema["instance_nodes_list"]
	assert.Equal(t, schema.TypeList, instanceNodesListSchema.Type)
	assert.True(t, instanceNodesListSchema.Computed)

	elemRes := instanceNodesListSchema.Elem.(*schema.Resource)
	expectedFields := []string{
		"ip", "spec", "core", "memory", "disk_type", "disk_size",
		"role", "status", "rip", "fe_role", "uuid", "zone",
		"virtual_zone", "create_time", "compute_group_id",
	}
	for _, field := range expectedFields {
		assert.Contains(t, elemRes.Schema, field)
	}

	nodeRolesSchema := res.Schema["node_roles"]
	assert.Equal(t, schema.TypeList, nodeRolesSchema.Type)
	assert.True(t, nodeRolesSchema.Computed)
}
