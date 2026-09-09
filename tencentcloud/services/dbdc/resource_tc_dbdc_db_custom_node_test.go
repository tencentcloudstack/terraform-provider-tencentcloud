package dbdc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcdbdc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

func TestAccTencentCloudDbdcDbCustomNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcDbCustomNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "zone", "ap-shanghai-5"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "node_type", "DB.AT5.8XLARGE128"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "node_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "tags.createBy", "Terraform"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "node_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "status"),
				),
			},
			{
				Config: testAccDbdcDbCustomNodeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "tags.createBy", "TerraformUpdate"),
				),
			},
			{
				ResourceName:            "tencentcloud_dbdc_db_custom_node.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"login_settings", "period", "node_count", "auto_voucher", "voucher_ids"},
			},
		},
	})
}

const testAccDbdcDbCustomNode = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone      = "ap-shanghai-5"
  image_id  = "img-xxxxxxxx"
  vpc_id    = "vpc-xxxxxxxx"
  subnet_id = "subnet-xxxxxxxx"
  node_type = "DB.AT5.8XLARGE128"
  period    = 1
  node_name = "tf-example"

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 0

  tags = {
    createBy = "Terraform"
  }
}
`

const testAccDbdcDbCustomNodeUpdate = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone      = "ap-shanghai-5"
  image_id  = "img-xxxxxxxx"
  vpc_id    = "vpc-xxxxxxxx"
  subnet_id = "subnet-xxxxxxxx"
  node_type = "DB.AT5.8XLARGE128"
  period    = 1
  node_name = "tf-example"

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 1

  tags = {
    createBy = "TerraformUpdate"
  }
}
`

// mockMetaDbdcDbCustomNode implements tccommon.ProviderMeta
type mockMetaDbdcDbCustomNode struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdcDbCustomNode) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdcDbCustomNode{}

func ptrDbdcDbCustomNodeInt64(i int64) *int64 {
	return &i
}

func ptrDbdcDbCustomNodeString(s string) *string {
	return &s
}

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomNodeNetworkMode" -v -count=1 -gcflags="all=-l"

// TestDbdcDbCustomNodeNetworkMode_Schema tests the schema definition of network_mode and eni_ip
func TestDbdcDbCustomNodeNetworkMode_Schema(t *testing.T) {
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()

	assert.NotNil(t, res)

	networkMode, ok := res.Schema["network_mode"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, networkMode.Type)
	assert.True(t, networkMode.Computed)
	assert.False(t, networkMode.Optional)

	eniIP, ok := res.Schema["eni_ip"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, eniIP.Type)
	assert.True(t, eniIP.Computed)
	assert.False(t, eniIP.Optional)
}

// TestDbdcDbCustomNodeNetworkMode_Read tests that network_mode and eni_ip are correctly read from DescribeDBCustomNodes response
func TestDbdcDbCustomNodeNetworkMode_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	// Patch DescribeDBCustomNodesWithContext to return response with NetworkMode and EniIP
	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:      ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName:    ptrDbdcDbCustomNodeString("test-node"),
					Zone:        ptrDbdcDbCustomNodeString("ap-shanghai-5"),
					NetworkMode: ptrDbdcDbCustomNodeString("NetworkModeCrossTenantENI"),
					EniIP:       ptrDbdcDbCustomNodeString("10.0.1.10"),
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":      "ap-shanghai-5",
		"node_type": "DB.AT5.8XLARGE128",
	})
	d.SetId("dbcn-test-node-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "NetworkModeCrossTenantENI", d.Get("network_mode").(string))
	assert.Equal(t, "10.0.1.10", d.Get("eni_ip").(string))
}

// TestDbdcDbCustomNodeNetworkMode_ReadNil tests that network_mode and eni_ip are not set when response fields are nil
func TestDbdcDbCustomNodeNetworkMode_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	// Patch DescribeDBCustomNodesWithContext to return response without NetworkMode and EniIP (nil)
	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:      ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName:    ptrDbdcDbCustomNodeString("test-node"),
					Zone:        ptrDbdcDbCustomNodeString("ap-shanghai-5"),
					NetworkMode: nil,
					EniIP:       nil,
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":      "ap-shanghai-5",
		"node_type": "DB.AT5.8XLARGE128",
	})
	d.SetId("dbcn-test-node-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("network_mode").(string))
	assert.Equal(t, "", d.Get("eni_ip").(string))
}

// TestDbdcDbCustomNodeDisasterRecoverGroup_Schema tests the schema definition of disaster_recover_group_ids and disaster_recover_group_id
func TestDbdcDbCustomNodeDisasterRecoverGroup_Schema(t *testing.T) {
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()

	assert.NotNil(t, res)

	disasterRecoverGroupIds, ok := res.Schema["disaster_recover_group_ids"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeList, disasterRecoverGroupIds.Type)
	assert.True(t, disasterRecoverGroupIds.Optional)
	assert.True(t, disasterRecoverGroupIds.ForceNew)
	assert.Equal(t, 1, disasterRecoverGroupIds.MaxItems)
	assert.False(t, disasterRecoverGroupIds.Computed)

	disasterRecoverGroupId, ok := res.Schema["disaster_recover_group_id"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, disasterRecoverGroupId.Type)
	assert.True(t, disasterRecoverGroupId.Computed)
	assert.False(t, disasterRecoverGroupId.Optional)
	assert.False(t, disasterRecoverGroupId.ForceNew)
}

// TestDbdcDbCustomNodeDisasterRecoverGroup_Read tests that disaster_recover_group_id is correctly read from DescribeDBCustomNodes response
func TestDbdcDbCustomNodeDisasterRecoverGroup_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:                 ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName:               ptrDbdcDbCustomNodeString("test-node"),
					Zone:                   ptrDbdcDbCustomNodeString("ap-shanghai-5"),
					DisasterRecoverGroupId: ptrDbdcDbCustomNodeString("dbrg-test-group-id"),
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                       "ap-shanghai-5",
		"node_type":                  "DB.AT5.8XLARGE128",
		"disaster_recover_group_ids": []interface{}{"dbrg-test-group-id"},
	})
	d.SetId("dbcn-test-node-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dbrg-test-group-id", d.Get("disaster_recover_group_id").(string))
}

// TestDbdcDbCustomNodeDisasterRecoverGroup_ReadNil tests that disaster_recover_group_id is not set when the response field is nil
func TestDbdcDbCustomNodeDisasterRecoverGroup_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:                 ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName:               ptrDbdcDbCustomNodeString("test-node"),
					Zone:                   ptrDbdcDbCustomNodeString("ap-shanghai-5"),
					DisasterRecoverGroupId: nil,
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":      "ap-shanghai-5",
		"node_type": "DB.AT5.8XLARGE128",
	})
	d.SetId("dbcn-test-node-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("disaster_recover_group_id").(string))
}

// TestDbdcDbCustomNodeDisasterRecoverGroup_Create tests that disaster_recover_group_ids is correctly mapped into the CreateDBCustomNodes request
func TestDbdcDbCustomNodeDisasterRecoverGroup_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	var capturedRequest *dbdcv20201029.CreateDBCustomNodesRequest
	patches.ApplyMethodFunc(dbdcClient, "CreateDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.CreateDBCustomNodesRequest) (*dbdcv20201029.CreateDBCustomNodesResponse, error) {
		capturedRequest = request
		resp := dbdcv20201029.NewCreateDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.CreateDBCustomNodesResponseParams{
			NodeIds:   []*string{ptrDbdcDbCustomNodeString("dbcn-test-node-id")},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch the read path invoked at the end of Create.
	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:                 ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName:               ptrDbdcDbCustomNodeString("test-node"),
					Zone:                   ptrDbdcDbCustomNodeString("ap-shanghai-5"),
					DisasterRecoverGroupId: ptrDbdcDbCustomNodeString("dbrg-test-group-id"),
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                       "ap-shanghai-5",
		"image_id":                   "img-xxxxxxxx",
		"vpc_id":                     "vpc-xxxxxxxx",
		"subnet_id":                  "subnet-xxxxxxxx",
		"node_type":                  "DB.AT5.8XLARGE128",
		"disaster_recover_group_ids": []interface{}{"dbrg-test-group-id"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dbcn-test-node-id", d.Id())
	// Verify the input parameter was passed through to the API request.
	assert.NotNil(t, capturedRequest.DisasterRecoverGroupIds)
	assert.Len(t, capturedRequest.DisasterRecoverGroupIds, 1)
	assert.Equal(t, "dbrg-test-group-id", *capturedRequest.DisasterRecoverGroupIds[0])
	// Verify the computed output was refreshed during the Create-time Read.
	assert.Equal(t, "dbrg-test-group-id", d.Get("disaster_recover_group_id").(string))
}

// TestDbdcDbCustomNodeDisasterRecoverGroup_CreateWithoutGroup tests that DisasterRecoverGroupIds is not set on the request when the input is absent
func TestDbdcDbCustomNodeDisasterRecoverGroup_CreateWithoutGroup(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseDbdcV20201029Client", dbdcClient)

	var capturedRequest *dbdcv20201029.CreateDBCustomNodesRequest
	patches.ApplyMethodFunc(dbdcClient, "CreateDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.CreateDBCustomNodesRequest) (*dbdcv20201029.CreateDBCustomNodesResponse, error) {
		capturedRequest = request
		resp := dbdcv20201029.NewCreateDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.CreateDBCustomNodesResponseParams{
			NodeIds:   []*string{ptrDbdcDbCustomNodeString("dbcn-test-node-id")},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodesWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodesRequest) (*dbdcv20201029.DescribeDBCustomNodesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodesResponseParams{
			TotalCount: ptrDbdcDbCustomNodeInt64(1),
			NodeSet: []*dbdcv20201029.DBCustomNode{
				{
					NodeId:   ptrDbdcDbCustomNodeString("dbcn-test-node-id"),
					NodeName: ptrDbdcDbCustomNodeString("test-node"),
					Zone:     ptrDbdcDbCustomNodeString("ap-shanghai-5"),
				},
			},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomNodeSecurityGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomNodeSecurityGroupsRequest) (*dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomNodeSecurityGroupsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomNodeSecurityGroupsResponseParams{
			Groups:    []*dbdcv20201029.SecurityGroup{},
			RequestId: ptrDbdcDbCustomNodeString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaDbdcDbCustomNode{client: mockClient}
	res := svcdbdc.ResourceTencentCloudDbdcDbCustomNode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":      "ap-shanghai-5",
		"image_id":  "img-xxxxxxxx",
		"vpc_id":    "vpc-xxxxxxxx",
		"subnet_id": "subnet-xxxxxxxx",
		"node_type": "DB.AT5.8XLARGE128",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dbcn-test-node-id", d.Id())
	assert.Nil(t, capturedRequest.DisasterRecoverGroupIds)
	assert.Equal(t, "", d.Get("disaster_recover_group_id").(string))
}
