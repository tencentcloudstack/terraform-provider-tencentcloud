package cat_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccat "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cat"
)

type mockMetaForCatNodeGroups struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCatNodeGroups) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCatNodeGroups{}

func newMockMetaForCatNodeGroups() *mockMetaForCatNodeGroups {
	return &mockMetaForCatNodeGroups{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCat(s string) *string {
	return &s
}

// go test ./tencentcloud/services/cat/ -run "TestCatNodeGroups" -v -count=1 -gcflags="all=-l"

// TestCatNodeGroups_Read_Success tests successful read with a populated response
func TestCatNodeGroups_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatNodeGroups().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeNodeGroupsWithContext", func(ctx context.Context, request *cat.DescribeNodeGroupsRequest) (*cat.DescribeNodeGroupsResponse, error) {
		resp := cat.NewDescribeNodeGroupsResponse()
		resp.Response = &cat.DescribeNodeGroupsResponseParams{
			NodeList: []*cat.NodeTree{
				{
					ID:      ptrStrCat("group-A"),
					Content: ptrStrCat("Group A"),
					Children: []*cat.NodeLeaf{
						{
							ID:      ptrStrCat("leaf-A1"),
							Content: ptrStrCat("Leaf A1"),
							Children: []*cat.NodeInfoBase{
								{
									ID:      ptrStrCat("node-A1-1"),
									Content: ptrStrCat("Node A1-1"),
								},
							},
						},
					},
				},
			},
			DistrictList: []*cat.DistinctOrNetServiceInfo{
				{
					ID:   ptrStrCat("1"),
					Name: ptrStrCat("Beijing"),
				},
			},
			NetServiceList: []*cat.DistinctOrNetServiceInfo{
				{
					ID:   ptrStrCat("1"),
					Name: ptrStrCat("China Telecom"),
				},
			},
			RequestId: ptrStrCat("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatNodeGroups()
	res := svccat.DataSourceTencentCloudCatNodeGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"node_type": []interface{}{1},
		"ip_type":   1,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeList := d.Get("node_list").([]interface{})
	assert.Equal(t, 1, len(nodeList))
	nodeTree := nodeList[0].(map[string]interface{})
	assert.Equal(t, "group-A", nodeTree["id"])
	assert.Equal(t, "Group A", nodeTree["content"])

	children := nodeTree["children"].([]interface{})
	assert.Equal(t, 1, len(children))
	nodeLeaf := children[0].(map[string]interface{})
	assert.Equal(t, "leaf-A1", nodeLeaf["id"])
	assert.Equal(t, "Leaf A1", nodeLeaf["content"])

	innerChildren := nodeLeaf["children"].([]interface{})
	assert.Equal(t, 1, len(innerChildren))
	nodeInfoBase := innerChildren[0].(map[string]interface{})
	assert.Equal(t, "node-A1-1", nodeInfoBase["id"])
	assert.Equal(t, "Node A1-1", nodeInfoBase["content"])

	districtList := d.Get("district_list").([]interface{})
	assert.Equal(t, 1, len(districtList))
	district := districtList[0].(map[string]interface{})
	assert.Equal(t, "1", district["id"])
	assert.Equal(t, "Beijing", district["name"])

	netServiceList := d.Get("net_service_list").([]interface{})
	assert.Equal(t, 1, len(netServiceList))
	netService := netServiceList[0].(map[string]interface{})
	assert.Equal(t, "1", netService["id"])
	assert.Equal(t, "China Telecom", netService["name"])
}

// TestCatNodeGroups_Read_NilResponse tests read when API returns nil response
func TestCatNodeGroups_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatNodeGroups().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeNodeGroupsWithContext", func(ctx context.Context, request *cat.DescribeNodeGroupsRequest) (*cat.DescribeNodeGroupsResponse, error) {
		resp := cat.NewDescribeNodeGroupsResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaForCatNodeGroups()
	res := svccat.DataSourceTencentCloudCatNodeGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestCatNodeGroups_Read_EmptyResponse tests read when API returns empty but non-nil lists
func TestCatNodeGroups_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatNodeGroups().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeNodeGroupsWithContext", func(ctx context.Context, request *cat.DescribeNodeGroupsRequest) (*cat.DescribeNodeGroupsResponse, error) {
		resp := cat.NewDescribeNodeGroupsResponse()
		resp.Response = &cat.DescribeNodeGroupsResponseParams{
			NodeList:       []*cat.NodeTree{},
			DistrictList:   []*cat.DistinctOrNetServiceInfo{},
			NetServiceList: []*cat.DistinctOrNetServiceInfo{},
			RequestId:      ptrStrCat("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatNodeGroups()
	res := svccat.DataSourceTencentCloudCatNodeGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	nodeList := d.Get("node_list").([]interface{})
	assert.Equal(t, 0, len(nodeList))
	districtList := d.Get("district_list").([]interface{})
	assert.Equal(t, 0, len(districtList))
	netServiceList := d.Get("net_service_list").([]interface{})
	assert.Equal(t, 0, len(netServiceList))
}

// TestCatNodeGroups_Read_APIError tests read when API returns error
func TestCatNodeGroups_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatNodeGroups().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeNodeGroupsWithContext", func(ctx context.Context, request *cat.DescribeNodeGroupsRequest) (*cat.DescribeNodeGroupsResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaForCatNodeGroups()
	res := svccat.DataSourceTencentCloudCatNodeGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestCatNodeGroups_Schema tests the schema definition
func TestCatNodeGroups_Schema(t *testing.T) {
	res := svccat.DataSourceTencentCloudCatNodeGroups()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check optional filter fields
	assert.True(t, res.Schema["node_type"].Optional)
	assert.Equal(t, schema.TypeList, res.Schema["node_type"].Type)
	assert.True(t, res.Schema["task_category"].Optional)
	assert.True(t, res.Schema["ip_type"].Optional)
	assert.True(t, res.Schema["name"].Optional)
	assert.True(t, res.Schema["region_id"].Optional)
	assert.True(t, res.Schema["district_id"].Optional)
	assert.True(t, res.Schema["net_service_id"].Optional)
	assert.True(t, res.Schema["node_group_type"].Optional)
	assert.True(t, res.Schema["task_type"].Optional)
	assert.True(t, res.Schema["probe_type"].Optional)
	assert.True(t, res.Schema["result_output_file"].Optional)

	// Check computed output fields
	assert.True(t, res.Schema["node_list"].Computed)
	assert.True(t, res.Schema["district_list"].Computed)
	assert.True(t, res.Schema["net_service_list"].Computed)
}
